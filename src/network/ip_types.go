package network

import (
	"github.com/google/go-cmp/cmp"
	"reflect"
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
	V4Broadcast = V4(255, 255, 255, 255)
	V4Zero      = V4(0, 0, 0, 0)
)

var v4InV6Prefix = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xff, 0xff}

// AddrFamily is IP address family.
type AddrFamily int

// An IP is a single IP address.
type IP []byte

// INTERNET PROTOCOL
// https://datatracker.ietf.org/doc/html/rfc791#page-13
// The number 576 is selected to allow a reasonable sized data block to be transmitted in addition to the required
// header information. For example, this size allows a data block of 512 octets plus 64 header octets to fit in a
// datagram. The maximal internet header is 60 octets, and a typical internet header is 20 octets, allowing a margin for
// headers of higher level protocols.

// Internet Header Format
// https://datatracker.ietf.org/doc/html/rfc791#section-3.1

// IpHeader is an internet protocol header
type IpHeader struct {
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

// ProtocolNumber is assigned internet protocol number
type ProtocolNumber uint8

func (v AddrFamily) String() string {
	return addrFamilies[v]
}

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
func (v IP) ToV4() IP {
	if len(v) == V6AddrLen && isZeros(v[0:10]) && v[10] == 0xff && v[11] == 0xff {
		return v[12:16]
	}
	return v
}

var addrFamilies = map[AddrFamily]string{
	V4AddrFamily: "IPv4",
	V6AddrFamily: "IPv6",
}

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
	63: "any local network",
	64: "SATNET and Backroom EXPAK",
	65: "MIT Subnet Support",
	// 66-68: Unassigned
	69: "SATNET Monitoring",
	71: "Internet Packet Core Utility",
	// 72-75: Unassigned
	76: "Backroom SATNET Monitoring",
	78: "WIDEBAND Monitoring",
	79: "WIDEBAND EXPAK",
	// 80-254: Unassigned
	// 255: Reserved
}
