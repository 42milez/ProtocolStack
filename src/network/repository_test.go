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

func TestDeviceRepo_Register_1(t *testing.T) {
	defer reset()

	psLog.DisableOutput()
	defer psLog.EnableOutput()

	dev := &ethernet.TapDevice{}

	got := DeviceRepo.Register(dev)
	if got != psErr.OK {
		t.Errorf("DeviceRepo.Register() = %s; want %s", got, psErr.OK)
	}
}

// Fail when it's trying to register same device.
func TestDeviceRepo_Register_2(t *testing.T) {
	defer reset()

	psLog.DisableOutput()
	defer psLog.EnableOutput()

	dev1 := &ethernet.TapDevice{Device: ethernet.Device{Name: "net0"}}
	dev2 := &ethernet.TapDevice{Device: ethernet.Device{Name: "net0"}}

	_ = DeviceRepo.Register(dev1)
	got := DeviceRepo.Register(dev2)
	if got != psErr.Error {
		t.Errorf("DeviceRepo.Register() = %s; want %s", got, psErr.Error)
	}
}

func TestIfaceRepo_Register_1(t *testing.T) {
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
	if got != psErr.OK {
		t.Errorf("IfaceRepo.Register() = %v; want %v", got, psErr.OK)
	}
}

// Fail when it's trying to register same interface.
func TestIfaceRepo_Register_2(t *testing.T) {
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
	if got != psErr.Error {
		t.Errorf("IfaceRepo.Register() = %s; want %s", got, psErr.Error)
	}
}

func TestUp_1(t *testing.T) {
	defer reset()

	psLog.DisableOutput()
	defer psLog.EnableOutput()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := ethernet.NewMockIDevice(ctrl)
	m.EXPECT().Open().Return(psErr.OK)
	m.EXPECT().Up()
	m.EXPECT().IsUp().Return(false)
	m.EXPECT().EthAddr().Return(ethernet.EthAddr{})
	m.EXPECT().BroadcastEthAddr().Return(ethernet.EthAddr{})
	m.EXPECT().PeerEthAddr().Return(ethernet.EthAddr{})
	m.EXPECT().DevType().Return(ethernet.DevTypeEthernet).AnyTimes()
	m.EXPECT().DevName().Return("net0").AnyTimes()
	m.EXPECT().PrivDevName().Return("tap0").AnyTimes()

	_ = DeviceRepo.Register(m)

	got := DeviceRepo.Up()
	if got != psErr.OK {
		t.Errorf("Up() = %v; want %v", got, psErr.OK)
	}
}

// Fail when device is already opened.
func TestUp_2(t *testing.T) {
	defer reset()

	psLog.DisableOutput()
	defer psLog.EnableOutput()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := ethernet.NewMockIDevice(ctrl)
	m.EXPECT().IsUp().Return(true)
	m.EXPECT().DevType().Return(ethernet.DevTypeEthernet).AnyTimes()
	m.EXPECT().DevName().Return("net0").AnyTimes()
	m.EXPECT().PrivDevName().Return("tap0").AnyTimes()
	m.EXPECT().EthAddr().Return(ethernet.EthAddr{})
	m.EXPECT().BroadcastEthAddr().Return(ethernet.EthAddr{})
	m.EXPECT().PeerEthAddr().Return(ethernet.EthAddr{})

	_ = DeviceRepo.Register(m)

	got := DeviceRepo.Up()
	if got != psErr.Error {
		t.Errorf("Up() = %s; want %s", got, psErr.Error)
	}
}

// Fail when it could not get device up.
func TestUp_3(t *testing.T) {
	defer reset()

	psLog.DisableOutput()
	defer psLog.EnableOutput()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := ethernet.NewMockIDevice(ctrl)
	m.EXPECT().Open().Return(psErr.Error)
	m.EXPECT().IsUp().Return(false)
	m.EXPECT().DevType().Return(ethernet.DevTypeEthernet).AnyTimes()
	m.EXPECT().DevName().Return("net0").AnyTimes()
	m.EXPECT().PrivDevName().Return("tap0").AnyTimes()
	m.EXPECT().EthAddr().Return(ethernet.EthAddr{})
	m.EXPECT().BroadcastEthAddr().Return(ethernet.EthAddr{})
	m.EXPECT().PeerEthAddr().Return(ethernet.EthAddr{})

	_ = DeviceRepo.Register(m)

	got := DeviceRepo.Up()
	if got != psErr.Error {
		t.Errorf("Up() = %s; want %s", got, psErr.Error)
	}
}
