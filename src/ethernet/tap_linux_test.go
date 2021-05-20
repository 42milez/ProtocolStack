package ethernet

import (
	"errors"
	psErr "github.com/42milez/ProtocolStack/src/error"
	psLog "github.com/42milez/ProtocolStack/src/log"
	mockSyscall "github.com/42milez/ProtocolStack/src/syscall/mock"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"syscall"
	"testing"
)

func TestGenTapDevice_A(t *testing.T) {
	devName := "tap0"
	devEthAddr := EthAddr{11, 12, 13, 14, 15, 16}
	want := &Device{
		Type:      DevTypeEthernet,
		MTU:       EthPayloadSizeMax,
		FLAG:      DevFlagBroadcast | DevFlagNeedArp,
		HeaderLen: EthHeaderSize,
		Addr:      devEthAddr,
		Broadcast: EthAddrBroadcast,
		Priv:      Privilege{FD: -1, Name: devName},
	}
	got, _ := GenTapDevice(devName, devEthAddr)
	if d := cmp.Diff(got, want); d != "" {
		t.Errorf("GenTapDevice() differs: (-got +want)\n%s", d)
	}
}

func TestGenTapDevice_B(t *testing.T) {
	devName := "tap000000000000000"
	devEthAddr := EthAddr{11, 12, 13, 14, 15, 16}
	want1 := (*Device)(nil)
	want2 := psErr.Error{
		Code: psErr.CantCreate,
		Msg:  "device name must be less than or equal to 16 characters",
	}
	got1, got2 := GenTapDevice(devName, devEthAddr)
	if !cmp.Equal(got1, want1) {
		t.Errorf("GenTapDevice() = %v; want %v", got1, want1)
	}
	if d := cmp.Diff(got2, want2); d != "" {
		t.Errorf("GenTapDevice() differs: (-got +want)\n%s", d)
	}
}

func TestTapOperation_Open_A(t *testing.T) {
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

	tapDev := TapDevice{Device{Syscall:m}}

	got := tapDev.Open()
	if got.Code != psErr.OK {
		t.Errorf("TapOperation.Open() = %v; want %v", got.Code, psErr.OK)
	}
}

func TestTapOperation_Open_B(t *testing.T) {
	psLog.DisableOutput()
	defer psLog.EnableOutput()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mockSyscall.NewMockISyscall(ctrl)
	m.EXPECT().Open(gomock.Any(), gomock.Any(), gomock.Any()).Return(-1, errors.New(""))

	tapDev := TapDevice{Device{Syscall:m}}

	got := tapDev.Open()
	if got.Code != psErr.CantOpen {
		t.Errorf("TapOperation.Open() = %v; want %v", got.Code, psErr.CantOpen)
	}
}

func TestTapOperation_Open_C(t *testing.T) {
	psLog.DisableOutput()
	defer psLog.EnableOutput()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ioctlRetVal := -1
	m := mockSyscall.NewMockISyscall(ctrl)
	m.EXPECT().Open(gomock.Any(), gomock.Any(), gomock.Any()).Return(10, nil)
	m.EXPECT().Ioctl(gomock.Any(), gomock.Any(), gomock.Any()).Return(uintptr(ioctlRetVal), uintptr(0), syscall.EBADF)
	m.EXPECT().Close(gomock.Any()).Return(nil)

	tapDev := TapDevice{Device{Syscall:m}}

	got := tapDev.Open()
	if got.Code != psErr.CantOpen {
		t.Errorf("TapOperation.Open() = %v; want %v", got.Code, psErr.CantOpen)
	}
}

func TestTapOperation_Open_D(t *testing.T) {
	psLog.DisableOutput()
	defer psLog.EnableOutput()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mockSyscall.NewMockISyscall(ctrl)
	m.EXPECT().Open(gomock.Any(), gomock.Any(), gomock.Any()).Return(10, nil)
	m.EXPECT().Ioctl(gomock.Any(), uintptr(syscall.TUNSETIFF), gomock.Any()).Return(uintptr(0), uintptr(0), syscall.Errno(0))
	m.EXPECT().Socket(gomock.Any(), gomock.Any(), gomock.Any()).Return(-1, errors.New(""))

	tapDev := TapDevice{Device{Syscall:m}}

	got := tapDev.Open()
	if got.Code != psErr.CantOpen {
		t.Errorf("TapOperation.Open() = %v; want %v", got.Code, psErr.CantOpen)
	}
}

