package network

import (
	"bytes"
	"encoding/binary"
	"fmt"
	psErr "github.com/42milez/ProtocolStack/src/error"
	"github.com/42milez/ProtocolStack/src/ethernet"
	psLog "github.com/42milez/ProtocolStack/src/log"
	"strings"
)

const IpVersionV4 = 0x04

// IpHeaderSizeMin
// IpHeaderSizeMax
// The number 576 is selected to allow a reasonable sized data block to be transmitted in addition to the required
// header information. For example, this size allows a data block of 512 octets plus 64 header octets to fit in a
// datagram. The maximal internet header is 60 octets, and a typical internet header is 20 octets, allowing a margin for
// headers of higher level protocols.
// see: https://datatracker.ietf.org/doc/html/rfc791
const IpHeaderSizeMin = 20
const IpHeaderSizeMax = 60

const (
	FamilyV4 AddrFamily = iota
	FamilyV6
)

// IP address lengths (bytes).
const (
	V4AddrLen = 4
	V6AddrLen = 16
)

var (
	V4Broadcast = V4(255, 255, 255, 255)
	V4Zero      = V4(0, 0, 0, 0)
)

// AddrFamily is IP address family.
type AddrFamily int

func (f AddrFamily) String() string {
	switch f {
	case FamilyV4:
		return "FAMILY_V4"
	case FamilyV6:
		return "FAMILY_V6"
	default:
		return "UNKNOWN"
	}
}

// An IP is a single IP address.
type IP []byte

func (ip IP) EqualV4(v4 [V4AddrLen]byte) bool {
	return ip[0] == v4[0] && ip[1] == v4[1] && ip[2] == v4[2] && ip[3] == v4[3]
}

// String returns the string form of IP.
func (ip IP) String() string {
	const maxIPv4StringLen = len("255.255.255.255")
	b := make(IP, maxIPv4StringLen)

	n := ubtoa(b, 0, ip[0])
	b[n] = '.'
	n++

	n += ubtoa(b, n, ip[1])
	b[n] = '.'
	n++

	n += ubtoa(b, n, ip[2])
	b[n] = '.'
	n++

	n += ubtoa(b, n, ip[3])

	return string(b[:n])
}

// ToV4 converts IP to 4 bytes representation.
func (ip IP) ToV4() IP {
	if len(ip) == V6AddrLen && isZeros(ip[0:10]) && ip[10] == 0xff && ip[11] == 0xff {
		return ip[12:16]
	}
	return ip
}

// Internet Header Format
// https://datatracker.ietf.org/doc/html/rfc791#section-3.1

type IpHeader struct {
	VHL      uint8
	TOS      uint8
	TotalLen uint16
	ID       uint16
	Offset   uint16
	TTL      uint8
	Protocol uint8
	Checksum uint16
	Src      [V4AddrLen]byte
	Dst      [V4AddrLen]byte
	Options  [0]byte
}

func IpInputHandler(payload []byte, dev ethernet.IDevice) psErr.E {
	payloadLen := len(payload)

	if payloadLen < IpHeaderSizeMin {
		psLog.E(fmt.Sprintf("IP packet size is too small: %d bytes", payloadLen))
		return psErr.InvalidPacket
	}

	buf := bytes.NewBuffer(payload)
	hdr := IpHeader{}
	if err := binary.Read(buf, binary.BigEndian, &hdr); err != nil {
		psLog.E(fmt.Sprintf("binary.Read() failed: %s", err))
		return psErr.Error
	}

	if version := hdr.VHL >> 4; version != IpVersionV4 {
		psLog.E(fmt.Sprintf("IP version %d is not supported", version))
		return psErr.UnsupportedVersion
	}

	hdrLen := int(hdr.VHL & 0x0f)

	if payloadLen < hdrLen {
		psLog.E(fmt.Sprintf("IP packet length is too short: IHL = %d, Actual Packet Size = %d", hdrLen, payloadLen))
		return psErr.InvalidPacket
	}

	if totalLen := int(hdr.TotalLen); payloadLen < totalLen {
		psLog.E(fmt.Sprintf("IP packet length is too short: Total Length = %d, Actual Length = %d", totalLen, payloadLen))
		return psErr.InvalidPacket
	}

	if hdr.TTL == 0 {
		psLog.E("TTL expired")
		return psErr.TtlExpired
	}

	if sum := checksum(payload[:20]); sum != hdr.Checksum {
		psLog.E(fmt.Sprintf("Checksum mismatch: Expect = 0x%04x, Actual = 0x%04x", hdr.Checksum, sum))
		return psErr.ChecksumMismatch
	}

	iface := IfaceRepo.Get(dev, FamilyV4)
	if iface == nil {
		psLog.E(fmt.Sprintf("Interface for %s is not registered", dev.DevName()))
		return psErr.InterfaceNotFound
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

// TODO: use IPv4-mapped address above
// V4 creates IP from bytes.
func V4(a, b, c, d byte) IP {
	p := make(IP, V4AddrLen)
	p[0] = a
	p[1] = b
	p[2] = c
	p[3] = d
	return p
}

func checksum(b []byte) uint16 {
	var sum uint32 = 0
	for i := 0; i < len(b); i += 2 {
		if i == 10 {
			continue
		}
		sum += uint32(uint16(b[i])<<8 | uint16(b[i+1]))
	}
	sum = ((sum & 0xffff0000)>>16) + (sum & 0x0000ffff)
	return ^(uint16(sum))
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
