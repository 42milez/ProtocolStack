package repo

import (
	psErr "github.com/42milez/ProtocolStack/src/error"
	psLog "github.com/42milez/ProtocolStack/src/log"
	"github.com/42milez/ProtocolStack/src/mw"
	"github.com/42milez/ProtocolStack/src/net/eth"
	"github.com/golang/mock/gomock"
	"testing"
)

func TestDeviceRepo_NextNumber(t *testing.T) {
	want := 0
	got := DeviceRepo.NextNumber()
	if got != want {
		t.Errorf("DeviceRepo.NextNumber() = %d; want %d", got, want)
	}
}

func TestDeviceRepo_Poll_1(t *testing.T) {
	ctrl, teardown := setup(t)
	defer teardown()

	m := mw.NewMockIDevice(ctrl)
	m.EXPECT().IsUp().Return(true)
	m.EXPECT().Poll().Return(psErr.OK)
	m.EXPECT().Type().Return(mw.EthernetDevice)
	m.EXPECT().Name().Return("net0")
	m.EXPECT().Priv().Return(mw.Privilege{Name: "tap0"})
	m.EXPECT().Addr().Return(mw.EthAddr{})
	_ = DeviceRepo.Register(m)

	want := psErr.OK
	got := DeviceRepo.Poll()
	if got != want {
		t.Errorf("DeviceRepo.Poll() = %s; want %s", got, want)
	}
}

// return OK when device is down
func TestDeviceRepo_Poll_2(t *testing.T) {
	ctrl, teardown := setup(t)
	defer teardown()

	m := mw.NewMockIDevice(ctrl)
	m.EXPECT().IsUp().Return(false)
	m.EXPECT().Type().Return(mw.EthernetDevice)
	m.EXPECT().Name().Return("net0")
	m.EXPECT().Priv().Return(mw.Privilege{Name: "tap0"})
	m.EXPECT().Addr().Return(mw.EthAddr{})
	_ = DeviceRepo.Register(m)

	want := psErr.OK
	got := DeviceRepo.Poll()
	if got != want {
		t.Errorf("DeviceRepo.Poll() = %s; want %s", got, want)
	}
}

// DeviceRepo.Poll() returns OK when IDevice.Poll() returns Interrupted
func TestDeviceRepo_Poll_3(t *testing.T) {
	ctrl, teardown := setup(t)
	defer teardown()

	m := mw.NewMockIDevice(ctrl)
	m.EXPECT().IsUp().Return(true)
	m.EXPECT().Poll().Return(psErr.Interrupted)
	m.EXPECT().Type().Return(mw.EthernetDevice)
	m.EXPECT().Name().Return("net0")
	m.EXPECT().Priv().Return(mw.Privilege{Name: "tap0"})
	m.EXPECT().Addr().Return(mw.EthAddr{})
	_ = DeviceRepo.Register(m)

	want := psErr.OK
	got := DeviceRepo.Poll()
	if got != want {
		t.Errorf("DeviceRepo.Poll() = %s; want %s", got, want)
	}
}

// DeviceRepo.Poll() returns Error when IDevice.Poll() returns Error
func TestDeviceRepo_Poll_4(t *testing.T) {
	ctrl, teardown := setup(t)
	defer teardown()

	m := mw.NewMockIDevice(ctrl)
	m.EXPECT().IsUp().Return(true)
	m.EXPECT().Poll().Return(psErr.Error)
	m.EXPECT().Type().Return(mw.EthernetDevice)
	m.EXPECT().Name().Return("net0")
	m.EXPECT().Priv().Return(mw.Privilege{Name: "tap0"})
	m.EXPECT().Addr().Return(mw.EthAddr{})
	_ = DeviceRepo.Register(m)

	want := psErr.Error
	got := DeviceRepo.Poll()
	if got != want {
		t.Errorf("DeviceRepo.Poll() = %s; want %s", got, want)
	}
}

func TestDeviceRepo_Register_1(t *testing.T) {
	_, teardown := setup(t)
	defer teardown()

	dev := &eth.TapDevice{}

	got := DeviceRepo.Register(dev)
	if got != psErr.OK {
		t.Errorf("DeviceRepo.Register() = %s; want %s", got, psErr.OK)
	}
}

// Fail when it's trying to register same device.
func TestDeviceRepo_Register_2(t *testing.T) {
	_, teardown := setup(t)
	defer teardown()

	dev1 := &eth.TapDevice{Device: mw.Device{Name_: "net0"}}
	dev2 := &eth.TapDevice{Device: mw.Device{Name_: "net0"}}

	_ = DeviceRepo.Register(dev1)
	got := DeviceRepo.Register(dev2)
	if got != psErr.Error {
		t.Errorf("DeviceRepo.Register() = %s; want %s", got, psErr.Error)
	}
}

