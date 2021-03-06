package mw

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestAddrFamily_String(t *testing.T) {
	got := V4AddrFamily.String()
	if got != "IPv4" {
		t.Errorf("AddrFamily.String() = %v; want %v", got, "IPv4")
	}

	got = V6AddrFamily.String()
	if got != "IPv6" {
		t.Errorf("AddrFamily.String() = %v; want %v", got, "IPv6")
	}
}

func TestIP_Equal(t *testing.T) {
	ip1 := IP{192, 168, 0, 1}
	ip2 := IP{192, 168, 0, 1}
	if !ip1.Equal(ip2) {
		t.Errorf("IP.Equal() failed")
	}

	ip2 = IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 192, 168, 0, 1}
	if !ip1.Equal(ip2) {
		t.Errorf("IP.Equal() failed")
	}

	ip1 = IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 192, 168, 0, 1}
	ip2 = IP{192, 168, 0, 1}
	if !ip1.Equal(ip2) {
		t.Errorf("IP.Equal() failed")
	}

	ip1 = IP{192, 168, 0, 2}
	if ip1.Equal(ip2) {
		t.Errorf("IP.Equal() failed")
	}
}

func TestIP_Mask(t *testing.T) {
	ip := IP{192, 168, 0, 1}
	mask := IP{255, 255, 255, 0}

	want := IP{192, 168, 0, 0}
	got := ip.Mask(mask)
	if !got.Equal(want) {
		t.Errorf("IP.Mask() = %s; want %s", got, want)
	}

	mask = IP{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 0}
	got = ip.Mask(mask)
	if !got.Equal(want) {
		t.Errorf("IP.Mask() = %s; want %s", got, want)
	}

	ip = IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 192, 168, 0, 1}
	mask = IP{255, 255, 255, 0}
	got = ip.Mask(mask)
	if !got.Equal(want) {
		t.Errorf("IP.Mask() = %s; want %s", got, want)
	}
}

func TestIP_EqualV4(t *testing.T) {
	ip1 := IP{192, 168, 0, 1}
	ip2 := [V4AddrLen]byte{192, 168, 0, 1}
	if !ip1.EqualV4(ip2) {
		t.Errorf("IP.Equal() failed")
	}

	ip2 = [V4AddrLen]byte{192, 168, 0, 2}
	if ip1.EqualV4(ip2) {
		t.Errorf("IP.Equal() failed")
	}
}

func TestIP_String(t *testing.T) {
	want := "192.168.0.1"
	got := IP{192, 168, 0, 1}.String()
	if got != want {
		t.Errorf("IP.String() = %v; want %v", got, want)
	}
}

func TestIP_ToV4(t *testing.T) {
	want := [V4AddrLen]byte{0xc0, 0xa8, 0x00, 0x01} // 192.168.0.1
	got := IP{192, 168, 0, 1}.ToV4()                // 192.168.0.1
	if d := cmp.Diff(got, want); d != "" {
		t.Errorf("IP.ToV4() differs: (-got +want)\n%s", d)
	}

	got = IP{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0xc0, 0xa8, 0x00, 0x01}.ToV4()
	if d := cmp.Diff(got, want); d != "" {
		t.Errorf("IP.ToV4() differs: (-got +want)\n%s", d)
	}
}

func TestProtocolNumber_String(t *testing.T) {
	want := "ICMP"
	got := ProtocolNumber(1).String()
	if got != want {
		t.Errorf("ProtocolNumber.String() = %s; want %s", got, want)
	}
}

func TestV4Addr_String(t *testing.T) {
	want := "192.168.0.1"
	got := V4Addr{192, 168, 0, 1}.String()
	if got != want {
		t.Errorf("V4Addr.String() = %s; want %s", got, want)
	}
}

func TestChecksum(t *testing.T) {
	b := []byte{
		0x45, 0x00, 0x00, 0x34, 0x51, 0x25, 0x40, 0x00, 0xff, 0x06,
		0x00, 0x00, 0x0a, 0x00, 0x0a, 0xbb, 0x0a, 0x00, 0x03, 0xc3,
	}
	want := uint16(0x0821)
	got := Checksum(b, 0)
	if got != want {
		t.Errorf("Checksum() = 0x%04x; want 0x%04x", got, want)
	}
}

func TestLongestIP(t *testing.T) {
	ip1 := IP{192, 168, 0, 0}
	ip2 := IP{192, 168, 1, 0}
	want := ip2
	got := LongestIP(ip1, ip2)
	if !got.Equal(want) {
		t.Errorf("LongestIP() = %s; want %s", got, want)
	}

	want = ip2
	got = LongestIP(ip2, ip1)
	if !got.Equal(want) {
		t.Errorf("LongestIP() = %s; want %s", got, want)
	}
}

func TestParseIP(t *testing.T) {
	want := IP{192, 168, 0, 1}
	got := ParseIP("192.168.0.1")
	if d := cmp.Diff(got, want); d != "" {
		t.Errorf("ParseIP() differs: (-got +want)\n%s", d)
	}

	want = nil
	got = ParseIP("2001:0db8:85a3:0000:0000:8a2e:0370:7334")
	if !cmp.Equal(got, want) {
		t.Errorf("ParseIP() = %v; want %v", got, want)
	}

	want = nil
	got = ParseIP("")
	if got != nil {
		t.Errorf("ParseIP() = %v; want %v", got, want)
	}
}

func TestV4(t *testing.T) {
	want := IP{192, 168, 0, 1}
	got := V4(192, 168, 0, 1)
	if d := cmp.Diff(got, want); d != "" {
		t.Errorf("V4() differs: (-got +want)\n%s", d)
	}
}
