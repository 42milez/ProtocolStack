package middleware

import (
	"github.com/42milez/ProtocolStack/src/ethernet"
	"github.com/42milez/ProtocolStack/src/network"
	psSyscall "github.com/42milez/ProtocolStack/src/syscall"
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestGenIface(t *testing.T) {
	want := &Iface{
		Family:    network.FamilyV4,
		Unicast:   network.ParseIP(ethernet.LoopbackIpAddr),
		Netmask:   network.ParseIP(ethernet.LoopbackNetmask),
		Broadcast: make(network.IP, 0),
		Network:   make(network.IP, 0),
	}
	got := GenIface(ethernet.LoopbackIpAddr, ethernet.LoopbackNetmask)
	if d := cmp.Diff(got, want); d != "" {
		t.Errorf("GenIface() differs: (-got +want)\n%s", d)
	}
}

func TestGenLoopbackDevice(t *testing.T) {
	want := &ethernet.LoopbackDevice{
		Device: ethernet.Device{
			Name:      "net0",
			Type:      ethernet.DevTypeLoopback,
			MTU:       ethernet.LoopbackMTU,
			HeaderLen: 0,
			FLAG:      ethernet.DevFlagLoopback,
			Syscall:   &psSyscall.Syscall{},
		},
	}
	got := GenLoopbackDevice()
	if d := cmp.Diff(*got, *want); d != "" {
		t.Errorf("GenLoopbackDevice() differs: (-got +want)\n%s", d)
	}
}

func TestGenTapDevice_A(t *testing.T) {
	devName := "tap0"
	devEthAddr := ethernet.EthAddr{11, 12, 13, 14, 15, 16}
	want := &ethernet.TapDevice{
		Device: ethernet.Device{
			Type:      ethernet.DevTypeEthernet,
			MTU:       ethernet.EthPayloadSizeMax,
			FLAG:      ethernet.DevFlagBroadcast | ethernet.DevFlagNeedArp,
			HeaderLen: ethernet.EthHeaderSize,
			Addr:      devEthAddr,
			Broadcast: ethernet.EthAddrBroadcast,
			Priv:      ethernet.Privilege{FD: -1, Name: devName},
			Syscall:   &psSyscall.Syscall{},
		},
	}
	got := GenTapDevice(0, devEthAddr)
	if d := cmp.Diff(*got, *want); d != "" {
		t.Errorf("GenTapDevice() differs: (-got +want)\n%s", d)
	}
}
