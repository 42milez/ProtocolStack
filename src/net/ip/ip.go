package ip

import (
	"bytes"
	"encoding/binary"
	"fmt"
	psErr "github.com/42milez/ProtocolStack/src/error"
	psLog "github.com/42milez/ProtocolStack/src/log"
	"github.com/42milez/ProtocolStack/src/monitor"
	"github.com/42milez/ProtocolStack/src/mw"
	"github.com/42milez/ProtocolStack/src/net"
	"github.com/42milez/ProtocolStack/src/net/arp"
	"github.com/42milez/ProtocolStack/src/repo"
	"github.com/42milez/ProtocolStack/src/worker"
	"sync"
)

const HdrLenMax = 60 // bytes
const HdrLenMin = 20 // bytes
const (
	ICMP mw.ProtocolNumber = 1
	TCP  mw.ProtocolNumber = 6
	UDP  mw.ProtocolNumber = 17
)
const xChBufSize = 5
const ipv4 = 4

var rcvMonCh chan *worker.Message
var rcvSigCh chan *worker.Message
var sndMonCh chan *worker.Message
var sndSigCh chan *worker.Message

var id *PacketID
var receiverID uint32
var senderID uint32

type PacketID struct {
	id  uint16
	mtx sync.Mutex
}

func (p *PacketID) Next() (id uint16) {
	p.mtx.Lock()
	defer p.mtx.Unlock()
	id = p.id
	p.id += 1
	return
}

func Receive(payload []byte, dev mw.IDevice) psErr.E {
	packetLen := len(payload)

	if packetLen < HdrLenMin {
		psLog.E(fmt.Sprintf("ip packet length is too short: %d bytes", packetLen))
		return psErr.InvalidPacket
	}

	buf := bytes.NewBuffer(payload)
	hdr := mw.IpHdr{}
	if err := binary.Read(buf, binary.BigEndian, &hdr); err != nil {
		return psErr.ReadFromBufError
	}

	if version := hdr.VHL >> 4; version != ipv4 {
		psLog.E(fmt.Sprintf("ip version %d is not supported", version))
		return psErr.InvalidProtocolVersion
	}

	hdrLen := int(hdr.VHL&0x0f) * 4
	if packetLen < hdrLen {
		psLog.E(fmt.Sprintf("ip packet length is too short: ihl = %d, actual = %d", hdrLen, packetLen))
		return psErr.InvalidPacket
	}

	if totalLen := int(hdr.TotalLen); packetLen < totalLen {
		psLog.E(fmt.Sprintf("ip packet length is too short: Total Length = %d, Actual Length = %d", totalLen, packetLen))
		return psErr.InvalidPacket
	}

	if hdr.TTL == 0 {
		psLog.E("ttl expired")
		return psErr.TtlExpired
	}

	cs1 := uint16(payload[10])<<8 | uint16(payload[11])
	payload[10] = 0x00 // assign 0 to Header Checksum field (16bit)
	payload[11] = 0x00
	if cs2 := mw.Checksum(payload); cs2 != cs1 {
		psLog.E(fmt.Sprintf("checksum mismatch: Expect = 0x%04x, Actual = 0x%04x", cs1, cs2))
		return psErr.ChecksumMismatch
	}

	iface := repo.IfaceRepo.Lookup(dev, mw.V4AddrFamily)
	if iface == nil {
		psLog.E(fmt.Sprintf("interface for %s is not registered", dev.Name()))
		return psErr.InterfaceNotFound
	}

	if !iface.Unicast.EqualV4(hdr.Dst) {
		if !iface.Broadcast.EqualV4(hdr.Dst) && mw.V4Broadcast.EqualV4(hdr.Dst) {
			psLog.I("ip packet was ignored (it was sent to different address)")
			return psErr.OK
		}
	}

	psLog.D("incoming ip packet", dump(payload)...)

	switch hdr.Protocol {
	case ICMP:
		mw.IcmpRxCh <- &mw.IcmpRxMessage{
			Payload: payload[hdrLen:],
			Dst:     hdr.Dst,
			Src:     hdr.Src,
			Dev:     dev,
		}
	case TCP:
		mw.TcpRxCh <- &mw.TcpRxMessage{
			ProtoNum: uint8(TCP),
			Payload: payload[hdrLen:],
			Dst:     hdr.Dst,
			Src:     hdr.Src,
			Dev:     dev,
		}
		return psErr.Error
	case UDP:
		psLog.E("currently NOT support UDP")
		return psErr.Error
	default:
		psLog.E(fmt.Sprintf("unsupported protocol: %d", hdr.Protocol))
		return psErr.UnsupportedProtocol
	}

	return psErr.OK
}

