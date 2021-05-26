package network

import (
	psErr "github.com/42milez/ProtocolStack/src/error"
	"github.com/42milez/ProtocolStack/src/ethernet"
	psLog "github.com/42milez/ProtocolStack/src/log"
	psSyscall "github.com/42milez/ProtocolStack/src/syscall"
	"github.com/golang/mock/gomock"
	"testing"
)

func reset() {
	DeviceRepo = &deviceRepo{}
	IfaceRepo = &ifaceRepo{}
}

func TestDeviceRepo_Register_SUCCESS(t *testing.T) {
	defer reset()

	psLog.DisableOutput()
	defer psLog.EnableOutput()

	dev := &ethernet.TapDevice{}

	got := DeviceRepo.Register(dev)
	if got.Code != psErr.OK {
		t.Errorf("DeviceRepo.Register() = %v; want %v", got.Code, psErr.OK)
	}
}

func TestDeviceRepo_Register_FAIL_WhenTryingToRegisterSameDevice(t *testing.T) {
	defer reset()

	psLog.DisableOutput()
	defer psLog.EnableOutput()

	dev1 := &ethernet.TapDevice{Device: ethernet.Device{Name: "net0"}}
	dev2 := &ethernet.TapDevice{Device: ethernet.Device{Name: "net0"}}

	_ = DeviceRepo.Register(dev1)
	got := DeviceRepo.Register(dev2)
	if got.Code != psErr.CantRegister {
		t.Errorf("DeviceRepo.Register() = %v; want %v", got.Code, psErr.CantRegister)
	}
}

func TestIfaceRepo_Register_SUCCESS(t *testing.T) {
	defer reset()

	psLog.DisableOutput()
	defer psLog.EnableOutput()

	iface := &Iface{
		Family:    FamilyV4,
		Unicast:   ParseIP(ethernet.LoopbackIpAddr),
		Netmask:   ParseIP(ethernet.LoopbackNetmask),
		Broadcast: make(IP, 0),
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

	got := IfaceRepo.Register(iface, dev)
	if got.Code != psErr.OK {
		t.Errorf("IfaceRepo.Register() = %v; want %v", got.Code, psErr.OK)
	}
}

func TestIfaceRepo_Register_FAIL_WhenTryingToRegisterSameInterface(t *testing.T) {
	defer reset()

	psLog.DisableOutput()
	defer psLog.EnableOutput()

	iface := &Iface{
		Family:    FamilyV4,
		Unicast:   ParseIP(ethernet.LoopbackIpAddr),
		Netmask:   ParseIP(ethernet.LoopbackNetmask),
		Broadcast: make(IP, 0),
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

	_ = IfaceRepo.Register(iface, dev)
	got := IfaceRepo.Register(iface, dev)
	if got.Code != psErr.CantRegister {
		t.Errorf("IfaceRepo.Register() = %v; want %v", got.Code, psErr.CantRegister)
	}
}

func TestUp_SUCCESS(t *testing.T) {
	defer reset()

	psLog.DisableOutput()
	defer psLog.EnableOutput()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := ethernet.NewMockIDevice(ctrl)
	m.EXPECT().Open().Return(psErr.Error{Code: psErr.OK})
	m.EXPECT().Up()
	m.EXPECT().IsUp().Return(false)
	m.EXPECT().Names().Return("net0", "tap0").AnyTimes()
	m.EXPECT().Typ().Return(ethernet.DevTypeEthernet.String()).AnyTimes()

	_ = DeviceRepo.Register(m)

	got := DeviceRepo.Up()
	if got.Code != psErr.OK {
		t.Errorf("Up() = %v; want %v", got.Code, psErr.OK)
	}
}

func TestUp_FailWhenDeviceIsAlreadyOpened(t *testing.T) {
	defer reset()

	psLog.DisableOutput()
	defer psLog.EnableOutput()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := ethernet.NewMockIDevice(ctrl)
	m.EXPECT().IsUp().Return(true)
	m.EXPECT().Names().Return("net0", "tap0")
	m.EXPECT().Typ().Return(ethernet.DevTypeEthernet.String())

	_ = DeviceRepo.Register(m)

	got := DeviceRepo.Up()
	if got.Code != psErr.AlreadyOpened {
		t.Errorf("Up() = %v; want %v", got.Code, psErr.AlreadyOpened)
	}
}

func TestUp_FailWhenCouldNotGetDeviceUp(t *testing.T) {
	defer reset()

	psLog.DisableOutput()
	defer psLog.EnableOutput()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := ethernet.NewMockIDevice(ctrl)
	m.EXPECT().Open().Return(psErr.Error{Code: psErr.CantOpen})
	m.EXPECT().IsUp().Return(false)
	m.EXPECT().Names().Return("net0", "tap0").AnyTimes()
	m.EXPECT().Typ().Return(ethernet.DevTypeEthernet.String()).AnyTimes()

	_ = DeviceRepo.Register(m)

	got := DeviceRepo.Up()
	if got.Code != psErr.CantOpen {
		t.Errorf("Up() = %v; want %v", got.Code, psErr.CantOpen)
	}
}
