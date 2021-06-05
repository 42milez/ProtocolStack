package eth

import (
	"github.com/golang/mock/gomock"
	"testing"
)

func TestDevType_String(t *testing.T) {
	devType := DevTypeEthernet
	want := "Ethernet"
	got := devType.String()
	if got != want {
		t.Errorf("DevType.String() = %v; want %v", got, want)
	}

	devType = DevTypeLoopback
	want = "Loopback"
	got = devType.String()
	if got != want {
		t.Errorf("DevType.String() = %v; want %v", got, want)
	}

	devType = DevTypeNull
	want = "Null"
	got = devType.String()
	if got != want {
		t.Errorf("DevType.String() = %v; want %v", got, want)
	}
}

func TestDevice_Up_Down(t *testing.T) {
	dev := Device{}

	dev.Up()
	want := DevFlagUp
	got := dev.Flag()
	if got != want {
		t.Errorf("Device.Up() = %v; want %v", got, want)
	}

	dev.Down()
	want = DevFlag(0)
	got = dev.Flag()
	if got != want {
		t.Errorf("Device.Down() = %v; want %v", got, want)
	}
}

func TestDevice_Equal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := NewMockIDevice(ctrl)
	m.EXPECT().Name().Return("net0").AnyTimes()

	dev := Device{Name_: "net0"}
	got := dev.Equal(m)
	if !got {
		t.Errorf("Device.Equal() = %t; want %t", got, true)
	}

	dev = Device{Name_: "net1"}
	got = dev.Equal(m)
	if got {
		t.Errorf("Device.Equal() = %t; want %t", got, false)
	}
}

func TestDevice_IsUp(t *testing.T) {
	dev := Device{}
	dev.Up()
	if got := dev.IsUp(); !got {
		t.Errorf("Device.IsUp() = %t; want %t", got, true)
	}

	dev = Device{}
	dev.Down()
	if got := dev.IsUp(); got {
		t.Errorf("Device.IsUp() = %t; want %t", got, false)
	}
}

func TestDevice_Type(t *testing.T) {
	dev := Device{
		Type_: DevTypeEthernet,
	}
	want := DevTypeEthernet
	got := dev.Type()
	if got != want {
		t.Errorf("Type() = %d; want %d", got, want)
	}
}

func TestDevice_Name(t *testing.T) {
	dev := Device{
		Name_: "net0",
	}
	want := "net0"
	got := dev.Name()
	if got != want {
		t.Errorf("Name() = %s; want %s", got, want)
	}
}

func TestDevice_MTU(t *testing.T) {
	dev := Device{
		MTU_: 1514,
	}
	var want uint16 = 1514
	got := dev.MTU()
	if got != want {
		t.Errorf("MTU() = %d; want %d", got, want)
	}
}
