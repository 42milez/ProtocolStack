package network

import "testing"

func TestAddrFamilyString(t *testing.T) {
	v4 := FamilyV4
	v6 := FamilyV6
	got := v4.String()
	if got != "FAMILY_V4" {
		t.Errorf("String() = %v; want \"FAMILY_V4\"", got)
	}
	got = v6.String()
	if got != "FAMILY_V6" {
		t.Errorf("String() = %v; want \"FAMILY_V6\"", got)
	}
}

func TestIsZeros(t *testing.T) {
	got := isZeros(IP{0, 0, 0, 0})
	if !got {
		t.Errorf("isZeros(IP{0, 0, 0, 0}) = %v; want true", got)
	}
}

func TestParseV4(t *testing.T) {
	ip := IP{192, 168, 0, 1}
	got := parseV4("192.168.0.1")
	if got[0] != ip[0] || got[1] != ip[1] || got[2] != ip[2] || got[3] != ip[3] {
		t.Errorf("parseV4(\"192.168.0.1\") = %v; want IP{192, 168, 0, 1}", got)
	}
}

func TestStoi(t *testing.T) {
	n, c, ok := stoi("192.168.0.1")
	if n != 192 || c != 3 || !ok {
		t.Errorf("stoi(\"192.168.0.1\") = %v, %v, %v; wnat 192, 3, true", n, c, ok)
	}
}

func TestUbtoa(t *testing.T) {
	b := make([]byte, 3)
	got := ubtoa(b, 0, 192)
	if got != 3 && b[0] == 49 && b[1] == 57 && b[2] == 50 {
		t.Errorf("ubtoa(b, 0, 192) = %v, b[0] = %v, b[1] = %v, b[2] = %v; want 3, 49, 57, 50", got, b[0], b[1], b[2])
	}
	got = ubtoa(b, 0, 11)
	if got != 2 {
		t.Errorf("ubtoa(b, 0, 11) = %v, b[0] = %v, b[1] = %v; want 2, 49, 49", got, b[0], b[1])
	}
	got = ubtoa(b, 0, 1)
	if got != 1 {
		t.Errorf("ubtoa(b, 0, 1) = %v, b[0] = %v; want 1, 49", got, b[0])
	}
}
