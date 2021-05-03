package network

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestAddrFamily_String(t *testing.T) {
	v4 := FamilyV4
	v6 := FamilyV6
	got := v4.String()
	if got != "FAMILY_V4" {
		t.Errorf("AddrFamily.String() = %v; want \"FAMILY_V4\"", got)
	}
	got = v6.String()
	if got != "FAMILY_V6" {
		t.Errorf("AddrFamily.String() = %v; want \"FAMILY_V6\"", got)
	}
}

func TestParseIP(t *testing.T) {
	want := IP{192, 168, 0, 1}
	got := ParseIP("192.168.0.1")
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("ParseIP(\"192.168.0.1\") = %v; want %v; diff %v", got, want, diff)
	}
}

func TestIP_String(t *testing.T) {
	want := "192.168.0.1"
	ip := IP{192, 168, 0, 1}
	got := ip.String()
	if got != want {
		t.Errorf("IP.String() = %v; want %v", got, want)
	}
}

func TestIP_ToV4(t *testing.T) {
	want := IP{192, 168, 0, 1}
	ip := IP{192, 168, 0, 1}
	got := ip.ToV4()
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("IP.ToV4() = %v; want %v; diff %v", got, want, diff)
	}
	ip = IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 192, 168, 0, 1}
	got = ip.ToV4()
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("IP.ToV4() = %v; want %v; diff %v", got, want, diff)
	}
}

func TestV4(t *testing.T) {
	want := IP{192, 168, 0, 1}
	got := V4(192, 168, 0, 1)
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("V4() = %v; want %v; diff %v", got, want, diff)
	}
}
