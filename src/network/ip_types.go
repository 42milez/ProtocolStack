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
func (v IP) ToV4() (ip [V4AddrLen]byte) {
	if len(v) == V6AddrLen && isZeros(v[0:10]) && v[10] == 0xff && v[11] == 0xff {
		copy(ip[:], v[12:16])
		return
	}
	copy(ip[:], v)
	return
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