func Send(protoNum mw.ProtocolNumber, payload []byte, src mw.IP, dst mw.IP) psErr.E {
	var iface *mw.Iface
	var nextHop mw.IP
	var err psErr.E

	// get a next hop
	if iface, nextHop, err = lookupRoute(dst, src); err != psErr.OK {
		psLog.E(fmt.Sprintf("route to %s not found", dst))
		return psErr.RouteNotFound
	}

	if packetLen := HdrLenMin + len(payload); int(iface.Dev.MTU()) < packetLen {
		psLog.E(fmt.Sprintf("ip packet length is too long: %d", packetLen))
		return psErr.PacketTooLong
	}

	packet := createPacket(protoNum, src, dst, payload)
	if packet == nil {
		psLog.E("can't create IP packet")
		return psErr.Error
	}

	psLog.D("outgoing ip packet", dump(packet)...)

	// get eth address from ip address
	var ethAddr mw.EthAddr
	if ethAddr, err = lookupEthAddr(iface, nextHop); err != psErr.OK {
		psLog.E(fmt.Sprintf("ethernet address was not found: %s", err))
		return psErr.NeedRetry
	}

	// send ip packet
	if err = net.Transmit(ethAddr, packet, mw.IPv4, iface); err != psErr.OK {
		return psErr.Error
	}

	return psErr.OK
}

func Start(wg *sync.WaitGroup) psErr.E {
	wg.Add(2)
	go receiver(wg)
	go sender(wg)
	psLog.D("ip service started")
	return psErr.OK
}

func Stop() {
	msg := &worker.Message{
		Desired: worker.Stopped,
	}
	rcvSigCh <- msg
	sndSigCh <- msg
}

func createPacket(protoNum mw.ProtocolNumber, src mw.IP, dst mw.IP, payload []byte) []byte {
	hdr := mw.IpHdr{}
	hdr.VHL = uint8(ipv4<<4) | uint8(HdrLenMin/4)
	hdr.TotalLen = uint16(HdrLenMin + len(payload))
	hdr.ID = id.Next()
	hdr.TTL = 0xff
	hdr.Protocol = protoNum
	copy(hdr.Src[:], src[:])
	copy(hdr.Dst[:], dst[:])

	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, &hdr); err != nil {
		return nil
	}
	if err := binary.Write(buf, binary.BigEndian, &payload); err != nil {
		return nil
	}
	packet := buf.Bytes()

	csum := mw.Checksum(packet)
	packet[10] = uint8((csum & 0xff00) >> 8)
	packet[11] = uint8(csum & 0x00ff)

	return packet
}

func dump(packet []byte) (ret []string) {
	ihl := packet[0] & 0x0f
	totalLen := uint16(packet[2])<<8 | uint16(packet[3])
	payloadLen := totalLen - uint16(4*ihl)

	ret = append(ret, fmt.Sprintf("version:             %d", packet[0]>>4))
	ret = append(ret, fmt.Sprintf("ihl:                 %d", ihl))
	ret = append(ret, fmt.Sprintf("type of service:     0b%08b", packet[1]))
	ret = append(ret, fmt.Sprintf("total length:        %d bytes (payload: %d bytes)", totalLen, payloadLen))
	ret = append(ret, fmt.Sprintf("id:                  %d", uint16(packet[4])<<8|uint16(packet[5])))
	ret = append(ret, fmt.Sprintf("flags:               0b%03b", (packet[6]&0xe0)>>5))
	ret = append(ret, fmt.Sprintf("fragment offset:     %d", uint16(packet[6]&0x1f)<<8|uint16(packet[7])))
	ret = append(ret, fmt.Sprintf("ttl:                 %d", packet[8]))
	ret = append(ret, fmt.Sprintf("protocol:            %s (%d)", mw.ProtocolNumber(packet[9]), packet[9]))
	ret = append(ret, fmt.Sprintf("checksum:            0x%04x", uint16(packet[10])<<8|uint16(packet[11])))
	ret = append(ret, fmt.Sprintf("source address:      %d.%d.%d.%d", packet[12], packet[13], packet[14], packet[15]))
	ret = append(ret, fmt.Sprintf("destination address: %d.%d.%d.%d", packet[16], packet[17], packet[18], packet[19]))

	s := "payload:             "
	for i, v := range packet[ihl:] {
		s += fmt.Sprintf("%02x ", v)
		if (i+1)%20 == 0 {
			s += "\n                                             "
		}
	}
	ret = append(ret, s)

	return
}

