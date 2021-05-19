package ethernet

import (
	psErr "github.com/42milez/ProtocolStack/src/error"
	psLog "github.com/42milez/ProtocolStack/src/log"
	mockSyscall "github.com/42milez/ProtocolStack/src/mock/syscall"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestGenLoopbackDevice(t *testing.T) {
	want := &Device{
		Type:      DevTypeLoopback,
		MTU:       LoopbackMTU,
		HeaderLen: 0,
		AddrLen:   0,
		FLAG:      DevFlagLoopback,
		Op:        LoopbackOperation{},
	}
	got := GenLoopbackDevice()
	if d := cmp.Diff(got, want); d != "" {
		t.Errorf("GenLoopbackDevice() differs: (-got +want)\n%s", d)
	}
}

func TestLoopbackOperation_Open(t *testing.T) {
	psLog.DisableOutput()
	defer psLog.EnableOutput()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mockSyscall.NewMockISyscall(ctrl)

	loopbackOp := LoopbackOperation{}
	dev := &Device{}

	got := loopbackOp.Open(dev, m)
	if got.Code != psErr.OK {
		t.Errorf("LoopbackOperation.Open() = %v; want %v", got.Code, psErr.OK)
	}
}

func TestLoopbackOperation_Close(t *testing.T) {
	psLog.DisableOutput()
	defer psLog.EnableOutput()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mockSyscall.NewMockISyscall(ctrl)

	loopbackOp := LoopbackOperation{}
	dev := &Device{}

	got := loopbackOp.Close(dev, m)
	if got.Code != psErr.OK {
		t.Errorf("LoopbackOperation.Close() = %v; want %v", got.Code, psErr.OK)
	}
}

func TestLoopbackOperation_Transmit(t *testing.T) {
	psLog.DisableOutput()
	defer psLog.EnableOutput()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mockSyscall.NewMockISyscall(ctrl)

	loopbackOp := LoopbackOperation{}
	dev := &Device{}

	got := loopbackOp.Transmit(dev, m)
	if got.Code != psErr.OK {
		t.Errorf("LoopbackOperation.Transmit() = %v; want %v", got.Code, psErr.OK)
	}
}

func TestLoopbackOperation_Poll(t *testing.T) {
	psLog.DisableOutput()
	defer psLog.EnableOutput()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mockSyscall.NewMockISyscall(ctrl)

	loopbackOp := LoopbackOperation{}
	dev := &Device{}

	got := loopbackOp.Poll(dev, m, false)
	if got.Code != psErr.OK {
		t.Errorf("LoopbackOperation.Poll() = %v; want %v", got.Code, psErr.OK)
	}
}
