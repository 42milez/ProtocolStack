package mw

import (
	"fmt"
	"github.com/google/go-cmp/cmp"
	"reflect"
	"strings"
	"syscall"
)

// IP address family
const (
	V4AddrFamily AddrFamily = syscall.AF_INET
	V6AddrFamily AddrFamily = syscall.AF_INET6
)

// IP address lengths (bytes).
const (
	V4AddrLen = 4
	V6AddrLen = 16
)

// IP address expressions
var (
	V4Any       = V4(0, 0, 0, 0)
	V4Broadcast = V4(255, 255, 255, 255)
)

// ASSIGNED INTERNET PROTOCOL NUMBERS
// https://datatracker.ietf.org/doc/html/rfc790#page-6

var protocolNumbers = map[ProtocolNumber]string{
	// 0: Reserved
	1:  "ICMP",
	3:  "Gateway-to-Gateway",
	4:  "CMCC Gateway Monitoring Message",
	5:  "ST",
	6:  "TCP",
	7:  "UCL",
	9:  "Secure",
	10: "BBN RCC Monitoring",
	11: "NVP",
	12: "PUP",
	13: "Pluribus",
	14: "Telenet",
	15: "XNET",
	16: "Chaos",
	17: "User Datagram",
	18: "Multiplexing",
	19: "DCN",
	20: "TAC Monitoring",
	// 21-62: Unassigned
	63: "any local net",
	64: "SATNET and Backroom EXPAK",
	65: "MIT Subnet Support",
	// 66-68: Unassigned
	69: "SATNET Monitoring",
	71: "Internet EthMessage Core Utility",
	// 72-75: Unassigned
	76: "Backroom SATNET Monitoring",
	78: "WIDEBAND Monitoring",
	79: "WIDEBAND EXPAK",
	// 80-254: Unassigned
	// 255: Reserved
}

var v4InV6Prefix = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xff, 0xff}

// AddrFamily is IP address family.
type AddrFamily int

func (v AddrFamily) String() string {
	return addrFamilies[v]
}

var addrFamilies = map[AddrFamily]string{
	V4AddrFamily: "IPv4",
	V6AddrFamily: "IPv6",
}

// INTERNET PROTOCOL
// https://datatracker.ietf.org/doc/html/rfc791#page-13
// The number 576 is selected to allow a reasonable sized data block to be transmitted in addition to the required
// header information. For example, this size allows a data block of 512 octets plus 64 header octets to fit in a
// datagram. The maximal internet header is 60 octets, and a typical internet header is 20 octets, allowing a margin for
// headers of higher level protocols.

// Internet Header Format
// https://datatracker.ietf.org/doc/html/rfc791#section-3.1

// IpHdr is an internet protocol header
type IpHdr struct {
	VHL      uint8
	TOS      uint8
	TotalLen uint16
	ID       uint16
	Offset   uint16
	TTL      uint8
	Protocol ProtocolNumber
	Checksum uint16
	Src      [V4AddrLen]byte
	Dst      [V4AddrLen]byte
	Options  [0]byte
}

type V4Addr [V4AddrLen]byte

func (v V4Addr) String() string {
	return fmt.Sprintf("%d.%d.%d.%d", v[0], v[1], v[2], v[3])
}

// ProtocolNumber is assigned internet protocol number
type ProtocolNumber uint8

func (v ProtocolNumber) String() string {
	return protocolNumbers[v]
}

// An IP is a single IP address.
type IP []byte

func (v IP) Equal(x IP) bool {
	if len(v) == len(x) {
		return reflect.DeepEqual(v, x)
	}
	if len(v) == V4AddrLen && len(x) == V6AddrLen {
		return reflect.DeepEqual(x[0:12], v4InV6Prefix) && reflect.DeepEqual(v, x[12:])
	}
	if len(v) == V6AddrLen && len(x) == V4AddrLen {
		return reflect.DeepEqual(v[0:12], v4InV6Prefix) && cmp.Equal(v[12:], x)
	}
	return false

}

func (v IP) EqualV4(v4 [V4AddrLen]byte) bool {
	return v[0] == v4[0] && v[1] == v4[1] && v[2] == v4[2] && v[3] == v4[3]
}

func (v IP) Mask(mask IP) IP {
	if len(mask) == V6AddrLen && len(v) == V4AddrLen && allFF(mask[:12]) {
		mask = mask[12:]
	}
	if len(mask) == V4AddrLen && len(v) == V6AddrLen && cmp.Equal(v[:12], v4InV6Prefix) {
		v = v[12:]
	}
	n := len(v)
	if n != len(mask) {
		return nil
	}
	ret := make(IP, n)
	for i := 0; i < n; i++ {
		ret[i] = v[i] & mask[i]
	}
	return ret

}

// String returns the string form of IP.
func (v IP) String() string {
	const maxIPv4StringLen = len("255.255.255.255")
	b := make(IP, maxIPv4StringLen)

	n := ubtoa(b, 0, v[0])
	b[n] = '.'
	n++

	n += ubtoa(b, n, v[1])
	b[n] = '.'
	n++

	n += ubtoa(b, n, v[2])
	b[n] = '.'
	n++

	n += ubtoa(b, n, v[3])

	return string(b[:n])
}

// ToV4 converts IP to 4 bytes representation.
func (v IP) ToV4() (ip [V4AddrLen]byte) {
	if len(v) == V6AddrLen && isZeros(v[0:10]) && v[10] == 0xff && v[11] == 0xff {
		copy(ip[:], v[12:16])
		return
	}
	copy(ip[:], v)
	return
}

// Computing the Internet Checksum
// https://datatracker.ietf.org/doc/html/rfc1071

func Checksum(b []byte) uint16 {
	var sum uint32
	// sum up all fields of IP header by each 16bits (except Header Checksum and Options)
	for i := 0; i < len(b); i += 2 {
		sum += uint32(uint16(b[i])<<8 | uint16(b[i+1]))
	}
	//
	sum = ((sum & 0xffff0000) >> 16) + (sum & 0x0000ffff)
	return ^(uint16(sum))
}

func LongestIP(ip1 IP, ip2 IP) IP {
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

func allFF(b []byte) bool {
	for _, c := range b {
		if c != 0xff {
			return false
		}
	}
	return true
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
