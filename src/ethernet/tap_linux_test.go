package ethernet

import (
	"errors"
	psErr "github.com/42milez/ProtocolStack/src/error"
	psLog "github.com/42milez/ProtocolStack/src/log"
	psSyscall "github.com/42milez/ProtocolStack/src/syscall"
	"github.com/golang/mock/gomock"
	"syscall"
	"testing"
)

func TestTapDevice_Open_SUCCESS(t *testing.T) {
	psLog.DisableOutput()
	defer psLog.EnableOutput()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := psSyscall.NewMockISyscall(ctrl)
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

func TestTapDevice_Open_FAIL_WhenOpenSyscallFailed(t *testing.T) {
	psLog.DisableOutput()
	defer psLog.EnableOutput()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := psSyscall.NewMockISyscall(ctrl)
	m.EXPECT().Open(gomock.Any(), gomock.Any(), gomock.Any()).Return(-1, errors.New(""))

	tapDev := TapDevice{Device{Syscall: m}}

	got := tapDev.Open()
	if got.Code != psErr.CantOpen {
		t.Errorf("TapDevice.Open() = %v; want %v", got.Code, psErr.CantOpen)
	}
}

func TestTapDevice_Open_FAIL_WhenIoctlSyscallFailed_A(t *testing.T) {
	psLog.DisableOutput()
	defer psLog.EnableOutput()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ioctlRetVal := -1
	m := psSyscall.NewMockISyscall(ctrl)
	m.EXPECT().Open(gomock.Any(), gomock.Any(), gomock.Any()).Return(10, nil)
	m.EXPECT().Ioctl(gomock.Any(), gomock.Any(), gomock.Any()).Return(uintptr(ioctlRetVal), uintptr(0), syscall.EBADF)
	m.EXPECT().Close(gomock.Any()).Return(nil)

	tapDev := TapDevice{Device{Syscall: m}}

	got := tapDev.Open()
	if got.Code != psErr.CantOpen {
		t.Errorf("TapDevice.Open() = %v; want %v", got.Code, psErr.CantOpen)
	}
}

func TestTapDevice_Open_FAIL_WhenSocketSyscallFailed(t *testing.T) {
	psLog.DisableOutput()
	defer psLog.EnableOutput()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := psSyscall.NewMockISyscall(ctrl)
	m.EXPECT().Open(gomock.Any(), gomock.Any(), gomock.Any()).Return(10, nil)
	m.EXPECT().Ioctl(gomock.Any(), uintptr(syscall.TUNSETIFF), gomock.Any()).Return(uintptr(0), uintptr(0), syscall.Errno(0))
	m.EXPECT().Socket(gomock.Any(), gomock.Any(), gomock.Any()).Return(-1, errors.New(""))

	tapDev := TapDevice{Device{Syscall: m}}

	got := tapDev.Open()
	if got.Code != psErr.CantOpen {
		t.Errorf("TapDevice.Open() = %v; want %v", got.Code, psErr.CantOpen)
	}
}

func TestTapDevice_Open_FAIL_WhenIoctlSyscallFailed_B(t *testing.T) {
	psLog.DisableOutput()
	defer psLog.EnableOutput()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ioctlRetVal := -1
	m := psSyscall.NewMockISyscall(ctrl)
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

func TestTapDevice_Open_FAIL_WhenEpollCreate1SyscallFailed(t *testing.T) {
	psLog.DisableOutput()
	defer psLog.EnableOutput()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := psSyscall.NewMockISyscall(ctrl)
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

func TestTapDevice_Open_FAIL_WhenEpollCtlSyscallFailed(t *testing.T) {
	psLog.DisableOutput()
	defer psLog.EnableOutput()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := psSyscall.NewMockISyscall(ctrl)
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

func TestTapDevice_Close_SUCCESS(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := psSyscall.NewMockISyscall(ctrl)
	m.EXPECT().Close(gomock.Any()).Return(nil)

	tapDev := TapDevice{Device{Syscall: m}}

	got := tapDev.Close()
	if got.Code != psErr.OK {
		t.Errorf("TapDevice.Close() = %v; want %v", got.Code, psErr.OK)
	}
}

func TestTapDevice_Transmit_SUCCESS(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tapDev := TapDevice{Device{Syscall: psSyscall.NewMockISyscall(ctrl)}}

	got := tapDev.Transmit(EthAddr{}, make([]byte, 0), EthTypeArp)
	if got.Code != psErr.OK {
		t.Errorf("TapDevice.Transmit() = %v; want %v", got.Code, psErr.OK)
	}
}

func TestTapDevice_Poll_SUCCESS_WhenNoEventOccurred(t *testing.T) {
	psLog.DisableOutput()
	defer psLog.EnableOutput()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := psSyscall.NewMockISyscall(ctrl)
	m.EXPECT().EpollWait(gomock.Any(), gomock.Any(), gomock.Any()).Return(0, nil)

	tapDev := TapDevice{Device{Syscall: m}}

	got := tapDev.Poll(false)
	if got.Code != psErr.OK {
		t.Errorf("TapDevice.Poll() = %v; want %v", got.Code, psErr.OK)
	}
}

func TestTapDevice_Poll_SUCCESS_WhenEventOccurred(t *testing.T) {
	psLog.DisableOutput()
	defer psLog.EnableOutput()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	nEvents := 1
	m := psSyscall.NewMockISyscall(ctrl)
	m.EXPECT().EpollWait(gomock.Any(), gomock.Any(), gomock.Any()).Return(nEvents, nil)
	m.EXPECT().Read(gomock.Any(), gomock.Any()).Return(150, nil)

	tapDev := TapDevice{Device{Syscall: m}}

	got := tapDev.Poll(false)
	if got.Code != psErr.OK {
		t.Errorf("TapDevice.Poll() = %v; want %v", got.Code, psErr.OK)
	}
}

func TestTapDevice_Poll_SUCCESS_WhenTerminated(t *testing.T) {
	psLog.DisableOutput()
	defer psLog.EnableOutput()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := psSyscall.NewMockISyscall(ctrl)
	m.EXPECT().Close(gomock.Any()).Return(nil)

	tapDev := TapDevice{Device{Syscall: m}}

	got := tapDev.Poll(true)
	if got.Code != psErr.OK {
		t.Errorf("TapDevice.Poll() = %v; want %v", got.Code, psErr.OK)
	}
}

func TestTapDevice_Poll_FAIL_WhenEpollWaitSyscallFailed(t *testing.T) {
	psLog.DisableOutput()
	defer psLog.EnableOutput()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := psSyscall.NewMockISyscall(ctrl)
	m.EXPECT().EpollWait(gomock.Any(), gomock.Any(), gomock.Any()).Return(-1, errors.New(""))
	m.EXPECT().Close(gomock.Any()).Return(nil)

	tapDev := TapDevice{Device{Syscall: m}}

	got := tapDev.Poll(false)
	if got.Code != psErr.Interrupted {
		t.Errorf("TapDevice.Poll() = %v; want %v", got.Code, psErr.Interrupted)
	}
}
