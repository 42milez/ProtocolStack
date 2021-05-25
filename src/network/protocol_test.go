package network

import "testing"

func TestProtocolType_String_SUCCESS(t *testing.T) {
	want := "PROTOCOL_TYPE_ARP"
	got := ProtocolTypeArp.String()
	if got != want {
		t.Errorf("ProtocolType.String() = %v; want %v", got, want)
	}

	want = "PROTOCOL_TYPE_ICMP"
	got = ProtocolTypeIcmp.String()
	if got != want {
		t.Errorf("ProtocolType.String() = %v; want %v", got, want)
	}

	want = "PROTOCOL_TYPE_IP"
	got = ProtocolTypeIp.String()
	if got != want {
		t.Errorf("ProtocolType.String() = %v; want %v", got, want)
	}

	want = "PROTOCOL_TYPE_TCP"
	got = ProtocolTypeTcp.String()
	if got != want {
		t.Errorf("ProtocolType.String() = %v; want %v", got, want)
	}

	want = "PROTOCOL_TYPE_UDP"
	got = ProtocolTypeUdp.String()
	if got != want {
		t.Errorf("ProtocolType.String() = %v; want %v", got, want)
	}

	want = "UNKNOWN"
	got = ProtocolType(999).String()
	if got != want {
		t.Errorf("ProtocolType.String() = %v; want %v", got, want)
	}
}
