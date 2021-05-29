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

const ErrnoSuccess = syscall.Errno(0)
const R2Zero = uintptr(0)

var RetValOnSuccess = 0
var RetValOnFail = -1

var ErrorWithNoMessage error

var Any = gomock.Any()

func TestTapDevice_Open_1(t *testing.T) {
	psLog.DisableOutput()
	defer psLog.EnableOutput()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := psSyscall.NewMockISyscall(ctrl)

	fd1 := 3
	fd2 := 4
	fd3 := 5

	m.EXPECT().Open(Any, Any, Any).Return(fd1, nil)
	m.EXPECT().Close(Any).Return(nil)
	m.EXPECT().Ioctl(Any, uintptr(syscall.TUNSETIFF), Any).Return(uintptr(RetValOnSuccess), R2Zero, ErrnoSuccess)
	m.EXPECT().Ioctl(Any, uintptr(syscall.SIOCGIFHWADDR), Any).Return(uintptr(RetValOnSuccess), R2Zero, ErrnoSuccess)
	m.EXPECT().Socket(Any, Any, Any).Return(fd2, nil)
	m.EXPECT().EpollCreate1(Any).Return(fd3, nil)
	m.EXPECT().EpollCtl(Any, Any, Any, Any).Return(nil)

	psSyscall.Syscall = m

	tapDev := TapDevice{}

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
	m.EXPECT().Open(Any, Any, Any).Return(RetValOnFail, ErrorWithNoMessage)

	psSyscall.Syscall = m

	tapDev := TapDevice{}

	got := tapDev.Open()
	if got != psErr.CantOpenIOResource {
		t.Errorf("TapDevice.Open() = %v; want %v", got, psErr.CantOpenIOResource)
	}
}

// Fail when Ioctl() returns error.
func TestTapDevice_Open_3(t *testing.T) {
	psLog.DisableOutput()
	defer psLog.EnableOutput()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := psSyscall.NewMockISyscall(ctrl)
	m.EXPECT().Open(Any, Any, Any).Return(10, nil)
	m.EXPECT().Ioctl(Any, Any, Any).Return(uintptr(RetValOnFail), R2Zero, syscall.EBADF)
	m.EXPECT().Close(Any).Return(nil)

	psSyscall.Syscall = m

	tapDev := TapDevice{}

	got := tapDev.Open()
	if got != psErr.CantModifyIOResourceParameter {
		t.Errorf("TapDevice.Open() = %v; want %v", got, psErr.CantModifyIOResourceParameter)
	}
}

// Fail when Socket() returns error.
func TestTapDevice_Open_34(t *testing.T) {
	psLog.DisableOutput()
	defer psLog.EnableOutput()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := psSyscall.NewMockISyscall(ctrl)

	fd := 3

	m.EXPECT().Open(Any, Any, Any).Return(fd, nil)
	m.EXPECT().Ioctl(Any, uintptr(syscall.TUNSETIFF), Any).Return(uintptr(RetValOnSuccess), R2Zero, ErrnoSuccess)
	m.EXPECT().Socket(Any, Any, Any).Return(RetValOnFail, ErrorWithNoMessage)
	m.EXPECT().Close(Any).Return(nil)

	psSyscall.Syscall = m

	tapDev := TapDevice{}

	got := tapDev.Open()
	if got != psErr.CantCreateEndpoint {
		t.Errorf("TapDevice.Open() = %v; want %v", got, psErr.CantCreateEndpoint)
	}
}

// Fail when Ioctl() returns error.
func TestTapDevice_Open_5(t *testing.T) {
	psLog.DisableOutput()
	defer psLog.EnableOutput()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := psSyscall.NewMockISyscall(ctrl)

	fd1 := 3
	fd2 := 4

	m.EXPECT().Open(Any, Any, Any).Return(fd1, nil)
	m.EXPECT().Ioctl(Any, uintptr(syscall.TUNSETIFF), Any).Return(uintptr(RetValOnSuccess), R2Zero, ErrnoSuccess)
	m.EXPECT().Socket(Any, Any, Any).Return(fd2, nil)
	m.EXPECT().Ioctl(Any, uintptr(syscall.SIOCGIFHWADDR), Any).Return(uintptr(RetValOnFail), R2Zero, syscall.EBADF)
	m.EXPECT().Close(Any).Return(nil)

	psSyscall.Syscall = m

	tapDev := TapDevice{}

	got := tapDev.Open()
	if got != psErr.CantModifyIOResourceParameter {
		t.Errorf("TapDevice.Open() = %v; want %v", got, psErr.CantModifyIOResourceParameter)
	}
}

