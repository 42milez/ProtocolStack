package middleware

import (
	psErr "github.com/42milez/ProtocolStack/src/error"
	"github.com/42milez/ProtocolStack/src/ethernet"
	psLog "github.com/42milez/ProtocolStack/src/log"
	"github.com/42milez/ProtocolStack/src/network"
	psSyscall "github.com/42milez/ProtocolStack/src/syscall"
	"testing"
)

func TestSetup(t *testing.T) {
	psLog.DisableOutput()
	defer psLog.EnableOutput()

	got := Setup()
	if got.Code != psErr.OK {
		t.Errorf("Setup() = %v; want %v", got, psErr.OK)
	}
}

func TestRegisterDevice(t *testing.T) {
	psLog.DisableOutput()
	defer psLog.EnableOutput()

	dev := &ethernet.TapDevice{
		Device: ethernet.Device{
			Type:      ethernet.DevTypeEthernet,
			MTU:       ethernet.EthPayloadSizeMax,
			FLAG:      ethernet.DevFlagBroadcast | ethernet.DevFlagNeedArp,
			HeaderLen: ethernet.EthHeaderSize,
			Addr:      ethernet.EthAddr{11, 12, 13, 14, 15, 16},
			Broadcast: ethernet.EthAddrBroadcast,
			Priv:      ethernet.Privilege{FD: -1, Name: "tap0"},
			Syscall:   &psSyscall.Syscall{},
		},
	}

	got := RegisterDevice(dev)
	if got.Code != psErr.OK {
		t.Errorf("RegisterDevice() = %v; want %v", got, psErr.OK)
	}
}

func TestRegisterInterface(t *testing.T) {
	psLog.DisableOutput()
	defer psLog.EnableOutput()

	iface := &Iface{
		Family:    network.FamilyV4,
		Unicast:   network.ParseIP(ethernet.LoopbackIpAddr),
		Netmask:   network.ParseIP(ethernet.LoopbackNetmask),
		Broadcast: make(network.IP, 0),
		Network:   make(network.IP, 0),
	}

	dev := &ethernet.TapDevice{
		Device: ethernet.Device{
			Type:      ethernet.DevTypeEthernet,
			MTU:       ethernet.EthPayloadSizeMax,
			FLAG:      ethernet.DevFlagBroadcast | ethernet.DevFlagNeedArp,
			HeaderLen: ethernet.EthHeaderSize,
			Addr:      ethernet.EthAddr{11, 12, 13, 14, 15, 16},
			Broadcast: ethernet.EthAddrBroadcast,
			Priv:      ethernet.Privilege{FD: -1, Name: "tap0"},
			Syscall:   &psSyscall.Syscall{},
		},
	}

	got := RegisterInterface(iface, dev)
	if got.Code != psErr.OK {
		t.Errorf("RegisterInterface() = %v; want %v", got, psErr.OK)
	}
}
