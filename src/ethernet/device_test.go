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

func TestDevice_Up(t *testing.T) {
	dev := Device{}
	dev.Up()

	want := DevFlagUp
	got := dev.FLAG & DevFlagUp

	if got != want {
		t.Errorf("Device.Up() = %v; want %v", got, want)
	}
}

func TestDevice_Down(t *testing.T) {
	dev := Device{}
	dev.Up()
	dev.Down()

	want := DevFlag(0)
	got := dev.FLAG & DevFlagUp

	if got != want {
		t.Errorf("Device.Down() = %v; want %v", got, want)
	}
}

func TestDevice_Equal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := NewMockIDevice(ctrl)
	m.EXPECT().Names().Return("net0", "").AnyTimes()

	dev := Device{Name: "net0"}
	got := dev.Equal(m)
	if !got {
		t.Errorf("Device.Equal() = %v; want %v", got, true)
	}

	dev = Device{Name: "net1"}
	got = dev.Equal(m)
	if !got {
		t.Errorf("Device.Equal() = %v; want %v", got, false)
	}
}

func TestDevice_IsUp(t *testing.T) {
	dev := Device{}
	dev.Up()
	if got := dev.IsUp(); !got {
		t.Errorf("Device.IsUp() = %v; want %v", got, true)
	}

	dev = Device{}
	dev.Down()
	if got := dev.IsUp(); got {
		t.Errorf("Device.IsUp() = %v; want %v", got, false)
	}
}
