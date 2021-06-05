package ip

import (
	"bytes"
	"encoding/binary"
	"fmt"
	psErr "github.com/42milez/ProtocolStack/src/error"
	psLog "github.com/42milez/ProtocolStack/src/log"
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
const ipv4 = 4

var RcvRxCh chan *worker.Message
var RcvTxCh chan *worker.Message
var SndRxCh chan *worker.Message
var SndTxCh chan *worker.Message

var id *PacketID

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
		psLog.E(fmt.Sprintf("IP packet length is too short: %d bytes", packetLen))
		return psErr.InvalidPacket
	}

	buf := bytes.NewBuffer(payload)
	hdr := mw.IpHdr{}
	if err := binary.Read(buf, binary.BigEndian, &hdr); err != nil {
		return psErr.ReadFromBufError
	}

	if version := hdr.VHL >> 4; version != ipv4 {
		psLog.E(fmt.Sprintf("IP version %d is not supported", version))
		return psErr.InvalidProtocolVersion
	}

	hdrLen := int(hdr.VHL&0x0f) * 4
	if packetLen < hdrLen {
		psLog.E(fmt.Sprintf("IP packet length is too short: ihl = %d, actual = %d", hdrLen, packetLen))
		return psErr.InvalidPacket
	}

	if totalLen := int(hdr.TotalLen); packetLen < totalLen {
		psLog.E(fmt.Sprintf("IP packet length is too short: Total Length = %d, Actual Length = %d", totalLen, packetLen))
		return psErr.InvalidPacket
	}

	if hdr.TTL == 0 {
		psLog.E("TTL expired")
		return psErr.TtlExpired
	}

	cs1 := uint16(payload[10])<<8 | uint16(payload[11])
	payload[10] = 0x00 // assign 0 to Header Checksum field (16bit)
	payload[11] = 0x00
	if cs2 := mw.Checksum(payload); cs2 != cs1 {
		psLog.E(fmt.Sprintf("Checksum mismatch: Expect = 0x%04x, Actual = 0x%04x", cs1, cs2))
		return psErr.ChecksumMismatch
	}

	iface := repo.IfaceRepo.Lookup(dev, mw.V4AddrFamily)
	if iface == nil {
		psLog.E(fmt.Sprintf("Interface for %s is not registered", dev.Name()))
		return psErr.InterfaceNotFound
	}

	if !iface.Unicast.EqualV4(hdr.Dst) {
		if !iface.Broadcast.EqualV4(hdr.Dst) && mw.V4Broadcast.EqualV4(hdr.Dst) {
			psLog.I("IP packet was ignored (It was sent to different address)")
			return psErr.OK
		}
	}

	psLog.I("Incoming IP packet")
	dumpIpPacket(payload)

	switch hdr.Protocol {
	case ICMP:
		msg := &mw.IcmpRxMessage{
			Payload: payload[hdrLen:],
			Dst:     hdr.Dst,
			Src:     hdr.Src,
			Dev:     dev,
		}
		mw.IcmpRxCh <- msg
	case TCP:
		psLog.E("Currently NOT support TCP")
		return psErr.Error
	case UDP:
		psLog.E("Currently NOT support UDP")
		return psErr.Error
	default:
		psLog.E(fmt.Sprintf("Unsupported protocol: %d", hdr.Protocol))
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
		psLog.E(fmt.Sprintf("Route was not found: %s", err))
		return psErr.Error
	}

	if packetLen := HdrLenMin + len(payload); int(iface.Dev.MTU()) < packetLen {
		psLog.E(fmt.Sprintf("IP packet length is too long: %d", packetLen))
		return psErr.PacketTooLong
	}

	packet := createPacket(protoNum, src, dst, payload)
	if packet == nil {
		psLog.E("Can't create IP packet")
		return psErr.Error
	}

	psLog.I("Outgoing IP packet")
	dumpIpPacket(packet)

	// get eth address from ip address
	var ethAddr mw.EthAddr
	if ethAddr, err = lookupEthAddr(iface, nextHop); err != psErr.OK {
		psLog.E(fmt.Sprintf("Ethernet address was not found: %s", err))
		return psErr.Error
	}

	// send ip packet
	if err = net.Transmit(ethAddr, packet, mw.IPv4, iface); err != psErr.OK {
		return psErr.Error
	}

	return psErr.OK
}

func StartService() {
	go receiver()
	go sender()
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

func dumpIpPacket(packet []byte) {
	ihl := packet[0] & 0x0f
	totalLen := uint16(packet[2])<<8 | uint16(packet[3])
	payloadLen := totalLen - uint16(4*ihl)
	psLog.I(fmt.Sprintf("\tversion:             %d", packet[0]>>4))
	psLog.I(fmt.Sprintf("\tihl:                 %d", ihl))
	psLog.I(fmt.Sprintf("\ttype of service:     0b%08b", packet[1]))
	psLog.I(fmt.Sprintf("\ttotal length:        %d bytes (payload: %d bytes)", totalLen, payloadLen))
	psLog.I(fmt.Sprintf("\tid:                  %d", uint16(packet[4])<<8|uint16(packet[5])))
	psLog.I(fmt.Sprintf("\tflags:               0b%03b", (packet[6]&0xe0)>>5))
	psLog.I(fmt.Sprintf("\tfragment offset:     %d", uint16(packet[6]&0x1f)<<8|uint16(packet[7])))
	psLog.I(fmt.Sprintf("\tttl:                 %d", packet[8]))
	psLog.I(fmt.Sprintf("\tprotocol:            %s (%d)", mw.ProtocolNumber(packet[9]), packet[9]))
	psLog.I(fmt.Sprintf("\tchecksum:            0x%04x", uint16(packet[10])<<8|uint16(packet[11])))
	psLog.I(fmt.Sprintf("\tsource address:      %d.%d.%d.%d", packet[12], packet[13], packet[14], packet[15]))
	psLog.I(fmt.Sprintf("\tdestination address: %d.%d.%d.%d", packet[16], packet[17], packet[18], packet[19]))
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

func receiver() {
	RcvTxCh <- &worker.Message{
		Current: worker.Running,
	}
	for {
		select {
		case msg := <-RcvRxCh:
			if msg.Desired == worker.Stopped {
				RcvTxCh <- &worker.Message{
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

func sender() {
	SndTxCh <- &worker.Message{
		Current: worker.Running,
	}
	for {
		select {
		case msg := <-SndRxCh:
			if msg.Desired == worker.Stopped {
				SndTxCh <- &worker.Message{
					Current: worker.Stopped,
				}
				return
			}
		case msg := <-mw.IpTxCh:
			if err := Send(msg.ProtoNum, msg.Packet, msg.Src, msg.Dst); err != psErr.OK {
				return
			}
		}
	}
}

func init() {
	RcvRxCh = make(chan *worker.Message, 5)
	RcvTxCh = make(chan *worker.Message, 5)
	SndRxCh = make(chan *worker.Message, 5)
	SndTxCh = make(chan *worker.Message, 5)
	id = &PacketID{}
}
