package ethernet

import (
	psErr "github.com/42milez/ProtocolStack/src/error"
	psLog "github.com/42milez/ProtocolStack/src/log"
	mockSyscall "github.com/42milez/ProtocolStack/src/syscall/mock"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestGenLoopbackDevice(t *testing.T) {
	want := &Device{
		Type:      DevTypeLoopback,
		MTU:       LoopbackMTU,
		HeaderLen: 0,
		FLAG:      DevFlagLoopback,
	}
	got, _ := GenLoopbackDevice()
	if d := cmp.Diff(got, want); d != "" {
		t.Errorf("GenLoopbackDevice() differs: (-got +want)\n%s", d)
	}
}

func TestLoopbackOperation_Open(t *testing.T) {
	psLog.DisableOutput()
	defer psLog.EnableOutput()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	loopbackDev := LoopbackDevice{Device{Syscall: mockSyscall.NewMockISyscall(ctrl)}}

	got := loopbackDev.Open()
	if got.Code != psErr.OK {
		t.Errorf("LoopbackOperation.Open() = %v; want %v", got.Code, psErr.OK)
	}
}

func TestLoopbackOperation_Close(t *testing.T) {
	psLog.DisableOutput()
	defer psLog.EnableOutput()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	loopbackDev := LoopbackDevice{Device{Syscall: mockSyscall.NewMockISyscall(ctrl)}}

	got := loopbackDev.Close()
	if got.Code != psErr.OK {
		t.Errorf("LoopbackOperation.Close() = %v; want %v", got.Code, psErr.OK)
	}
}

func TestLoopbackOperation_Transmit(t *testing.T) {
	psLog.DisableOutput()
	defer psLog.EnableOutput()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	loopbackDev := LoopbackDevice{Device{Syscall: mockSyscall.NewMockISyscall(ctrl)}}

	got := loopbackDev.Transmit()
	if got.Code != psErr.OK {
		t.Errorf("LoopbackOperation.Transmit() = %v; want %v", got.Code, psErr.OK)
	}
}

func TestLoopbackOperation_Poll(t *testing.T) {
	psLog.DisableOutput()
	defer psLog.EnableOutput()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	loopbackDev := LoopbackDevice{Device{Syscall: mockSyscall.NewMockISyscall(ctrl)}}

	got := loopbackDev.Poll(false)
	if got.Code != psErr.OK {
		t.Errorf("LoopbackOperation.Poll() = %v; want %v", got.Code, psErr.OK)
	}
}
