package network

import (
	"bytes"
	"encoding/binary"
	"fmt"
	psErr "github.com/42milez/ProtocolStack/src/error"
	"github.com/42milez/ProtocolStack/src/ethernet"
	psLog "github.com/42milez/ProtocolStack/src/log"
	"strings"
	"sync"
)

const IpHdrLenMin = 20 // bytes
const IpHdrLenMax = 60 // bytes
const ProtoNumICMP = 1
const ProtoNumTCP = 6
const ProtoNumUDP = 17

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

func IpReceive(payload []byte, dev ethernet.IDevice) psErr.E {
	packetLen := len(payload)

	if packetLen < IpHdrLenMin {
		psLog.E(fmt.Sprintf("IP packet length is too short: %d bytes", packetLen))
		return psErr.InvalidPacket
	}

	buf := bytes.NewBuffer(payload)
	hdr := IpHdr{}
	if err := binary.Read(buf, binary.BigEndian, &hdr); err != nil {
		return psErr.Error
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
	if cs2 := checksum(payload); cs2 != cs1 {
		psLog.E(fmt.Sprintf("Checksum mismatch: Expect = 0x%04x, Actual = 0x%04x", cs1, cs2))
		return psErr.ChecksumMismatch
	}

	iface := IfaceRepo.Lookup(dev, V4AddrFamily)
	if iface == nil {
		psLog.E(fmt.Sprintf("Interface for %s is not registered", dev.Name()))
		return psErr.InterfaceNotFound
	}

	if !iface.Unicast.EqualV4(hdr.Dst) {
		if !iface.Broadcast.EqualV4(hdr.Dst) && V4Broadcast.EqualV4(hdr.Dst) {
			psLog.I("Ignored IP packet (It was sent to different address)")
			return psErr.OK
		}
	}

	psLog.I("Incoming IP packet")
	dumpIpPacket(payload)

	switch hdr.Protocol {
	case ProtoNumICMP:
		if err := IcmpReceive(payload[hdrLen:], hdr.Dst, hdr.Src, dev); err != psErr.OK {
			return psErr.Error
		}
	case ProtoNumTCP:
		psLog.E("Currently NOT support TCP")
		return psErr.Error
	case ProtoNumUDP:
		psLog.E("Currently NOT support UDP")
		return psErr.Error
	default:
		psLog.E(fmt.Sprintf("Unsupported protocol: %d", hdr.Protocol))
		return psErr.UnsupportedProtocol
	}

	return psErr.OK
}

func IpSend(protoNum ProtocolNumber, payload []byte, dst IP, src IP) psErr.E {
	var iface *Iface
	var nextHop IP
	var err psErr.E

	// get a next hop
	if iface, nextHop, err = lookupRouting(dst, src); err != psErr.OK {
		psLog.E(fmt.Sprintf("Route was not found: %s", err))
		return psErr.Error
	}

	if packetLen := IpHdrLenMin + len(payload); int(iface.Dev.MTU()) < packetLen {
		psLog.E(fmt.Sprintf("IP packet length is too long: %d", packetLen))
		return psErr.PacketTooLong
	}

	packet := createIpPacket(protoNum, src, dst, payload)
	if packet == nil {
		psLog.E("IP packet was not created")
		return psErr.Error
	}

	psLog.I("Outgoing IP packet")
	dumpIpPacket(packet)

	// get ethernet address from ip address
	var ethAddr ethernet.EthAddr
	if ethAddr, err = lookupEthAddr(iface, nextHop); err != psErr.OK {
		psLog.E(fmt.Sprintf("Ethernet address was not found: %s", err))
		return psErr.Error
	}

	// send ip packet
	if err = Transmit(ethAddr, packet, ethernet.EthTypeIpv4, iface); err != psErr.OK {
		return psErr.Error
	}

	return psErr.OK
}

// ParseIP parses string as IPv4 or IPv6 address by detecting its format.
func ParseIP(s string) IP {
	if strings.Contains(s, ".") {
		return parseV4(s)
	}
	if strings.Contains(s, ":") {
		return parseV6(s)
	}
	return nil
}

// The prefix for the special addresses described in RFC5952.
//var v4InV6Prefix = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xff, 0xff}

// v4 creates IP from bytes.
//func v4(a, b, c, d byte) IP {
//	p := make(IP, V6AddrLen)
//	copy(p, v4InV6Prefix)
//	p[12] = a
//	p[13] = b
//	p[14] = c
//	p[15] = d
//	return p
//}

// V4 creates IP from bytes. TODO: use IPv4-mapped address above
func V4(a, b, c, d byte) IP {
	p := make(IP, V4AddrLen)
	p[0] = a
	p[1] = b
	p[2] = c
	p[3] = d
	return p
}

const ipv4 = 4

var id *PacketID

func allFF(b []byte) bool {
	for _, c := range b {
		if c != 0xff {
			return false
		}
	}
	return true
}

// Computing the Internet Checksum
// https://datatracker.ietf.org/doc/html/rfc1071

func checksum(b []byte) uint16 {
	var sum uint32
	// sum up all fields of IP header by each 16bits (except Header Checksum and Options)
	for i := 0; i < len(b); i += 2 {
		sum += uint32(uint16(b[i])<<8 | uint16(b[i+1]))
	}
	//
	sum = ((sum & 0xffff0000) >> 16) + (sum & 0x0000ffff)
	return ^(uint16(sum))
}

func createIpPacket(protoNum ProtocolNumber, src IP, dst IP, payload []byte) []byte {
	hdr := IpHdr{}
	hdr.VHL = uint8(ipv4<<4) | uint8(IpHdrLenMin/4)
	hdr.TotalLen = uint16(IpHdrLenMin + len(payload))
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

	csum := checksum(packet)
	packet[10] = uint8((csum & 0xff00) >> 8)
	packet[11] = uint8(csum & 0x00ff)

	return packet
}

func dumpIpPacket(packet []byte) {
	ihl := packet[0] & 0x0f
	totalLen := uint16(packet[2])<<8 | uint16(packet[3])
	payloadLen := totalLen - uint16(4*ihl)
	psLog.I(fmt.Sprintf("\tversion:             IPv%d", packet[0]>>4))
	psLog.I(fmt.Sprintf("\tihl:                 %d", ihl))
	psLog.I(fmt.Sprintf("\ttype of service:     0b%08b", packet[1]))
	psLog.I(fmt.Sprintf("\ttotal length:        %d bytes (payload: %d bytes)", totalLen, payloadLen))
	psLog.I(fmt.Sprintf("\tid:                  %d", uint16(packet[4])<<8|uint16(packet[5])))
	psLog.I(fmt.Sprintf("\tflags:               0b%03b", (packet[6]&0xe0)>>5))
	psLog.I(fmt.Sprintf("\tfragment offset:     %d", uint16(packet[6]&0x1f)<<8|uint16(packet[7])))
	psLog.I(fmt.Sprintf("\tttl:                 %d", packet[8]))
	psLog.I(fmt.Sprintf("\tprotocol:            %d (%s)", packet[9], protocolNumbers[ProtocolNumber(packet[9])]))
	psLog.I(fmt.Sprintf("\tchecksum:            0x%04x", uint16(packet[10])<<8|uint16(packet[11])))
	psLog.I(fmt.Sprintf("\tsource address:      %d.%d.%d.%d", packet[12], packet[13], packet[14], packet[15]))
	psLog.I(fmt.Sprintf("\tdestination address: %d.%d.%d.%d", packet[16], packet[17], packet[18], packet[19]))
}

func lookupEthAddr(iface *Iface, nextHop IP) (ethernet.EthAddr, psErr.E) {
	var addr ethernet.EthAddr
	if iface.Dev.Flag()&ethernet.DevFlagNeedArp != 0 {
		if nextHop.Equal(iface.Broadcast) || nextHop.Equal(V4Broadcast) {
			addr = ethernet.EthAddrBroadcast
		} else {
			var status ArpStatus
			if addr, status = arpResolve(iface, nextHop); status != ArpStatusComplete {
				return ethernet.EthAddr{}, psErr.ArpIncomplete
			}
		}
	}
	return addr, psErr.OK
}

func lookupRouting(dst IP, src IP) (*Iface, IP, psErr.E) {
	var iface *Iface
	var nextHop IP

	if src.Equal(V4Zero) {
		// Can't determine network address (0.0.0.0 is a non-routable meta-address), so lookup appropriate interface to
		// send IP packet.
		route := RouteRepo.Get(dst)
		if route == nil {
			psLog.E("Route to destination was not found")
			return nil, IP{}, psErr.RouteNotFound
		}
		iface = route.Iface
		if route.NextHop.Equal(V4Zero) {
			nextHop = dst
		} else {
			nextHop = route.NextHop
		}
	} else {
		// Source address isn't equal to V4Zero means it can determine network address.
		iface = IfaceRepo.Get(src)
		if iface == nil {
			psLog.E(fmt.Sprintf("Interface for %s was not found", src))
			return nil, IP{}, psErr.InterfaceNotFound
		}
		// Don't send IP packet when network address of both destination and iface is not matched each other or
		// destination address is not matched to the broadcast address.
		if !dst.Mask(iface.Netmask).Equal(iface.Unicast.Mask(iface.Netmask)) && !dst.Equal(V4Broadcast) {
			psLog.E(fmt.Sprintf("IP packet can't reach %s (Network address is not matched)", dst.String()))
			return nil, IP{}, psErr.NetworkAddressNotMatch
		}
		nextHop = dst
	}

	return iface, nextHop, psErr.OK
}

// isZeros checks if ip all zeros.
func isZeros(ip IP) bool {
	for i := 0; i < len(ip); i++ {
		if ip[i] != 0 {
			return false
		}
	}
	return true
}

func longestIP(ip1 IP, ip2 IP) IP {
	if len(ip1) != len(ip2) {
		return nil
	}
	for i, v := range ip1 {
		if v < ip2[i] {
			return ip2
		}
	}
	return ip1
}

// parseV4 parses string as IPv4 address.
func parseV4(s string) IP {
	var p [V4AddrLen]byte
	for i := 0; i < V4AddrLen; i++ {
		if i > 0 {
			if s[0] != '.' {
				return nil
			}
			s = s[1:]
		}
		n, c, ok := stoi(s)
		if !ok || n > 0xff {
			return nil
		}
		s = s[c:]
		p[i] = byte(n)
	}
	return V4(p[0], p[1], p[2], p[3])
}

// parseV6 parses string as IPv6 address.
func parseV6(s string) IP {
	// TODO: parse the string as IPv6 address
	return nil
}

// stoi converts string to integer and returns number, characters consumed, and success.
func stoi(s string) (n int, c int, ok bool) {
	n = 0
	for c = 0; c < len(s) && '0' <= s[c] && s[c] <= '9'; c++ {
		n = n*10 + int(s[c]-'0')
	}
	if c == 0 {
		return 0, 0, false
	}
	return n, c, true
}

// ubtoa encodes the string form of the integer v to dst[start:] and
// returns the number of bytes written to dst.
func ubtoa(dst []byte, start int, v byte) int {
	if v < 10 {
		dst[start] = v + '0' // convert a decimal number into ASCII code
		return 1
	}
	if v < 100 {
		dst[start+1] = v%10 + '0'
		dst[start] = v/10 + '0'
		return 2
	}
	dst[start+2] = (v % 10) + '0'
	dst[start+1] = ((v / 10) % 10) + '0'
	dst[start] = (v / 100) + '0'
	return 3
}

func init() {
	id = &PacketID{}
}
