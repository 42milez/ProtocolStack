package network

import "testing"

func TestIsZeros(t *testing.T) {
	ip := make(IP, V4AddrLen)
	ret := isZeros(ip)
	if !ret {
		t.Errorf("isZeros() -> expects: %v / actually: %v", ret, true)
	}
}

func TestParseV4(t *testing.T) {
	s := "192.168.0.1"
	ip := make(IP, V4AddrLen)
	ip[0] = 192
	ip[1] = 168
	ip[2] = 0
	ip[3] = 1
	ret := parseV4(s)
	if ret[0] != ip[0] || ret[1] != ip[1] || ret[2] != ip[2] || ret[3] != ip[3] {
		t.Errorf("parseV4() -> expects: %v / actually: %v", ret, ip)
	}
}

func TestStoi(t *testing.T) {
	n, c, ok := stoi("192.168.0.1")
	if n != 192 || c != 3 || !ok {
		t.Errorf("stoi() -> expects: %v, %v, %v / actually: %v, %v, %v", 192, 3, true, n, c, ok)
	}
}
