package ethernet

import (
	"github.com/golang/mock/gomock"
	"testing"
)

func TestDevType_String(t *testing.T) {
	devType := DevTypeEthernet
	want := "ETHERNET"
	got := devType.String()
	if got != want {
		t.Errorf("DevType.String() = %v; want %v", got, want)
	}

	devType = DevTypeLoopback
	want = "LOOPBACK"
	got = devType.String()
	if got != want {
		t.Errorf("DevType.String() = %v; want %v", got, want)
	}

	devType = DevTypeNull
	want = "NULL"
	got = devType.String()
	if got != want {
		t.Errorf("DevType.String() = %v; want %v", got, want)
	}

	devType = DevType(99)
	want = "UNKNOWN"
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

func TestDevice_AddrLen(t *testing.T) {
	dev := Device{
		AddrLen_: 0,
	}
	want := uint16(0)
	got := dev.AddrLen()
	if got != want {
		t.Errorf("Name() = %d; want %d", got, want)
	}
}

func TestDevice_HdrLen(t *testing.T) {
	dev := Device{
		HdrLen_: EthFrameSizeMin,
	}
	want := uint16(EthFrameSizeMin)
	got := dev.HdrLen()
	if got != want {
		t.Errorf("Name() = %d; want %d", got, want)
	}
}
