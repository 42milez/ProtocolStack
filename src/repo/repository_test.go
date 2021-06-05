package repo

import (
	psErr "github.com/42milez/ProtocolStack/src/error"
	"github.com/42milez/ProtocolStack/src/eth"
	psLog "github.com/42milez/ProtocolStack/src/log"
	"github.com/42milez/ProtocolStack/src/mw"
	"github.com/golang/mock/gomock"
	"testing"
)

func reset() {
	DeviceRepo = &deviceRepo{}
	IfaceRepo = &ifaceRepo{}
}

func setupRepositoryTest(t *testing.T) (ctrl *gomock.Controller, teardown func()) {
	psLog.DisableOutput()
	ctrl = gomock.NewController(t)
	teardown = func() {
		ctrl.Finish()
		psLog.EnableOutput()
	}
	return
}

func TestDeviceRepo_Register_1(t *testing.T) {
	_, teardown := setupRepositoryTest(t)
	defer teardown()
	defer reset()

	dev := &eth.TapDevice{}

	got := DeviceRepo.Register(dev)
	if got != psErr.OK {
		t.Errorf("DeviceRepo.Register() = %s; want %s", got, psErr.OK)
	}
}

// Fail when it's trying to register same device.
func TestDeviceRepo_Register_2(t *testing.T) {
	_, teardown := setupRepositoryTest(t)
	defer teardown()
	defer reset()

	dev1 := &eth.TapDevice{Device: mw.Device{Name_: "net0"}}
	dev2 := &eth.TapDevice{Device: mw.Device{Name_: "net0"}}

	_ = DeviceRepo.Register(dev1)
	got := DeviceRepo.Register(dev2)
	if got != psErr.Error {
		t.Errorf("DeviceRepo.Register() = %s; want %s", got, psErr.Error)
	}
}

func TestIfaceRepo_Register_1(t *testing.T) {
	_, teardown := setupRepositoryTest(t)
	defer teardown()
	defer reset()

	iface := &mw.Iface{
		Family:    mw.V4AddrFamily,
		Unicast:   mw.ParseIP(mw.LoopbackIpAddr),
		Netmask:   mw.ParseIP(mw.LoopbackNetmask),
		Broadcast: make(mw.IP, 0),
	}

	dev := &eth.TapDevice{
		Device: mw.Device{
			Type_: mw.DevTypeEthernet,
			MTU_:  mw.PayloadLenMax,
			Flag_: mw.DevFlagBroadcast | mw.DevFlagNeedArp,
			Addr_: mw.Addr{11, 12, 13, 14, 15, 16},
			Priv_: mw.Privilege{FD: -1, Name: "tap0"},
		},
	}

	got := IfaceRepo.Register(iface, dev)
	if got != psErr.OK {
		t.Errorf("IfaceRepo.Register() = %v; want %v", got, psErr.OK)
	}
}

// Fail when it's trying to register same interface.
func TestIfaceRepo_Register_2(t *testing.T) {
	_, teardown := setupRepositoryTest(t)
	defer teardown()
	defer reset()

	iface := &mw.Iface{
		Family:    mw.V4AddrFamily,
		Unicast:   mw.ParseIP(mw.LoopbackIpAddr),
		Netmask:   mw.ParseIP(mw.LoopbackNetmask),
		Broadcast: make(mw.IP, 0),
	}

	dev := &eth.TapDevice{
		Device: mw.Device{
			Type_: mw.DevTypeEthernet,
			MTU_:  mw.PayloadLenMax,
			Flag_: mw.DevFlagBroadcast | mw.DevFlagNeedArp,
			Addr_: mw.Addr{11, 12, 13, 14, 15, 16},
			Priv_: mw.Privilege{FD: -1, Name: "tap0"},
		},
	}

	_ = IfaceRepo.Register(iface, dev)

	got := IfaceRepo.Register(iface, dev)
	if got != psErr.Error {
		t.Errorf("IfaceRepo.Register() = %s; want %s", got, psErr.Error)
	}
}

func TestUp_1(t *testing.T) {
	ctrl, teardown := setupRepositoryTest(t)
	defer teardown()
	defer reset()

	m := mw.NewMockIDevice(ctrl)
	m.EXPECT().Open().Return(psErr.OK)
	m.EXPECT().Up()
	m.EXPECT().IsUp().Return(false)
	m.EXPECT().Addr().Return(mw.Addr{})
	m.EXPECT().Type().Return(mw.DevTypeEthernet).AnyTimes()
	m.EXPECT().Name().Return("net0").AnyTimes()
	m.EXPECT().Priv().Return(mw.Privilege{Name: "tap0"}).AnyTimes()

	_ = DeviceRepo.Register(m)

	got := DeviceRepo.Up()
	if got != psErr.OK {
		t.Errorf("Up() = %v; want %v", got, psErr.OK)
	}
}

// Fail when device is already opened.
func TestUp_2(t *testing.T) {
	ctrl, teardown := setupRepositoryTest(t)
	defer teardown()
	defer reset()

	m := mw.NewMockIDevice(ctrl)
	m.EXPECT().IsUp().Return(true)
	m.EXPECT().Type().Return(mw.DevTypeEthernet).AnyTimes()
	m.EXPECT().Name().Return("net0").AnyTimes()
	m.EXPECT().Priv().Return(mw.Privilege{Name: "tap0"}).AnyTimes()
	m.EXPECT().Addr().Return(mw.Addr{})

	_ = DeviceRepo.Register(m)

	got := DeviceRepo.Up()
	if got != psErr.Error {
		t.Errorf("Up() = %s; want %s", got, psErr.Error)
	}
}

// Fail when it could not get device up.
func TestUp_3(t *testing.T) {
	ctrl, teardown := setupRepositoryTest(t)
	defer teardown()
	defer reset()

	m := mw.NewMockIDevice(ctrl)
	m.EXPECT().Open().Return(psErr.Error)
	m.EXPECT().IsUp().Return(false)
	m.EXPECT().Type().Return(mw.DevTypeEthernet).AnyTimes()
	m.EXPECT().Name().Return("net0").AnyTimes()
	m.EXPECT().Priv().Return(mw.Privilege{Name: "tap0"}).AnyTimes()
	m.EXPECT().Addr().Return(mw.Addr{})

	_ = DeviceRepo.Register(m)

	got := DeviceRepo.Up()
	if got != psErr.Error {
		t.Errorf("Up() = %s; want %s", got, psErr.Error)
	}
}
