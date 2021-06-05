package eth

import (
	psErr "github.com/42milez/ProtocolStack/src/error"
	psLog "github.com/42milez/ProtocolStack/src/log"
	"github.com/42milez/ProtocolStack/src/mw"
	psSyscall "github.com/42milez/ProtocolStack/src/syscall"
	"github.com/golang/mock/gomock"
	"testing"
)

func setupLoopbackTest(t *testing.T) (ctrl *gomock.Controller, teardown func()) {
	psLog.DisableOutput()
	ctrl = gomock.NewController(t)
	teardown = func() {
		ctrl.Finish()
		psLog.EnableOutput()
	}
	return
}

func TestLoopbackDevice_Open(t *testing.T) {
	ctrl, teardown := setupLoopbackTest(t)
	defer teardown()

	psSyscall.Syscall = psSyscall.NewMockISyscall(ctrl)

	loopbackDev := LoopbackDevice{}

	got := loopbackDev.Open()
	if got != psErr.OK {
		t.Errorf("LoopbackDevice.Open() = %v; want %v", got, psErr.OK)
	}
}

func TestLoopbackDevice_Close(t *testing.T) {
	ctrl, teardown := setupLoopbackTest(t)
	defer teardown()

	psSyscall.Syscall = psSyscall.NewMockISyscall(ctrl)

	loopbackDev := LoopbackDevice{}

	got := loopbackDev.Close()
	if got != psErr.OK {
		t.Errorf("LoopbackDevice.Close() = %v; want %v", got, psErr.OK)
	}
}

func TestLoopbackDevice_Poll(t *testing.T) {
	ctrl, teardown := setupLoopbackTest(t)
	defer teardown()

	psSyscall.Syscall = psSyscall.NewMockISyscall(ctrl)

	loopbackDev := LoopbackDevice{}

	got := loopbackDev.Poll(false)
	if got != psErr.OK {
		t.Errorf("LoopbackDevice.Poll() = %v; want %v", got, psErr.OK)
	}
}

func TestLoopbackDevice_Transmit(t *testing.T) {
	ctrl, teardown := setupLoopbackTest(t)
	defer teardown()

	psSyscall.Syscall = psSyscall.NewMockISyscall(ctrl)

	loopbackDev := LoopbackDevice{}

	got := loopbackDev.Transmit(mw.EthAddr{}, make([]byte, 0), mw.ARP)
	if got != psErr.OK {
		t.Errorf("LoopbackDevice.Transmit() = %v; want %v", got, psErr.OK)
	}
}
