package ethernet

import (
	"github.com/golang/mock/gomock"
	"testing"
)

func TestDevType_String_SUCCESS_A(t *testing.T) {
	devType := DevTypeEthernet
	want := "DEVICE_TYPE_ETHERNET"
	got := devType.String()
	if got != want {
		t.Errorf("DevType.String() = %v; want %v", got, want)
	}
}

func TestDevType_String_SUCCESS_B(t *testing.T) {
	devType := DevTypeLoopback
	want := "DEVICE_TYPE_LOOPBACK"
	got := devType.String()
	if got != want {
		t.Errorf("DevType.String() = %v; want %v", got, want)
	}
}

func TestDevType_String_SUCCESS_C(t *testing.T) {
	devType := DevTypeNull
	want := "DEVICE_TYPE_NULL"
	got := devType.String()
	if got != want {
		t.Errorf("DevType.String() = %v; want %v", got, want)
	}
}

func TestDevType_String_SUCCESS_D(t *testing.T) {
	devType := DevType(99)
	want := "UNKNOWN"
	got := devType.String()
	if got != want {
		t.Errorf("DevType.String() = %v; want %v", got, want)
	}
}

func TestDevice_Enable_SUCCESS(t *testing.T) {
	dev := Device{}
	dev.Up()

	want := DevFlagUp
	got := dev.FLAG & DevFlagUp

	if got != want {
		t.Errorf("Device.Up() = %v; want %v", got, want)
	}
}

func TestDevice_Disable_SUCCESS(t *testing.T) {
	dev := Device{}
	dev.Up()
	dev.Down()

	want := DevFlag(0)
	got := dev.FLAG & DevFlagUp

	if got != want {
		t.Errorf("Device.Down() = %v; want %v", got, want)
	}
}

func TestDevice_Equal_SUCCESS(t *testing.T) {
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

func TestDevice_Info_SUCCESS(t *testing.T) {
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

func TestDevice_IsUp_SUCCESS_WhenDeviceIsUp(t *testing.T) {
	dev := &Device{}
	dev.Up()
	if got := dev.IsUp(); !got {
		t.Errorf("Device.IsUp() = %v; want %v", got, true)
	}
}

func TestDevice_IsUp_SUCCESS_WhenDeviceIsDown(t *testing.T) {
	dev := &Device{}
	dev.Down()
	if got := dev.IsUp(); got {
		t.Errorf("Device.IsUp() = %v; want %v", got, false)
	}
}
