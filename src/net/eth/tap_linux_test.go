package eth

import (
	"errors"
	psErr "github.com/42milez/ProtocolStack/src/error"
	psLog "github.com/42milez/ProtocolStack/src/log"
	"github.com/42milez/ProtocolStack/src/mw"
	psSyscall "github.com/42milez/ProtocolStack/src/syscall"
	"github.com/golang/mock/gomock"
	"syscall"
	"testing"
)

const ErrnoSuccess = syscall.Errno(0)

var ErrorWithNoMessage error
var RetValOnFail = -1
var any = gomock.Any()

func setupTapLinuxTest(t *testing.T) (ctrl *gomock.Controller, teardown func()) {
	psLog.DisableOutput()
	ctrl = gomock.NewController(t)
	teardown = func() {
		ctrl.Finish()
		psLog.EnableOutput()
	}
	return
}

func TestTapDevice_Open_1(t *testing.T) {
	ctrl, teardown := setupTapLinuxTest(t)
	defer teardown()

	m := psSyscall.NewMockISyscall(ctrl)
	m.EXPECT().Open(any, any, any).Return(3, nil)
	m.EXPECT().Ioctl(any, syscall.TUNSETIFF, any).Return(ErrnoSuccess)
	m.EXPECT().Socket(any, any, any).Return(4, nil)
	m.EXPECT().Ioctl(any, syscall.SIOCGIFHWADDR, any).Return(ErrnoSuccess)
	m.EXPECT().Close(any).Return(nil)
	m.EXPECT().EpollCreate1(any).Return(5, nil)
	m.EXPECT().EpollCtl(any, any, any, any).Return(nil)
	psSyscall.Syscall = m

	tapDev := TapDevice{}

	got := tapDev.Open()
	if got != psErr.OK {
		t.Errorf("TapDevice.Open() = %v; want %v", got, psErr.OK)
	}
}

// Fail when Open() returns error.
func TestTapDevice_Open_2(t *testing.T) {
	ctrl, teardown := setupTapLinuxTest(t)
	defer teardown()

	m := psSyscall.NewMockISyscall(ctrl)
	m.EXPECT().Open(any, any, any).Return(RetValOnFail, ErrorWithNoMessage)

	psSyscall.Syscall = m

	tapDev := TapDevice{}

	got := tapDev.Open()
	if got != psErr.CantOpenIOResource {
		t.Errorf("TapDevice.Open() = %v; want %v", got, psErr.CantOpenIOResource)
	}
}

// Fail when Ioctl() returns error.
func TestTapDevice_Open_3(t *testing.T) {
	ctrl, teardown := setupTapLinuxTest(t)
	defer teardown()

	m := psSyscall.NewMockISyscall(ctrl)
	m.EXPECT().Open(any, any, any).Return(10, nil)
	m.EXPECT().Ioctl(any, any, any).Return(syscall.EBADF)
	m.EXPECT().Close(any).Return(nil)

	psSyscall.Syscall = m

	tapDev := TapDevice{}

	got := tapDev.Open()
	if got != psErr.CantModifyIOResourceParameter {
		t.Errorf("TapDevice.Open() = %v; want %v", got, psErr.CantModifyIOResourceParameter)
	}
}

// Fail when Socket() returns error.
func TestTapDevice_Open_4(t *testing.T) {
	ctrl, teardown := setupTapLinuxTest(t)
	defer teardown()

	m := psSyscall.NewMockISyscall(ctrl)

	fd := 3

	m.EXPECT().Open(any, any, any).Return(fd, nil)
	m.EXPECT().Ioctl(any, syscall.TUNSETIFF, any).Return(ErrnoSuccess)
	m.EXPECT().Socket(any, any, any).Return(RetValOnFail, ErrorWithNoMessage)
	m.EXPECT().Close(any).Return(nil)

	psSyscall.Syscall = m

	tapDev := TapDevice{}

	got := tapDev.Open()
	if got != psErr.CantCreateEndpoint {
		t.Errorf("TapDevice.Open() = %v; want %v", got, psErr.CantCreateEndpoint)
	}
}

// Fail when Ioctl() returns error.
func TestTapDevice_Open_5(t *testing.T) {
	ctrl, teardown := setupTapLinuxTest(t)
	defer teardown()

	m := psSyscall.NewMockISyscall(ctrl)

	fd1 := 3
	fd2 := 4

	m.EXPECT().Open(any, any, any).Return(fd1, nil)
	m.EXPECT().Ioctl(any, syscall.TUNSETIFF, any).Return(ErrnoSuccess)
	m.EXPECT().Socket(any, any, any).Return(fd2, nil)
	m.EXPECT().Ioctl(any, syscall.SIOCGIFHWADDR, any).Return(syscall.EBADF)
	m.EXPECT().Close(any).Return(nil)

	psSyscall.Syscall = m

	tapDev := TapDevice{}

	got := tapDev.Open()
	if got != psErr.CantModifyIOResourceParameter {
		t.Errorf("TapDevice.Open() = %v; want %v", got, psErr.CantModifyIOResourceParameter)
	}
}

