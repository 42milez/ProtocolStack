package ethernet

import "testing"

func TestEthHeader_TypeAsString(t *testing.T) {
	h := EthHeader{Type: 0x0608}
	want := "ARP"
	got := h.TypeAsString()
	if got != want {
		t.Errorf("EthHeader.TypeAsString() = %v; want %v", got, want)
	}

	h = EthHeader{Type: 0x0008}
	want = "IPv4"
	got = h.TypeAsString()
	if got != want {
		t.Errorf("EthHeader.TypeAsString() = %v; want %v", got, want)
	}

	h = EthHeader{Type: 0xdd86}
	want = "IPv6"
	got = h.TypeAsString()
	if got != want {
		t.Errorf("EthHeader.TypeAsString() = %v; want %v", got, want)
	}
}
