package middleware

import (
	"github.com/42milez/ProtocolStack/src/ethernet"
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestGenLoopbackDevice(t *testing.T) {
	want := &ethernet.Device{
		Type:      ethernet.DevTypeLoopback,
		MTU:       ethernet.LoopbackMTU,
		HeaderLen: 0,
		FLAG:      ethernet.DevFlagLoopback,
	}
	got := GenLoopbackDevice()
	if d := cmp.Diff(got, want); d != "" {
		t.Errorf("GenLoopbackDevice() differs: (-got +want)\n%s", d)
	}
}

func TestGenTapDevice_A(t *testing.T) {
	devName := "tap0"
	devEthAddr := ethernet.EthAddr{11, 12, 13, 14, 15, 16}
	want := &ethernet.Device{
		Type:      ethernet.DevTypeEthernet,
		MTU:       ethernet.EthPayloadSizeMax,
		FLAG:      ethernet.DevFlagBroadcast | ethernet.DevFlagNeedArp,
		HeaderLen: ethernet.EthHeaderSize,
		Addr:      devEthAddr,
		Broadcast: ethernet.EthAddrBroadcast,
		Priv:      ethernet.Privilege{FD: -1, Name: devName},
	}
	got := GenTapDevice(0, devEthAddr)
	if d := cmp.Diff(got, want); d != "" {
		t.Errorf("GenTapDevice() differs: (-got +want)\n%s", d)
	}
}
