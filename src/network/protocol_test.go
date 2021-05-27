package network

import "testing"

func TestProtocolType_String_SUCCESS(t *testing.T) {
	want := "ARP"
	got := ProtocolTypeArp.String()
	if got != want {
		t.Errorf("ProtocolType.String() = %v; want %v", got, want)
	}

	want = "ICMP"
	got = ProtocolTypeIcmp.String()
	if got != want {
		t.Errorf("ProtocolType.String() = %v; want %v", got, want)
	}

	want = "IP"
	got = ProtocolTypeIp.String()
	if got != want {
		t.Errorf("ProtocolType.String() = %v; want %v", got, want)
	}

	want = "TCP"
	got = ProtocolTypeTcp.String()
	if got != want {
		t.Errorf("ProtocolType.String() = %v; want %v", got, want)
	}

	want = "UDP"
	got = ProtocolTypeUdp.String()
	if got != want {
		t.Errorf("ProtocolType.String() = %v; want %v", got, want)
	}

	want = "UNKNOWN"
	got = ProtocolType(100).String()
	if got != want {
		t.Errorf("ProtocolType.String() = %v; want %v", got, want)
	}
}
