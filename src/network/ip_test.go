package network

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

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
