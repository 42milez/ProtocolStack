package ethernet

import (
	"testing"
)

func TestEthType_String(t *testing.T) {
	ethType := EthType(0x0608)
	want := "ARP"
	got := ethType.String()
	if got != want {
		t.Errorf("EthType.String() = %v; want %v", got, want)
	}

	ethType = EthType(0x0008)
	want = "IPv4"
	got = ethType.String()
	if got != want {
		t.Errorf("EthType.String() = %v; want %v", got, want)
	}

	ethType = EthType(0xdd86)
	want = "IPv6"
	got = ethType.String()
	if got != want {
		t.Errorf("EthType.String() = %v; want %v", got, want)
	}

	ethType = EthType(0x0000)
	want = "UNKNOWN"
	got = ethType.String()
	if got != want {
		t.Errorf("EthType.String() = %v; want %v", got, want)
	}
}