// Fail when EpollCreate1() returns error.
func TestTapDevice_Open_6(t *testing.T) {
	ctrl, teardown := setupTapLinuxTest(t)
	defer teardown()

	m := psSyscall.NewMockISyscall(ctrl)

	fd1 := 3
	fd2 := 4

	m.EXPECT().Open(any, any, any).Return(fd1, nil)
	m.EXPECT().Ioctl(any, syscall.TUNSETIFF, any).Return(ErrnoSuccess)
	m.EXPECT().Socket(any, any, any).Return(fd2, nil)
	m.EXPECT().Ioctl(any, syscall.SIOCGIFHWADDR, any).Return(ErrnoSuccess)
	m.EXPECT().EpollCreate1(any).Return(RetValOnFail, ErrorWithNoMessage)
	m.EXPECT().Close(any).Return(nil)

	psSyscall.Syscall = m

	tapDev := TapDevice{}

	got := tapDev.Open()
	if got != psErr.CantCreateEpollInstance {
		t.Errorf("TapDevice.Open() = %v; want %v", got, psErr.CantCreateEpollInstance)
	}
}

// Fail when EpollCtl() returns error.
func TestTapDevice_Open_7(t *testing.T) {
	ctrl, teardown := setupTapLinuxTest(t)
	defer teardown()

	m := psSyscall.NewMockISyscall(ctrl)
	m.EXPECT().Open(any, any, any).Return(10, nil)
	m.EXPECT().Ioctl(any, syscall.TUNSETIFF, any).Return(ErrnoSuccess)
	m.EXPECT().Socket(any, any, any).Return(11, nil)
	m.EXPECT().Ioctl(any, syscall.SIOCGIFHWADDR, any).Return(ErrnoSuccess)
	m.EXPECT().EpollCreate1(any).Return(12, nil)
	m.EXPECT().EpollCtl(any, any, any, any).Return(ErrorWithNoMessage)
	m.EXPECT().Close(any).Return(nil).AnyTimes()

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
	m.EXPECT().Close(any).Return(nil)

	psSyscall.Syscall = m

	tapDev := TapDevice{}

	got := tapDev.Close()
	if got != psErr.OK {
		t.Errorf("TapDevice.Close() = %v; want %v", got, psErr.OK)
	}
}

func TestTapDevice_Transmit(t *testing.T) {
	ctrl, teardown := setupTapLinuxTest(t)
	defer teardown()

	m := psSyscall.NewMockISyscall(ctrl)
	m.EXPECT().Write(any, any).Return(0, nil)

	psSyscall.Syscall = m

	tapDev := TapDevice{}

	got := tapDev.Transmit(mw.EthAddr{}, make([]byte, 0), mw.ARP)
	if got != psErr.OK {
		t.Errorf("TapDevice.Transmit() = %v; want %v", got, psErr.OK)
	}
}

// Success when no event occurs.
func TestTapDevice_Poll_1(t *testing.T) {
	ctrl, teardown := setupTapLinuxTest(t)
	defer teardown()

	m := psSyscall.NewMockISyscall(ctrl)
	m.EXPECT().EpollWait(any, any, any).Return(0, nil)

	psSyscall.Syscall = m

	tapDev := TapDevice{}

	got := tapDev.Poll()
	if got != psErr.OK {
		t.Errorf("TapDevice.Poll() = %v; want %v", got, psErr.OK)
	}
}

// Success when an event occurs.
func TestTapDevice_Poll_2(t *testing.T) {
	ctrl, teardown := setupTapLinuxTest(t)
	defer teardown()

	m := psSyscall.NewMockISyscall(ctrl)
	m.EXPECT().EpollWait(any, any, any).Return(1, nil)
	m.EXPECT().Read(any, any).Return(150, nil)

	psSyscall.Syscall = m

	tapDev := TapDevice{}

	got := tapDev.Poll()
	if got != psErr.OK {
		t.Errorf("TapDevice.Poll() = %v; want %v", got, psErr.OK)
	}
}

// Fail when EpollWait() is interrupted.
func TestTapDevice_Poll_4(t *testing.T) {
	ctrl, teardown := setupTapLinuxTest(t)
	defer teardown()

	m := psSyscall.NewMockISyscall(ctrl)
	m.EXPECT().EpollWait(any, any, any).Return(RetValOnFail, syscall.EINTR)

	psSyscall.Syscall = m

	tapDev := TapDevice{}

	got := tapDev.Poll()
	if got != psErr.Interrupted {
		t.Errorf("TapDevice.Poll() = %v; want %v", got, psErr.Interrupted)
	}
}

// Fail when EpollWait() returns EFAULT.
func TestTapDevice_Poll_5(t *testing.T) {
	ctrl, teardown := setupTapLinuxTest(t)
	defer teardown()

	m := psSyscall.NewMockISyscall(ctrl)
	m.EXPECT().EpollWait(any, any, any).Return(RetValOnFail, syscall.EFAULT)

	psSyscall.Syscall = m

	tapDev := TapDevice{}

	got := tapDev.Poll()
	if got != psErr.SyscallError {
		t.Errorf("TapDevice.Poll() = %v; want %v", got, psErr.SyscallError)
	}
}

// Fail when ReadFrame() failed.
func TestTapDevice_Poll_6(t *testing.T) {
	ctrl, teardown := setupTapLinuxTest(t)
	defer teardown()

	m := psSyscall.NewMockISyscall(ctrl)
	m.EXPECT().EpollWait(any, any, any).Return(1, nil)
	m.EXPECT().Read(any, any).Return(RetValOnFail, syscall.EIO)

	psSyscall.Syscall = m

	tapDev := TapDevice{}

	got := tapDev.Poll()
	if got != psErr.Error {
		t.Errorf("TapDevice.Poll() = %v; want %v", got, psErr.Error)
	}
}

func init() {
	ErrorWithNoMessage = errors.New("")
}