// Fail when EpollCreate1() returns error.
func TestTapDevice_Open_6(t *testing.T) {
	psLog.DisableOutput()
	defer psLog.EnableOutput()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := psSyscall.NewMockISyscall(ctrl)

	fd1 := 3
	fd2 := 4

	m.EXPECT().Open(Any, Any, Any).Return(fd1, nil)
	m.EXPECT().Ioctl(Any, uintptr(syscall.TUNSETIFF), Any).Return(uintptr(RetValOnSuccess), R2Zero, ErrnoSuccess)
	m.EXPECT().Socket(Any, Any, Any).Return(fd2, nil)
	m.EXPECT().Ioctl(Any, uintptr(syscall.SIOCGIFHWADDR), Any).Return(uintptr(RetValOnSuccess), R2Zero, ErrnoSuccess)
	m.EXPECT().EpollCreate1(Any).Return(RetValOnFail, ErrorWithNoMessage)
	m.EXPECT().Close(Any).Return(nil)

	psSyscall.Syscall = m

	tapDev := TapDevice{}

	got := tapDev.Open()
	if got != psErr.CantCreateEpollInstance {
		t.Errorf("TapDevice.Open() = %v; want %v", got, psErr.CantCreateEpollInstance)
	}
}

// Fail when EpollCtl() returns error.
func TestTapDevice_Open_7(t *testing.T) {
	psLog.DisableOutput()
	defer psLog.EnableOutput()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := psSyscall.NewMockISyscall(ctrl)
	m.EXPECT().Open(Any, Any, Any).Return(10, nil)
	m.EXPECT().Ioctl(Any, uintptr(syscall.TUNSETIFF), Any).Return(uintptr(RetValOnSuccess), R2Zero, ErrnoSuccess)
	m.EXPECT().Socket(Any, Any, Any).Return(11, nil)
	m.EXPECT().Ioctl(Any, uintptr(syscall.SIOCGIFHWADDR), Any).Return(uintptr(RetValOnSuccess), R2Zero, ErrnoSuccess)
	m.EXPECT().EpollCreate1(Any).Return(12, nil)
	m.EXPECT().EpollCtl(Any, Any, Any, Any).Return(ErrorWithNoMessage)
	m.EXPECT().Close(Any).Return(nil).AnyTimes()

	psSyscall.Syscall = m

	tapDev := TapDevice{}

	got := tapDev.Open()
	if got != psErr.CantModifyIOResourceParameter {
		t.Errorf("TapDevice.Open() = %v; want %v", got, psErr.CantModifyIOResourceParameter)
	}
}

func TestTapDevice_Close(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := psSyscall.NewMockISyscall(ctrl)
	m.EXPECT().Close(Any).Return(nil)

	psSyscall.Syscall = m

	tapDev := TapDevice{}

	got := tapDev.Close()
	if got != psErr.OK {
		t.Errorf("TapDevice.Close() = %v; want %v", got, psErr.OK)
	}
}

func TestTapDevice_Transmit(t *testing.T) {
	psLog.DisableOutput()
	defer psLog.EnableOutput()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := psSyscall.NewMockISyscall(ctrl)
	m.EXPECT().Write(Any, Any).Return(0, nil)

	psSyscall.Syscall = m

	tapDev := TapDevice{}

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
	m.EXPECT().EpollWait(Any, Any, Any).Return(0, nil)

	psSyscall.Syscall = m

	tapDev := TapDevice{}

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
	m.EXPECT().EpollWait(Any, Any, Any).Return(1, nil)
	m.EXPECT().Read(Any, Any).Return(150, nil)

	psSyscall.Syscall = m

	tapDev := TapDevice{}

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
	m.EXPECT().Close(Any).Return(nil)

	psSyscall.Syscall = m

	tapDev := TapDevice{}

	got := tapDev.Poll(true)
	if got != psErr.Terminated {
		t.Errorf("TapDevice.Poll() = %v; want %v", got, psErr.Terminated)
	}
}

// Fail when EpollWait() is interrupted.
func TestTapDevice_Poll_4(t *testing.T) {
	psLog.DisableOutput()
	defer psLog.EnableOutput()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := psSyscall.NewMockISyscall(ctrl)
	m.EXPECT().EpollWait(Any, Any, Any).Return(RetValOnFail, syscall.EINTR)

	psSyscall.Syscall = m

	tapDev := TapDevice{}

	got := tapDev.Poll(false)
	if got != psErr.Interrupted {
		t.Errorf("TapDevice.Poll() = %v; want %v", got, psErr.Interrupted)
	}
}

func init() {
	ErrorWithNoMessage = errors.New("")
}
