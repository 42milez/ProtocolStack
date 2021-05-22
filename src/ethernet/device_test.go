package ethernet

import (
	"github.com/golang/mock/gomock"
	"testing"
)

func TestDevType_String(t *testing.T) {
	devType := DevTypeEthernet
	want := "DEVICE_TYPE_ETHERNET"
	got := devType.String()
	if got != want {
		t.Errorf("DevType.String() = %v; want %v", got, want)
	}

	devType = DevTypeLoopback
	want = "DEVICE_TYPE_LOOPBACK"
	got = devType.String()
	if got != want {
		t.Errorf("DevType.String() = %v; want %v", got, want)
	}

	devType = DevTypeNull
	want = "DEVICE_TYPE_NULL"
	got = devType.String()
	if got != want {
		t.Errorf("DevType.String() = %v; want %v", got, want)
	}

	devType = 99
	want = "UNKNOWN"
	got = devType.String()
	if got != want {
		t.Errorf("DevType.String() = %v; want %v", got, want)
	}
}

func TestDevice_Enable(t *testing.T) {
	dev := Device{}
	dev.Enable()

	want := DevFlagUp
	got := dev.FLAG & DevFlagUp

	if got != want {
		t.Errorf("Device.Enable() = %v; want %v", got, want)
	}
}

func TestDevice_Disable(t *testing.T) {
	dev := Device{}
	dev.Enable()
	dev.Disable()

	want := DevFlag(0)
	got := dev.FLAG & DevFlagUp

	if got != want {
		t.Errorf("Device.Disable() = %v; want %v", got, want)
	}
}

func TestDevice_Equal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := NewMockIDevice(ctrl)
	m.EXPECT().Info().Return("", "net0", "")

	dev1 := &Device{Name: "net0"}

	got := dev1.Equal(m)
	if !got {
		t.Errorf("Device.Equal() = %v; want %v", got, true)
	}
}

func TestDevice_Info(t *testing.T) {
	devType := DevTypeEthernet
	devName1 := "net0"
	devName2 := "tap0"
	dev := &Device{
		Type: devType,
		Name: devName1,
		Priv: Privilege{
			Name: devName2,
		},
	}
	typ, n1, n2 := dev.Info()
	if typ != devType.String() || n1 != devName1 || n2 != devName2 {
		t.Errorf("Device.Equal() = %v, %v, %v; want %v, %v, %v", typ, n1, n2, devType, devName1, devName2)
	}
}

func TestDevice_IsUp(t *testing.T) {
	dev := &Device{}

	dev.Enable()
	if got := dev.IsUp(); !got {
		t.Errorf("Device.IsUp() = %v; want %v", got, true)
	}

	dev.Disable()
	if got := dev.IsUp(); got {
		t.Errorf("Device.IsUp() = %v; want %v", got, false)
	}
}