func TestUp_1(t *testing.T) {
	ctrl, teardown := setup(t)
	defer teardown()

	m := mw.NewMockIDevice(ctrl)
	m.EXPECT().Open().Return(psErr.OK)
	m.EXPECT().Up()
	m.EXPECT().IsUp().Return(false)
	m.EXPECT().Addr().Return(mw.EthAddr{})
	m.EXPECT().Type().Return(mw.EthernetDevice).AnyTimes()
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
	ctrl, teardown := setup(t)
	defer teardown()

	m := mw.NewMockIDevice(ctrl)
	m.EXPECT().IsUp().Return(true)
	m.EXPECT().Type().Return(mw.EthernetDevice).AnyTimes()
	m.EXPECT().Name().Return("net0").AnyTimes()
	m.EXPECT().Priv().Return(mw.Privilege{Name: "tap0"}).AnyTimes()
	m.EXPECT().Addr().Return(mw.EthAddr{})

	_ = DeviceRepo.Register(m)

	got := DeviceRepo.Up()
	if got != psErr.Error {
		t.Errorf("Up() = %s; want %s", got, psErr.Error)
	}
}

// Fail when it could not get device up.
func TestUp_3(t *testing.T) {
	ctrl, teardown := setup(t)
	defer teardown()

	m := mw.NewMockIDevice(ctrl)
	m.EXPECT().Open().Return(psErr.Error)
	m.EXPECT().IsUp().Return(false)
	m.EXPECT().Type().Return(mw.EthernetDevice).AnyTimes()
	m.EXPECT().Name().Return("net0").AnyTimes()
	m.EXPECT().Priv().Return(mw.Privilege{Name: "tap0"}).AnyTimes()
	m.EXPECT().Addr().Return(mw.EthAddr{})

	_ = DeviceRepo.Register(m)

	got := DeviceRepo.Up()
	if got != psErr.Error {
		t.Errorf("Up() = %s; want %s", got, psErr.Error)
	}
}

func TestIfaceRepo_Get_1(t *testing.T) {
	_, teardown := setup(t)
	defer teardown()

	iface := &mw.Iface{
		Family:    mw.V4AddrFamily,
		Unicast:   mw.ParseIP(mw.LoopbackIpAddr),
		Netmask:   mw.ParseIP(mw.LoopbackNetmask),
		Broadcast: make(mw.IP, 0),
	}
	dev := &eth.TapDevice{
		Device: mw.Device{
			Type_: mw.EthernetDevice,
			MTU_:  mw.EthPayloadLenMax,
			Flag_: mw.BroadcastFlag | mw.NeedArpFlag,
			Addr_: mw.EthAddr{11, 12, 13, 14, 15, 16},
			Priv_: mw.Privilege{FD: -1, Name: "tap0"},
		},
	}
	_ = IfaceRepo.Register(iface, dev)

	if IfaceRepo.Get(iface.Unicast) != iface {
		t.Errorf("IfaceRepo.Get() returns invalid Iface")
	}
}

func TestIfaceRepo_Get_2(t *testing.T) {
	_, teardown := setup(t)
	defer teardown()

	if IfaceRepo.Get(mw.ParseIP(mw.LoopbackIpAddr)) != nil {
		t.Errorf("IfaceRepo.Get() returns invalid Iface")
	}
}

func TestIfaceRepo_Lookup_1(t *testing.T) {
	_, teardown := setup(t)
	defer teardown()

	iface := &mw.Iface{
		Family:    mw.V4AddrFamily,
		Unicast:   mw.ParseIP(mw.LoopbackIpAddr),
		Netmask:   mw.ParseIP(mw.LoopbackNetmask),
		Broadcast: make(mw.IP, 0),
	}
	dev := &eth.TapDevice{
		Device: mw.Device{
			Type_: mw.EthernetDevice,
			MTU_:  mw.EthPayloadLenMax,
			Flag_: mw.BroadcastFlag | mw.NeedArpFlag,
			Addr_: mw.EthAddr{11, 12, 13, 14, 15, 16},
			Priv_: mw.Privilege{FD: -1, Name: "tap0"},
		},
	}
	_ = IfaceRepo.Register(iface, dev)

	if IfaceRepo.Lookup(dev, mw.V4AddrFamily) != iface {
		t.Errorf("IfaceRepo.Lookup() returns invalid Iface")
	}
}

func TestIfaceRepo_Lookup_2(t *testing.T) {
	_, teardown := setup(t)
	defer teardown()

	dev := &eth.TapDevice{
		Device: mw.Device{
			Type_: mw.EthernetDevice,
			MTU_:  mw.EthPayloadLenMax,
			Flag_: mw.BroadcastFlag | mw.NeedArpFlag,
			Addr_: mw.EthAddr{11, 12, 13, 14, 15, 16},
			Priv_: mw.Privilege{FD: -1, Name: "tap0"},
		},
	}

	if IfaceRepo.Lookup(dev, mw.V4AddrFamily) != nil {
		t.Errorf("IfaceRepo.Lookup() returns invalid Iface")
	}
}

