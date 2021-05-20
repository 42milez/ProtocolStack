package ethernet

import (
	"errors"
	psErr "github.com/42milez/ProtocolStack/src/error"
	psLog "github.com/42milez/ProtocolStack/src/log"
	mockSyscall "github.com/42milez/ProtocolStack/src/mock/syscall"
	"github.com/golang/mock/gomock"
	"syscall"
	"testing"
)

func TestTapDevice_Open_A(t *testing.T) {
	psLog.DisableOutput()
	defer psLog.EnableOutput()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mockSyscall.NewMockISyscall(ctrl)
	m.EXPECT().Open(gomock.Any(), gomock.Any(), gomock.Any()).Return(10, nil)
	m.EXPECT().Ioctl(gomock.Any(), uintptr(syscall.TUNSETIFF), gomock.Any()).Return(uintptr(0), uintptr(0), syscall.Errno(0))
	m.EXPECT().Socket(gomock.Any(), gomock.Any(), gomock.Any()).Return(11, nil)
	m.EXPECT().Ioctl(gomock.Any(), uintptr(syscall.SIOCGIFHWADDR), gomock.Any()).Return(uintptr(0), uintptr(0), syscall.Errno(0))
	m.EXPECT().Close(gomock.Any()).Return(nil)
	m.EXPECT().EpollCreate1(gomock.Any()).Return(12, nil)
	m.EXPECT().EpollCtl(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	tapDev := TapDevice{Device{Syscall: m}}

	got := tapDev.Open()
	if got.Code != psErr.OK {
		t.Errorf("TapDevice.Open() = %v; want %v", got.Code, psErr.OK)
	}
}

func TestTapDevice_Open_B(t *testing.T) {
	psLog.DisableOutput()
	defer psLog.EnableOutput()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mockSyscall.NewMockISyscall(ctrl)
	m.EXPECT().Open(gomock.Any(), gomock.Any(), gomock.Any()).Return(-1, errors.New(""))

	tapDev := TapDevice{Device{Syscall: m}}

	got := tapDev.Open()
	if got.Code != psErr.CantOpen {
		t.Errorf("TapDevice.Open() = %v; want %v", got.Code, psErr.CantOpen)
	}
}

func TestTapDevice_Open_C(t *testing.T) {
	psLog.DisableOutput()
	defer psLog.EnableOutput()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ioctlRetVal := -1
	m := mockSyscall.NewMockISyscall(ctrl)
	m.EXPECT().Open(gomock.Any(), gomock.Any(), gomock.Any()).Return(10, nil)
	m.EXPECT().Ioctl(gomock.Any(), gomock.Any(), gomock.Any()).Return(uintptr(ioctlRetVal), uintptr(0), syscall.EBADF)
	m.EXPECT().Close(gomock.Any()).Return(nil)

	tapDev := TapDevice{Device{Syscall: m}}

	got := tapDev.Open()
	if got.Code != psErr.CantOpen {
		t.Errorf("TapDevice.Open() = %v; want %v", got.Code, psErr.CantOpen)
	}
}

func TestTapDevice_Open_D(t *testing.T) {
	psLog.DisableOutput()
	defer psLog.EnableOutput()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mockSyscall.NewMockISyscall(ctrl)
	m.EXPECT().Open(gomock.Any(), gomock.Any(), gomock.Any()).Return(10, nil)
	m.EXPECT().Ioctl(gomock.Any(), uintptr(syscall.TUNSETIFF), gomock.Any()).Return(uintptr(0), uintptr(0), syscall.Errno(0))
	m.EXPECT().Socket(gomock.Any(), gomock.Any(), gomock.Any()).Return(-1, errors.New(""))

	tapDev := TapDevice{Device{Syscall: m}}

	got := tapDev.Open()
	if got.Code != psErr.CantOpen {
		t.Errorf("TapDevice.Open() = %v; want %v", got.Code, psErr.CantOpen)
	}
}

func TestTapDevice_Open_E(t *testing.T) {
	psLog.DisableOutput()
	defer psLog.EnableOutput()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ioctlRetVal := -1
	m := mockSyscall.NewMockISyscall(ctrl)
	m.EXPECT().Open(gomock.Any(), gomock.Any(), gomock.Any()).Return(10, nil)
	m.EXPECT().Ioctl(gomock.Any(), uintptr(syscall.TUNSETIFF), gomock.Any()).Return(uintptr(0), uintptr(0), syscall.Errno(0))
	m.EXPECT().Socket(gomock.Any(), gomock.Any(), gomock.Any()).Return(11, nil)
	m.EXPECT().Ioctl(gomock.Any(), uintptr(syscall.SIOCGIFHWADDR), gomock.Any()).Return(uintptr(ioctlRetVal), uintptr(0), syscall.EBADF)
	m.EXPECT().Close(gomock.Any()).Return(nil)

	tapDev := TapDevice{Device{Syscall: m}}

	got := tapDev.Open()
	if got.Code != psErr.CantOpen {
		t.Errorf("TapDevice.Open() = %v; want %v", got.Code, psErr.CantOpen)
	}
}

func TestTapDevice_Open_F(t *testing.T) {
	psLog.DisableOutput()
	defer psLog.EnableOutput()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mockSyscall.NewMockISyscall(ctrl)
	m.EXPECT().Open(gomock.Any(), gomock.Any(), gomock.Any()).Return(10, nil)
	m.EXPECT().Ioctl(gomock.Any(), uintptr(syscall.TUNSETIFF), gomock.Any()).Return(uintptr(0), uintptr(0), syscall.Errno(0))
	m.EXPECT().Socket(gomock.Any(), gomock.Any(), gomock.Any()).Return(11, nil)
	m.EXPECT().Ioctl(gomock.Any(), uintptr(syscall.SIOCGIFHWADDR), gomock.Any()).Return(uintptr(0), uintptr(0), syscall.Errno(0))
	m.EXPECT().Close(gomock.Any()).Return(nil)
	m.EXPECT().EpollCreate1(gomock.Any()).Return(-1, errors.New(""))

	tapDev := TapDevice{Device{Syscall: m}}

	got := tapDev.Open()
	if got.Code != psErr.CantOpen {
		t.Errorf("TapDevice.Open() = %v; want %v", got.Code, psErr.CantOpen)
	}
}

func TestTapDevice_Open_G(t *testing.T) {
	psLog.DisableOutput()
	defer psLog.EnableOutput()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mockSyscall.NewMockISyscall(ctrl)
	m.EXPECT().Open(gomock.Any(), gomock.Any(), gomock.Any()).Return(10, nil)
	m.EXPECT().Ioctl(gomock.Any(), uintptr(syscall.TUNSETIFF), gomock.Any()).Return(uintptr(0), uintptr(0), syscall.Errno(0))
	m.EXPECT().Socket(gomock.Any(), gomock.Any(), gomock.Any()).Return(11, nil)
	m.EXPECT().Ioctl(gomock.Any(), uintptr(syscall.SIOCGIFHWADDR), gomock.Any()).Return(uintptr(0), uintptr(0), syscall.Errno(0))
	m.EXPECT().Close(gomock.Any()).Return(nil)
	m.EXPECT().EpollCreate1(gomock.Any()).Return(12, nil)
	m.EXPECT().EpollCtl(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New(""))

	tapDev := TapDevice{Device{Syscall: m}}

	got := tapDev.Open()
	if got.Code != psErr.CantOpen {
		t.Errorf("TapDevice.Open() = %v; want %v", got.Code, psErr.CantOpen)
	}
}

func TestTapDevice_Close(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mockSyscall.NewMockISyscall(ctrl)
	m.EXPECT().Close(gomock.Any()).Return(nil)

	tapDev := TapDevice{Device{Syscall: m}}

	got := tapDev.Close()
	if got.Code != psErr.OK {
		t.Errorf("TapDevice.Close() = %v; want %v", got.Code, psErr.OK)
	}
}

func TestTapDevice_Transmit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tapDev := TapDevice{Device{Syscall: mockSyscall.NewMockISyscall(ctrl)}}

	got := tapDev.Transmit()
	if got.Code != psErr.OK {
		t.Errorf("TapDevice.Transmit() = %v; want %v", got.Code, psErr.OK)
	}
}

