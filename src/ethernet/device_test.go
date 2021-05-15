package ethernet

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestMAC_Byte(t *testing.T) {
	mac := MAC("11:22:33:44:55:66")
	want := []byte{11, 22, 33, 44, 55, 66}
	got := mac.Byte()
	if cmp.Equal(got, want) {
		t.Errorf("MAC.Byte() = %v; want %v", got, want)
	}
}
