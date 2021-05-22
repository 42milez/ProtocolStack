package ethernet

import (
	psErr "github.com/42milez/ProtocolStack/src/error"
	psLog "github.com/42milez/ProtocolStack/src/log"
	psSyscall "github.com/42milez/ProtocolStack/src/syscall"
	"github.com/golang/mock/gomock"
	"testing"
)

func TestLoopbackDevice_Open_OK(t *testing.T) {
	psLog.DisableOutput()
	defer psLog.EnableOutput()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	loopbackDev := LoopbackDevice{Device{Syscall: psSyscall.NewMockISyscall(ctrl)}}

	got := loopbackDev.Open()
	if got.Code != psErr.OK {
		t.Errorf("LoopbackDevice.Open() = %v; want %v", got.Code, psErr.OK)
	}
}

func TestLoopbackDevice_Close_OK(t *testing.T) {
	psLog.DisableOutput()
	defer psLog.EnableOutput()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	loopbackDev := LoopbackDevice{Device{Syscall: psSyscall.NewMockISyscall(ctrl)}}

	got := loopbackDev.Close()
	if got.Code != psErr.OK {
		t.Errorf("LoopbackDevice.Close() = %v; want %v", got.Code, psErr.OK)
	}
}

func TestLoopbackDevice_Poll_OK(t *testing.T) {
	psLog.DisableOutput()
	defer psLog.EnableOutput()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	loopbackDev := LoopbackDevice{Device{Syscall: psSyscall.NewMockISyscall(ctrl)}}

	got := loopbackDev.Poll(false)
	if got.Code != psErr.OK {
		t.Errorf("LoopbackDevice.Poll() = %v; want %v", got.Code, psErr.OK)
	}
}

func TestLoopbackDevice_Transmit_OK(t *testing.T) {
	psLog.DisableOutput()
	defer psLog.EnableOutput()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	loopbackDev := LoopbackDevice{Device{Syscall: psSyscall.NewMockISyscall(ctrl)}}

	got := loopbackDev.Transmit()
	if got.Code != psErr.OK {
		t.Errorf("LoopbackDevice.Transmit() = %v; want %v", got.Code, psErr.OK)
	}
}
