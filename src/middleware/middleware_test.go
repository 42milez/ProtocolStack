package middleware

import (
	psErr "github.com/42milez/ProtocolStack/src/error"
	"github.com/42milez/ProtocolStack/src/ethernet"
	psLog "github.com/42milez/ProtocolStack/src/log"
	mockEthernet "github.com/42milez/ProtocolStack/src/mock/ethernet"
	"github.com/42milez/ProtocolStack/src/network"
	psSyscall "github.com/42milez/ProtocolStack/src/syscall"
	"github.com/golang/mock/gomock"
	"testing"
)

func reset() {
	devices = make([]ethernet.IDevice, 0)
	interfaces = make([]*Iface, 0)
	protocols = make([]Protocol, 0)
}

func TestSetup(t *testing.T) {
	defer reset()

	psLog.DisableOutput()
	defer psLog.EnableOutput()

	got := Setup()
	if got.Code != psErr.OK {
		t.Errorf("Setup() = %v; want %v", got, psErr.OK)
	}
}

func TestRegisterDevice(t *testing.T) {
	defer reset()

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
	//defer reset()

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

func TestUp(t *testing.T) {
	defer reset()

	psLog.DisableOutput()
	defer psLog.EnableOutput()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mockEthernet.NewMockIDevice(ctrl)
	m.EXPECT().Open().Return(psErr.Error{Code: psErr.OK})
	m.EXPECT().Enable()
	m.EXPECT().Info().Return(ethernet.DevTypeEthernet.String(), "net0", "tap0").AnyTimes()
	m.EXPECT().IsUp().Return(false)

	_ = RegisterDevice(m)

	got := Up()
	if got.Code != psErr.OK {
		t.Errorf("Up() = %v; want %v", got, psErr.OK)
	}
}
