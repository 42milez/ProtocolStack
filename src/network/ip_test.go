package network

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestAddrFamily_String(t *testing.T) {
	got := FamilyV4.String()
	if got != "FAMILY_V4" {
		t.Errorf("AddrFamily.String() = %v; want %v", got, "FAMILY_V4")
	}

	got = FamilyV6.String()
	if got != "FAMILY_V6" {
		t.Errorf("AddrFamily.String() = %v; want %v", got, "FAMILY_V6")
	}

	got = AddrFamily(999).String()
	if got != "UNKNOWN" {
		t.Errorf("AddrFamily.String() = %v; want %v", got, "UNKNOWN")
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
	want := IP{192, 168, 0, 1}
	got := IP{192, 168, 0, 1}.ToV4()
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("IP.ToV4() = %v; want %v; diff %v", got, want, diff)
	}

	got = IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 192, 168, 0, 1}.ToV4()
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("IP.ToV4() = %v; want %v; diff %v", got, want, diff)
	}
}

func TestParseIP(t *testing.T) {
	want := IP{192, 168, 0, 1}
	got := ParseIP("192.168.0.1")
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("ParseIP() = %v; want %v; diff %v", got, want, diff)
	}

	want = nil
	got = ParseIP("")
	if d := cmp.Diff(got, want); d != "" {
		t.Errorf("ParseIP() differs: (-got +want)\n%s", d)
	}
}

func TestV4(t *testing.T) {
	want := IP{192, 168, 0, 1}
	got := V4(192, 168, 0, 1)
	if d := cmp.Diff(got, want); d != "" {
		t.Errorf("V4() differs: (-got +want)\n%s", d)
	}
}
