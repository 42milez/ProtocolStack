package ip

import "testing"

func TestPacketID_Next(t *testing.T) {
	want := uint16(0)
	got := id.Next()
	if got != want {
		t.Errorf("PacketID.Next() = %d; want %d", got, want)
	}
}
