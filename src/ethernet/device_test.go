package ethernet

import (
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

func TestDevice_Disable(t *testing.T) {
	dev := Device{}
	dev.Enable()
	dev.Disable()

	want := DevFlag(0)
	got := dev.FLAG&DevFlagUp

	if got != want {
		t.Errorf("Device.Disable() = %v; want %v", got, want)
	}
}

func TestDevice_Enable(t *testing.T) {
	dev := Device{}
	dev.Enable()

	want := DevFlagUp
	got := dev.FLAG&DevFlagUp

	if got != want {
		t.Errorf("Device.Enable() = %v; want %v", got, want)
	}
}

//func TestDevice_Equal(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	m := mockEthernet.NewMockIDevice(ctrl)
//	m.EXPECT().Info().Return("", "net0", "")
//
//	dev1 := &Device{Name: "net0"}
//
//	got := dev1.Equal(m)
//	if ! got {
//		t.Errorf("Device.Equal() = %v; want %v", got, true)
//	}
//}