func lookupEthAddr(iface *mw.Iface, nextHop mw.IP) (mw.EthAddr, psErr.E) {
	var addr mw.EthAddr
	if iface.Dev.Flag()&mw.NeedArpFlag != 0 {
		if nextHop.Equal(iface.Broadcast) || nextHop.Equal(mw.V4Broadcast) {
			addr = mw.EthBroadcast
		} else {
			var status arp.Status
			if addr, status = arp.Resolve(iface, nextHop); status != arp.Complete {
				return mw.EthAddr{}, psErr.ArpIncomplete
			}
		}
	}
	return addr, psErr.OK
}

func lookupRoute(dst mw.IP, src mw.IP) (*mw.Iface, mw.IP, psErr.E) {
	var iface *mw.Iface
	var nextHop mw.IP

	if src.Equal(mw.V4Any) {
		// Can't determine net address (0.0.0.0 is a non-routable meta-address), so lookup appropriate interface to
		// send IP packet.
		route := repo.RouteRepo.Get(dst)
		if route == nil {
			psLog.E("Route to destination was not found")
			return nil, mw.IP{}, psErr.RouteNotFound
		}
		iface = route.Iface
		if route.NextHop.Equal(mw.V4Any) {
			nextHop = dst
		} else {
			nextHop = route.NextHop
		}
	} else {
		// Source address isn't equal to V4Any means it can determine net address.
		iface = repo.IfaceRepo.Get(src)
		if iface == nil {
			psLog.E(fmt.Sprintf("Interface for %s was not found", src))
			return nil, mw.IP{}, psErr.InterfaceNotFound
		}
		// Don't send IP packet when net address of both destination and iface is not matched each other or
		// destination address is not matched to the broadcast address.
		if !dst.Mask(iface.Netmask).Equal(iface.Unicast.Mask(iface.Netmask)) && !dst.Equal(mw.V4Broadcast) {
			psLog.E(fmt.Sprintf("IP packet can't reach %s (Network address is not matched)", dst.String()))
			return nil, mw.IP{}, psErr.NetworkAddressNotMatch
		}
		nextHop = dst
	}

	return iface, nextHop, psErr.OK
}

func receiver(wg *sync.WaitGroup) {
	defer wg.Done()

	rcvMonCh <- &worker.Message{
		ID:      receiverID,
		Current: worker.Running,
	}

	for {
		select {
		case msg := <-rcvSigCh:
			if msg.Desired == worker.Stopped {
				return
			}
		case msg := <-mw.IpRxCh:
			if err := Receive(msg.Content, msg.Dev); err != psErr.OK {
				return
			}
		}
	}
}

func sender(wg *sync.WaitGroup) {
	defer wg.Done()

	sndMonCh <- &worker.Message{
		ID:      senderID,
		Current: worker.Running,
	}

	for {
		select {
		case msg := <-sndSigCh:
			if msg.Desired == worker.Stopped {
				return
			}
		case msg := <-mw.IpTxCh:
			switch Send(msg.ProtoNum, msg.Packet, msg.Src, msg.Dst) {
			case psErr.OK:
			case psErr.RouteNotFound:
			case psErr.NeedRetry:
				switch msg.ProtoNum {
				case ICMP:
					mw.IcmpDeadLetterQueue <- &mw.IcmpQueueEntry{
						Payload: msg.Packet,
					}
				default:
					psLog.W(fmt.Sprintf("currently NOT support to process unsent message of protocol number %d", msg.ProtoNum))
				}
			default:
				sndMonCh <- &worker.Message{
					ID:      senderID,
					Current: worker.Stopped,
				}
				return
			}
		}
	}
}

func init() {
	rcvMonCh = make(chan *worker.Message, xChBufSize)
	rcvSigCh = make(chan *worker.Message, xChBufSize)
	receiverID = monitor.Register("IP Receiver", rcvMonCh, rcvSigCh)

	sndMonCh = make(chan *worker.Message, xChBufSize)
	sndSigCh = make(chan *worker.Message, xChBufSize)
	senderID = monitor.Register("IP Sender", sndMonCh, sndSigCh)

	id = &PacketID{}
}
