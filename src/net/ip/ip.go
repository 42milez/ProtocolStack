package ip

import (
	"bytes"
	"encoding/binary"
	"fmt"
	psBinary "github.com/42milez/ProtocolStack/src/binary"
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

func Receive(packet []byte, dev mw.IDevice) psErr.E {
	packetLen := len(packet)

	if packetLen < HdrLenMin {
		psLog.E(fmt.Sprintf("ip packet length is too short: %d bytes", packetLen))
		return psErr.InvalidPacketLength
	}

	buf := bytes.NewBuffer(packet)
	hdr := mw.IpHdr{}
	if err := binary.Read(buf, binary.BigEndian, &hdr); err != nil {
		return psErr.ReadFromBufError
	}

	if version := hdr.VHL >> 4; version != ipv4 {
		psLog.E(fmt.Sprintf("ip version %d is not supported", version))
		return psErr.InvalidProtocolVersion
	}

	hdrLen := int(hdr.VHL&0x0f) << 2
	if packetLen < hdrLen {
		psLog.E(fmt.Sprintf("ip packet length is too short: ihl = %d, actual = %d", hdrLen, packetLen))
		return psErr.InvalidPacketLength
	}

	if totalLen := int(hdr.TotalLen); packetLen < totalLen {
		psLog.E(fmt.Sprintf("ip packet length is too short: Total Length = %d, Actual Length = %d", totalLen, packetLen))
		return psErr.InvalidPacketLength
	}

	if hdr.TTL == 0 {
		psLog.E("ttl expired")
		return psErr.TtlExpired
	}

	if mw.Checksum(packet[:hdrLen], 0) != 0 {
		psLog.E("checksum mismatch (ip)")
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

	psLog.D("incoming ip packet", dump(packet)...)

	switch hdr.Protocol {
	case mw.PnICMP:
		mw.IcmpRxCh <- &mw.IcmpRxMessage{
			Packet: packet[hdrLen:],
			Dst:    hdr.Dst,
			Src:    hdr.Src,
			Dev:    dev,
		}
	case mw.PnTCP:
		mw.TcpRxCh <- &mw.TcpRxMessage{
			ProtoNum:   uint8(mw.PnTCP),
			RawSegment: packet[hdrLen:],
			Dst:        hdr.Dst,
			Src:        hdr.Src,
			Iface:      iface,
		}
	case mw.PnUDP:
		psLog.E("currently NOT support UDP")
		return psErr.Error
	default:
		psLog.E(fmt.Sprintf("unsupported protocol: %d", hdr.Protocol))
		return psErr.UnsupportedProtocol
	}

	return psErr.OK
}

func Send(protoNum mw.ProtocolNumber, data []byte, src mw.IP, dst mw.IP) psErr.E {
	var iface *mw.Iface
	var nextHop mw.IP
	var err psErr.E

	// get a next hop
	if iface, nextHop, err = lookupRoute(dst, src); err != psErr.OK {
		psLog.E(fmt.Sprintf("route to %s not found", dst))
		return psErr.RouteNotFound
	}

	if packetLen := HdrLenMin + len(data); int(iface.Dev.MTU()) < packetLen {
		psLog.E(fmt.Sprintf("ip packet length is too long: %d", packetLen))
		return psErr.PacketTooLong
	}

	packet := createPacket(protoNum, src, dst, data)
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
	if err = net.Transmit(ethAddr, packet, mw.EtIPV4, iface); err != psErr.OK {
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

func createPacket(protoNum mw.ProtocolNumber, src mw.IP, dst mw.IP, data []byte) []byte {
	hdr := mw.IpHdr{}
	hdr.VHL = uint8(ipv4<<4) | uint8(HdrLenMin/4)
	hdr.TotalLen = uint16(HdrLenMin + len(data))
	hdr.ID = id.Next()
	hdr.TTL = 0xff
	hdr.Protocol = protoNum
	copy(hdr.Src[:], src[:])
	copy(hdr.Dst[:], dst[:])

	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, &hdr); err != nil {
		return nil
	}
	if err := binary.Write(buf, binary.BigEndian, &data); err != nil {
		return nil
	}
	packet := buf.Bytes()

	hdrLen := (hdr.VHL & 0x0f) << 2
	csum := mw.Checksum(packet[:hdrLen], 0)
	packet[10] = uint8((csum & 0xff00) >> 8)
	packet[11] = uint8(csum & 0x00ff)

	return packet
}

func dump(packet []byte) (ret []string) {
	hdr := mw.IpHdr{}
	buf := bytes.NewBuffer(packet)
	if err := binary.Read(buf, psBinary.Endian, &hdr); err != nil {
		return nil
	}
	ihl := hdr.VHL & 0x0f
	hdrLen := 4*ihl
	dataLen := hdr.TotalLen - uint16(hdrLen)
	data := buf.Bytes()[hdrLen:]

	v4AddrToString := func(addr [mw.V4AddrLen]byte) string {
		return fmt.Sprintf("%d.%d.%d.%d", addr[0], addr[1], addr[2], addr[3])
	}

	ret = append(ret, fmt.Sprintf("version:             %d", hdr.VHL>>4))
	ret = append(ret, fmt.Sprintf("ihl:                 %d", ihl))
	ret = append(ret, fmt.Sprintf("type of service:     0b%08b", hdr.TOS))
	ret = append(ret, fmt.Sprintf("total length:        %d bytes (data: %d bytes)", hdr.TotalLen, dataLen))
	ret = append(ret, fmt.Sprintf("id:                  %d", hdr.ID))
	ret = append(ret, fmt.Sprintf("flags:               0b%03b", (hdr.Offset&0xe0)>>13))
	ret = append(ret, fmt.Sprintf("fragment offset:     %d", hdr.Offset&0x1f))
	ret = append(ret, fmt.Sprintf("ttl:                 %d", hdr.TTL))
	ret = append(ret, fmt.Sprintf("protocol:            %s (%d)", hdr.Protocol, uint8(hdr.Protocol)))
	ret = append(ret, fmt.Sprintf("checksum:            0x%04x", hdr.Checksum))
	ret = append(ret, fmt.Sprintf("source address:      %s", v4AddrToString(hdr.Src)))
	ret = append(ret, fmt.Sprintf("destination address: %s", v4AddrToString(hdr.Dst)))

	s := "data:                "
	for i, v := range data {
		s += fmt.Sprintf("%02x ", v)
		if (i+1)%20 == 0 && i+1 != len(data) {
			s += "\n                                      "
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
			if addr, status = arp.Resolver.Resolve(iface, nextHop); status != arp.Complete {
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
	defer func() {
		psLog.D("ip receiver stopped")
		wg.Done()
	}()

	rcvMonCh <- &worker.Message{
		ID:      receiverID,
		Current: worker.Running,
	}

	for {
		select {
		case msg := <-rcvSigCh:
			if msg.Desired == worker.Stopped {
				rcvMonCh <- &worker.Message{
					ID:      receiverID,
					Current: worker.Stopped,
				}
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
	defer func() {
		psLog.D("ip sender stopped")
		wg.Done()
	}()

	sndMonCh <- &worker.Message{
		ID:      senderID,
		Current: worker.Running,
	}

	for {
		select {
		case msg := <-sndSigCh:
			if msg.Desired == worker.Stopped {
				sndMonCh <- &worker.Message{
					ID:      senderID,
					Current: worker.Stopped,
				}
				return
			}
		case msg := <-mw.IpTxCh:
			switch Send(msg.ProtoNum, msg.Packet, mw.V4FromByte(msg.Src), mw.V4FromByte(msg.Dst)) {
			case psErr.OK:
			case psErr.RouteNotFound:
			case psErr.NeedRetry:
				switch msg.ProtoNum {
				case mw.PnICMP:
					mw.IcmpDeadLetterQueue <- &mw.IcmpQueueEntry{
						Packet: msg.Packet,
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
