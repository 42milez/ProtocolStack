package network

import (
	"github.com/42milez/ProtocolStack/src/device"
	"strings"
)

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

// String returns the string form of IP.
func (ip IP) String() string {
	if len(ip) == 0 {
		return "<nil>"
	}

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
	if len(ip) == V4AddrLen {
		return ip
	}
	if len(ip) == V6AddrLen && isZeros(ip[0:10]) && ip[10] == 0xff && ip[11] == 0xff {
		return ip[12:16]
	}
	return nil
}

// IpInputHandler handles incoming datagram.
func IpInputHandler(data []uint8, dev device.Device) {}

// ParseIP parses string as IPv4 or IPv6 address by detecting its format.
func ParseIP(s string) IP {
	if len(s) == 0 {
		return nil
	}
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
		n = n * 10 + int(s[c] - '0')
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
		dst[start + 1] = v % 10 + '0'
		dst[start] = v / 10 + '0'
		return 2
	}
	dst[start + 2] = (v % 10) + '0'
	dst[start + 1] = ((v / 10) % 10) + '0'
	dst[start] = (v / 100) + '0'
	return 3
}
