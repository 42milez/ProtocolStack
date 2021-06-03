package ethernet

import (
	psErr "github.com/42milez/ProtocolStack/src/error"
	psSyscall "github.com/42milez/ProtocolStack/src/syscall"
	"testing"
)

func TestLoopbackDevice_Open(t *testing.T) {
	ctrl, teardown := setup(t)
	defer teardown()

	psSyscall.Syscall = psSyscall.NewMockISyscall(ctrl)

	loopbackDev := LoopbackDevice{}

	got := loopbackDev.Open()
	if got != psErr.OK {
		t.Errorf("LoopbackDevice.Open() = %v; want %v", got, psErr.OK)
	}
}

func TestLoopbackDevice_Close(t *testing.T) {
	ctrl, teardown := setup(t)
	defer teardown()

	psSyscall.Syscall = psSyscall.NewMockISyscall(ctrl)

	loopbackDev := LoopbackDevice{}

	got := loopbackDev.Close()
	if got != psErr.OK {
		t.Errorf("LoopbackDevice.Close() = %v; want %v", got, psErr.OK)
	}
}

func TestLoopbackDevice_Poll(t *testing.T) {
	ctrl, teardown := setup(t)
	defer teardown()

	psSyscall.Syscall = psSyscall.NewMockISyscall(ctrl)

	loopbackDev := LoopbackDevice{}

	got := loopbackDev.Poll(false)
	if got != psErr.OK {
		t.Errorf("LoopbackDevice.Poll() = %v; want %v", got, psErr.OK)
	}
}

func TestLoopbackDevice_Transmit(t *testing.T) {
	ctrl, teardown := setup(t)
	defer teardown()

	psSyscall.Syscall = psSyscall.NewMockISyscall(ctrl)

	loopbackDev := LoopbackDevice{}

	got := loopbackDev.Transmit(EthAddr{}, make([]byte, 0), EthTypeArp)
	if got != psErr.OK {
		t.Errorf("LoopbackDevice.Transmit() = %v; want %v", got, psErr.OK)
	}
}
