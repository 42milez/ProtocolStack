package network

import (
	"github.com/42milez/ProtocolStack/src/ethernet"
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestGenIface_SUCCESS(t *testing.T) {
	want := &Iface{
		Family:    V4AddrFamily,
		Unicast:   ParseIP(ethernet.LoopbackIpAddr),
		Netmask:   ParseIP(ethernet.LoopbackNetmask),
		Broadcast: ParseIP(ethernet.LoopbackBroadcast),
	}
	got := GenIface(ethernet.LoopbackIpAddr, ethernet.LoopbackNetmask, ethernet.LoopbackBroadcast)
	if d := cmp.Diff(got, want); d != "" {
		t.Errorf("GenIface() differs: (-got +want)\n%s", d)
	}
}

func TestGenLoopbackDevice_SUCCESS(t *testing.T) {
	want := &ethernet.LoopbackDevice{
		Device: ethernet.Device{
			Type:      ethernet.DevTypeLoopback,
			Name:      "net0",
			MTU:       ethernet.LoopbackMTU,
			HeaderLen: 0,
			FLAG:      ethernet.DevFlagLoopback,
		},
	}
	got := GenLoopbackDevice()
	if d := cmp.Diff(*got, *want); d != "" {
		t.Errorf("GenLoopbackDevice() differs: (-got +want)\n%s", d)
	}
}

func TestGenTapDevice_SUCCESS(t *testing.T) {
	devName := "tap0"
	devEthAddr := ethernet.EthAddr{11, 12, 13, 14, 15, 16}
	want := &ethernet.TapDevice{
		Device: ethernet.Device{
			Type:      ethernet.DevTypeEthernet,
			Name:      "net0",
			MTU:       ethernet.EthPayloadSizeMax,
			FLAG:      ethernet.DevFlagBroadcast | ethernet.DevFlagNeedArp,
			HeaderLen: ethernet.EthHeaderSize,
			Addr:      devEthAddr,
			Broadcast: ethernet.EthAddrBroadcast,
			Priv:      ethernet.Privilege{FD: -1, Name: devName},
		},
	}
	got := GenTapDevice(0, devEthAddr)
	if d := cmp.Diff(*got, *want); d != "" {
		t.Errorf("GenTapDevice() differs: (-got +want)\n%s", d)
	}
}
