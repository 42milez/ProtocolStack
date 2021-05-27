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

func TestTapDevice_Open_1(t *testing.T) {
	psLog.DisableOutput()
	defer psLog.EnableOutput()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := psSyscall.NewMockISyscall(ctrl)

	m.EXPECT().
		Open(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(10, nil)
	m.EXPECT().
		Close(gomock.Any()).Return(nil)
	m.EXPECT().
		Ioctl(gomock.Any(), uintptr(syscall.TUNSETIFF), gomock.Any()).
		Return(uintptr(0), uintptr(0), syscall.Errno(0))
	m.EXPECT().
		Ioctl(gomock.Any(), uintptr(syscall.SIOCGIFHWADDR), gomock.Any()).
		Return(uintptr(0), uintptr(0), syscall.Errno(0))
	m.EXPECT().
		Socket(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(11, nil)
	m.EXPECT().
		EpollCreate1(gomock.Any()).
		Return(12, nil)
	m.EXPECT().
		EpollCtl(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil)

	tapDev := TapDevice{Device{Syscall: m}}

	got := tapDev.Open()
	if got != psErr.OK {
		t.Errorf("TapDevice.Open() = %v; want %v", got, psErr.OK)
	}
}

// Fail when Open() returns error.
func TestTapDevice_Open_2(t *testing.T) {
	psLog.DisableOutput()
	defer psLog.EnableOutput()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := psSyscall.NewMockISyscall(ctrl)

	m.EXPECT().
		Open(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(-1, errors.New(""))

	tapDev := TapDevice{Device{Syscall: m}}

	got := tapDev.Open()
	if got != psErr.CantOpenIOResource {
		t.Errorf("TapDevice.Open() = %v; want %v", got, psErr.CantOpenIOResource)
	}
}

// Fail when Ioctl() returns error.
func TestTapDevice_Open_FAIL_WhenIoctlSyscallFailed_A(t *testing.T) {
	psLog.DisableOutput()
	defer psLog.EnableOutput()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := psSyscall.NewMockISyscall(ctrl)

	m.EXPECT().
		Open(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(10, nil)
	m.EXPECT().
		Ioctl(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(uintptr(-1), uintptr(0), syscall.EBADF)
	m.EXPECT().
		Close(gomock.Any()).
		Return(nil)

	tapDev := TapDevice{Device{Syscall: m}}

	got := tapDev.Open()
	if got != psErr.CantModifyIOResourceParameter {
		t.Errorf("TapDevice.Open() = %v; want %v", got, psErr.CantModifyIOResourceParameter)
	}
}

// Fail when Socket() returns error.
func TestTapDevice_Open_3(t *testing.T) {
	psLog.DisableOutput()
	defer psLog.EnableOutput()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := psSyscall.NewMockISyscall(ctrl)

	m.EXPECT().
		Open(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(10, nil)
	m.EXPECT().
		Ioctl(gomock.Any(), uintptr(syscall.TUNSETIFF), gomock.Any()).
		Return(uintptr(0), uintptr(0), syscall.Errno(0))
	m.EXPECT().
		Socket(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(-1, errors.New(""))
	m.EXPECT().
		Close(gomock.Any()).
		Return(nil)

	tapDev := TapDevice{Device{Syscall: m}}

	got := tapDev.Open()
	if got != psErr.CantCreateEndpoint {
		t.Errorf("TapDevice.Open() = %v; want %v", got, psErr.CantCreateEndpoint)
	}
}

// Fail when Ioctl() returns error.
func TestTapDevice_Open_4(t *testing.T) {
	psLog.DisableOutput()
	defer psLog.EnableOutput()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := psSyscall.NewMockISyscall(ctrl)

	m.EXPECT().
		Open(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(10, nil)
	m.EXPECT().
		Ioctl(gomock.Any(), uintptr(syscall.TUNSETIFF), gomock.Any()).
		Return(uintptr(0), uintptr(0), syscall.Errno(0))
	m.EXPECT().
		Socket(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(11, nil)
	m.EXPECT().
		Ioctl(gomock.Any(), uintptr(syscall.SIOCGIFHWADDR), gomock.Any()).
		Return(uintptr(-1), uintptr(0), syscall.EBADF)
	m.EXPECT().
		Close(gomock.Any()).
		Return(nil)

	tapDev := TapDevice{Device{Syscall: m}}

	got := tapDev.Open()
	if got != psErr.CantModifyIOResourceParameter {
		t.Errorf("TapDevice.Open() = %v; want %v", got, psErr.CantModifyIOResourceParameter)
	}
}

// Fail when EpollCreate1() returns error.
func TestTapDevice_Open_5(t *testing.T) {
	psLog.DisableOutput()
	defer psLog.EnableOutput()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := psSyscall.NewMockISyscall(ctrl)

	m.EXPECT().
		Open(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(10, nil)
	m.EXPECT().
		Ioctl(gomock.Any(), uintptr(syscall.TUNSETIFF), gomock.Any()).
		Return(uintptr(0), uintptr(0), syscall.Errno(0))
	m.EXPECT().
		Socket(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(11, nil)
	m.EXPECT().
		Ioctl(gomock.Any(), uintptr(syscall.SIOCGIFHWADDR), gomock.Any()).
		Return(uintptr(0), uintptr(0), syscall.Errno(0))
	m.EXPECT().
		EpollCreate1(gomock.Any()).
		Return(-1, errors.New(""))
	m.EXPECT().
		Close(gomock.Any()).
		Return(nil)

	tapDev := TapDevice{Device{Syscall: m}}

	got := tapDev.Open()
	if got != psErr.CantCreateEpollInstance {
		t.Errorf("TapDevice.Open() = %v; want %v", got, psErr.CantCreateEpollInstance)
	}
}

// Fail when EpollCtl() returns error.
func TestTapDevice_Open_6(t *testing.T) {
	psLog.DisableOutput()
	defer psLog.EnableOutput()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := psSyscall.NewMockISyscall(ctrl)

	m.EXPECT().
		Open(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(10, nil)
	m.EXPECT().
		Ioctl(gomock.Any(), uintptr(syscall.TUNSETIFF), gomock.Any()).
		Return(uintptr(0), uintptr(0), syscall.Errno(0))
	m.EXPECT().
		Socket(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(11, nil)
	m.EXPECT().
		Ioctl(gomock.Any(), uintptr(syscall.SIOCGIFHWADDR), gomock.Any()).
		Return(uintptr(0), uintptr(0), syscall.Errno(0))
	m.EXPECT().
		EpollCreate1(gomock.Any()).
		Return(12, nil)
	m.EXPECT().
		EpollCtl(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(errors.New(""))
	m.EXPECT().
		Close(gomock.Any()).
		Return(nil)

	tapDev := TapDevice{Device{Syscall: m}}

	got := tapDev.Open()
	if got != psErr.CantModifyIOResourceParameter {
		t.Errorf("TapDevice.Open() = %v; want %v", got, psErr.CantModifyIOResourceParameter)
	}
}

func TestTapDevice_Close(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := psSyscall.NewMockISyscall(ctrl)

	m.EXPECT().
		Close(gomock.Any()).
		Return(nil)

	tapDev := TapDevice{
		Device{
			Syscall: m,
		},
	}

	got := tapDev.Close()
	if got != psErr.OK {
		t.Errorf("TapDevice.Close() = %v; want %v", got, psErr.OK)
	}
}

func TestTapDevice_Transmit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tapDev := TapDevice{
		Device{
			Syscall: psSyscall.NewMockISyscall(ctrl),
		},
	}

	got := tapDev.Transmit(EthAddr{}, make([]byte, 0), EthTypeArp)
	if got != psErr.OK {
		t.Errorf("TapDevice.Transmit() = %v; want %v", got, psErr.OK)
	}
}

// Success when no event occurs.
func TestTapDevice_Poll_1(t *testing.T) {
	psLog.DisableOutput()
	defer psLog.EnableOutput()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := psSyscall.NewMockISyscall(ctrl)

	m.EXPECT().
		EpollWait(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(0, nil)

	tapDev := TapDevice{Device{Syscall: m}}

	got := tapDev.Poll(false)
	if got != psErr.OK {
		t.Errorf("TapDevice.Poll() = %v; want %v", got, psErr.OK)
	}
}

// Success when an event occurs.
func TestTapDevice_Poll_2(t *testing.T) {
	psLog.DisableOutput()
	defer psLog.EnableOutput()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := psSyscall.NewMockISyscall(ctrl)

	m.EXPECT().
		EpollWait(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(1, nil)
	m.EXPECT().
		Read(gomock.Any(), gomock.Any()).
		Return(150, nil)

	tapDev := TapDevice{Device{Syscall: m}}

	got := tapDev.Poll(false)
	if got != psErr.OK {
		t.Errorf("TapDevice.Poll() = %v; want %v", got, psErr.OK)
	}
}

// Success when Poll() is terminated.
func TestTapDevice_Poll_3(t *testing.T) {
	psLog.DisableOutput()
	defer psLog.EnableOutput()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := psSyscall.NewMockISyscall(ctrl)

	m.EXPECT().
		Close(gomock.Any()).
		Return(nil)

	tapDev := TapDevice{Device{Syscall: m}}

	got := tapDev.Poll(true)
	if got != psErr.OK {
		t.Errorf("TapDevice.Poll() = %v; want %v", got, psErr.OK)
	}
}

// Fail when EpollWait() is interrupted.
func TestTapDevice_Poll_4(t *testing.T) {
	psLog.DisableOutput()
	defer psLog.EnableOutput()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := psSyscall.NewMockISyscall(ctrl)
	m.EXPECT().
		EpollWait(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(-1, errors.New(""))
	m.EXPECT().
		Close(gomock.Any()).
		Return(nil)

	tapDev := TapDevice{Device{Syscall: m}}

	got := tapDev.Poll(false)
	if got != psErr.Interrupted {
		t.Errorf("TapDevice.Poll() = %v; want %v", got, psErr.Interrupted)
	}
}