func TestIfaceRepo_Register_1(t *testing.T) {
	_, teardown := setup(t)
	defer teardown()

	iface := &mw.Iface{
		Family:    mw.V4AddrFamily,
		Unicast:   mw.ParseIP(mw.LoopbackIpAddr),
		Netmask:   mw.ParseIP(mw.LoopbackNetmask),
		Broadcast: make(mw.IP, 0),
	}

	dev := &eth.TapDevice{
		Device: mw.Device{
			Type_: mw.EthernetDevice,
			MTU_:  mw.EthPayloadLenMax,
			Flag_: mw.BroadcastFlag | mw.NeedArpFlag,
			Addr_: mw.EthAddr{11, 12, 13, 14, 15, 16},
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
	_, teardown := setup(t)
	defer teardown()

	iface := &mw.Iface{
		Family:    mw.V4AddrFamily,
		Unicast:   mw.ParseIP(mw.LoopbackIpAddr),
		Netmask:   mw.ParseIP(mw.LoopbackNetmask),
		Broadcast: make(mw.IP, 0),
	}

	dev := &eth.TapDevice{
		Device: mw.Device{
			Type_: mw.EthernetDevice,
			MTU_:  mw.EthPayloadLenMax,
			Flag_: mw.BroadcastFlag | mw.NeedArpFlag,
			Addr_: mw.EthAddr{11, 12, 13, 14, 15, 16},
			Priv_: mw.Privilege{FD: -1, Name: "tap0"},
		},
	}

	_ = IfaceRepo.Register(iface, dev)

	got := IfaceRepo.Register(iface, dev)
	if got != psErr.Error {
		t.Errorf("IfaceRepo.Register() = %s; want %s", got, psErr.Error)
	}
}

func TestRouteRepo_Get_1(t *testing.T) {
	_, teardown := setup(t)
	defer teardown()

	iface := &mw.Iface{
		Family:    mw.V4AddrFamily,
		Unicast:   mw.IP{192, 0, 0, 1},
		Netmask:   mw.IP{255, 255, 255, 0},
		Broadcast: mw.IP{192, 0, 0, 255},
	}
	dev := &eth.TapDevice{
		Device: mw.Device{
			Type_: mw.EthernetDevice,
			MTU_:  mw.EthPayloadLenMax,
			Flag_: mw.BroadcastFlag | mw.NeedArpFlag,
			Addr_: mw.EthAddr{11, 12, 13, 14, 15, 16},
			Priv_: mw.Privilege{FD: -1, Name: "tap0"},
		},
	}
	_ = IfaceRepo.Register(iface, dev)
	RouteRepo.Register(mw.ParseIP("192.0.0.0"), mw.V4Any, iface)

	// return valid route
	if RouteRepo.Get(mw.IP{192, 0, 0, 1}) == nil {
		t.Errorf("RouteRepo.Get() return invalid route")
	}

	// return nil (route not exist)
	if RouteRepo.Get(mw.IP{192, 0, 2, 1}) != nil {
		t.Errorf("RouteRepo.Get() return invalid route")
	}
}

func TestRouteRepo_Get_2(t *testing.T) {
	_, teardown := setup(t)
	defer teardown()

	iface := &mw.Iface{
		Family:    mw.V4AddrFamily,
		Unicast:   mw.IP{192, 0, 2, 2},
		Netmask:   mw.IP{255, 255, 255, 0},
		Broadcast: mw.IP{192, 0, 0, 255},
	}
	dev := &eth.TapDevice{
		Device: mw.Device{
			Type_: mw.EthernetDevice,
			MTU_:  mw.EthPayloadLenMax,
			Flag_: mw.BroadcastFlag | mw.NeedArpFlag,
			Addr_: mw.EthAddr{11, 12, 13, 14, 15, 16},
			Priv_: mw.Privilege{FD: -1, Name: "tap0"},
		},
	}
	_ = IfaceRepo.Register(iface, dev)
	RouteRepo.RegisterDefaultGateway(iface, mw.ParseIP("192.0.2.1"))

	// return valid route
	if RouteRepo.Get(mw.IP{192, 0, 2, 1}) == nil {
		t.Errorf("RouteRepo.Get() return invalid route")
	}
}

func setup(t *testing.T) (ctrl *gomock.Controller, teardown func()) {
	ctrl = gomock.NewController(t)
	psLog.DisableOutput()
	reset := func() {
		psLog.EnableOutput()
		DeviceRepo = &deviceRepo{}
		IfaceRepo = &ifaceRepo{}
	}
	teardown = func() {
		ctrl.Finish()
		reset()
	}
	return
}