func TestTapOperation_Open_E(t *testing.T) {
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

	tapDev := TapDevice{Device{Syscall:m}}

	got := tapDev.Open()
	if got.Code != psErr.CantOpen {
		t.Errorf("TapOperation.Open() = %v; want %v", got.Code, psErr.CantOpen)
	}
}

func TestTapOperation_Open_F(t *testing.T) {
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

	tapDev := TapDevice{Device{Syscall:m}}

	got := tapDev.Open()
	if got.Code != psErr.CantOpen {
		t.Errorf("TapOperation.Open() = %v; want %v", got.Code, psErr.CantOpen)
	}
}

func TestTapOperation_Open_G(t *testing.T) {
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

	tapDev := TapDevice{Device{Syscall:m}}

	got := tapDev.Open()
	if got.Code != psErr.CantOpen {
		t.Errorf("TapOperation.Open() = %v; want %v", got.Code, psErr.CantOpen)
	}
}

func TestTapOperation_Close(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mockSyscall.NewMockISyscall(ctrl)
	m.EXPECT().Close(gomock.Any()).Return(nil)

	tapDev := TapDevice{Device{Syscall:m}}

	got := tapDev.Close()
	if got.Code != psErr.OK {
		t.Errorf("TapOperation.Close() = %v; want %v", got.Code, psErr.OK)
	}
}

func TestTapOperation_Transmit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tapDev := TapDevice{Device{Syscall:mockSyscall.NewMockISyscall(ctrl)}}

	got := tapDev.Transmit()
	if got.Code != psErr.OK {
		t.Errorf("TapOperation.Transmit() = %v; want %v", got.Code, psErr.OK)
	}
}

func TestTapOperation_Poll_NoEventOccurred(t *testing.T) {
	psLog.DisableOutput()
	defer psLog.EnableOutput()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mockSyscall.NewMockISyscall(ctrl)
	m.EXPECT().EpollWait(gomock.Any(), gomock.Any(), gomock.Any()).Return(0, nil)

	tapDev := TapDevice{Device{Syscall:mockSyscall.NewMockISyscall(ctrl)}}

	got := tapDev.Poll(false)
	if got.Code != psErr.OK {
		t.Errorf("TapOperation.Poll() = %v; want %v", got.Code, psErr.OK)
	}
}

func TestTapOperation_Poll_EventOccurred(t *testing.T) {
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

	tapDev := TapDevice{Device{Syscall:m}}

	got := tapDev.Poll(false)
	if got.Code != psErr.OK {
		t.Errorf("TapOperation.Poll() = %v; want %v", got.Code, psErr.OK)
	}
}

func TestTapOperation_Poll_Terminated(t *testing.T) {
	psLog.DisableOutput()
	defer psLog.EnableOutput()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mockSyscall.NewMockISyscall(ctrl)
	m.EXPECT().Close(gomock.Any()).Return(nil)

	tapDev := TapDevice{Device{Syscall:m}}

	got := tapDev.Poll(true)
	if got.Code != psErr.OK {
		t.Errorf("TapOperation.Poll() = %v; want %v", got.Code, psErr.OK)
	}
}

func TestTapOperation_Poll_FailOnEpollWait(t *testing.T) {
	psLog.DisableOutput()
	defer psLog.EnableOutput()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mockSyscall.NewMockISyscall(ctrl)
	m.EXPECT().EpollWait(gomock.Any(), gomock.Any(), gomock.Any()).Return(-1, errors.New(""))
	m.EXPECT().Close(gomock.Any()).Return(nil)

	tapDev := TapDevice{Device{Syscall:m}}

	got := tapDev.Poll(false)
	if got.Code != psErr.Interrupted {
		t.Errorf("TapOperation.Poll() = %v; want %v", got.Code, psErr.Interrupted)
	}
}
