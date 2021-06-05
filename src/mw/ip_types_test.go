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