func TestTapDevice_Poll_NoEventOccurred(t *testing.T) {
	psLog.DisableOutput()
	defer psLog.EnableOutput()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mockSyscall.NewMockISyscall(ctrl)
	m.EXPECT().EpollWait(gomock.Any(), gomock.Any(), gomock.Any()).Return(0, nil)

	tapDev := TapDevice{Device{Syscall: mockSyscall.NewMockISyscall(ctrl)}}

	got := tapDev.Poll(false)
	if got.Code != psErr.OK {
		t.Errorf("TapDevice.Poll() = %v; want %v", got.Code, psErr.OK)
	}
}

func TestTapDevice_Poll_EventOccurred(t *testing.T) {
	psLog.DisableOutput()
	defer psLog.EnableOutput()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	nEvents := 1
	ethHdrLen := EthHeaderSize * 8
	m := mockSyscall.NewMockISyscall(ctrl)
	m.EXPECT().EpollWait(gomock.Any(), gomock.Any(), gomock.Any()).Return(nEvents, nil)
	m.EXPECT().
		Read(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(uintptr(ethHdrLen), uintptr(0), syscall.Errno(0))

	tapDev := TapDevice{Device{Syscall: m}}

	got := tapDev.Poll(false)
	if got.Code != psErr.OK {
		t.Errorf("TapDevice.Poll() = %v; want %v", got.Code, psErr.OK)
	}
}

func TestTapDevice_Poll_Terminated(t *testing.T) {
	psLog.DisableOutput()
	defer psLog.EnableOutput()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mockSyscall.NewMockISyscall(ctrl)
	m.EXPECT().Close(gomock.Any()).Return(nil)

	tapDev := TapDevice{Device{Syscall: m}}

	got := tapDev.Poll(true)
	if got.Code != psErr.OK {
		t.Errorf("TapDevice.Poll() = %v; want %v", got.Code, psErr.OK)
	}
}

func TestTapDevice_Poll_FailOnEpollWait(t *testing.T) {
	psLog.DisableOutput()
	defer psLog.EnableOutput()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mockSyscall.NewMockISyscall(ctrl)
	m.EXPECT().EpollWait(gomock.Any(), gomock.Any(), gomock.Any()).Return(-1, errors.New(""))
	m.EXPECT().Close(gomock.Any()).Return(nil)

	tapDev := TapDevice{Device{Syscall: m}}

	got := tapDev.Poll(false)
	if got.Code != psErr.Interrupted {
		t.Errorf("TapDevice.Poll() = %v; want %v", got.Code, psErr.Interrupted)
	}
}
