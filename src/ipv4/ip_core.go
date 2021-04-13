package ipv4

import "strings"

// IP address lengths (bytes).
const (
	V4AddrLen = 4
	V6AddrLen = 16
)

// An IP is a single IP address.
type IP []byte

// stoi converts string to integer and returns number, characters consumed, and success.
func stoi(s string) (n int, c int, ok bool) {
	n = 0
	c = 0
	for c = 0; c < len(s) && '0' <= s[c] && s[c] <= '9'; c++ {
		n = n * 10 + int(s[c] - '0')
	}
	if c == 0 {
		return 0, 0, false
	}
	return n, c, true
}

// The prefix for the special addresses described in RFC5952.
var v4InV6Prefix = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xff, 0xff}

// v4 creates IP from bytes.
func v4(a, b, c, d byte) IP {
	p := make(IP, V6AddrLen)
	copy(p, v4InV6Prefix)
	p[12] = a
	p[13] = b
	p[14] = c
	p[15] = d
	return p
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
	return v4(p[0], p[1], p[2], p[3])
}

// parseV6 parses string as IPv6 address.
func parseV6(s string) IP {
	// TODO: parse the string as IPv6 address
	return nil
}

// ParseIP parses string as IPv4 or IPv6 address by detecting its format.
func ParseIP(s string) IP {
	if len(s) == 0 {
		return nil
	}

	if strings.Contains(s, ".") {
		return parseV4(s)
	} else {
		return parseV6(s)
	}
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

// ToV4 converts IP to 4 bytes representation.
func (ip IP) ToV4() IP {
	if len(ip) == V4AddrLen {
		return ip
	}
	if len(ip) == V6AddrLen && isZeros(ip) && ip[10] == 0xff && ip[11] == 0xff {
		return ip[12:16]
	}
	return nil
}
