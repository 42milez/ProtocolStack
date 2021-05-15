package ethernet

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestDevType_String(t *testing.T) {
	devType := DevTypeEthernet
	want := "DEVICE_TYPE_ETHERNET"
	got := devType.String()
	if got != want {
		t.Errorf("DevType.String() = %v; want %v", got, want)
	}

	devType = DevTypeLoopback
	want = "DEVICE_TYPE_LOOPBACK"
	got = devType.String()
	if got != want {
		t.Errorf("DevType.String() = %v; want %v", got, want)
	}

	devType = DevTypeNull
	want = "DEVICE_TYPE_NULL"
	got = devType.String()
	if got != want {
		t.Errorf("DevType.String() = %v; want %v", got, want)
	}
}

func TestMAC_Byte(t *testing.T) {
	mac := MAC("11:22:33:44:55:66")
	want := []byte{11, 22, 33, 44, 55, 66}
	got := mac.Byte()
	if cmp.Equal(got, want) {
		t.Errorf("MAC.Byte() = %v; want %v", got, want)
	}
}
