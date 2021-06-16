package ip

import (
	psLog "github.com/42milez/ProtocolStack/src/log"
	"github.com/42milez/ProtocolStack/src/worker"
	"github.com/golang/mock/gomock"
	"sync"
	"testing"
)

func TestPacketID_Next(t *testing.T) {
	want := uint16(0)
	got := id.Next()
	if got != want {
		t.Errorf("PacketID.Next() = %d; want %d", got, want)
	}
}

func TestStart(t *testing.T) {
	_, teardown := setupIpTest(t)
	defer teardown()

	var wg sync.WaitGroup
	_ = Start(&wg)
	rcvMonMsg := <-rcvMonCh
	sndMonMsg := <-sndMonCh

	if rcvMonMsg.Current != worker.Running || sndMonMsg.Current != worker.Running {
		t.Errorf("Start() failed")
	}
}

func TestStop(t *testing.T) {
	_, teardown := setupIpTest(t)
	defer teardown()

	var wg sync.WaitGroup
	_ = Start(&wg)
	<-rcvMonCh
	<-sndMonCh
	Stop()
	rcvMonMsg := <-rcvMonCh
	sndMonMsg := <-sndMonCh

	if rcvMonMsg.Current != worker.Stopped || sndMonMsg.Current != worker.Stopped {
		t.Errorf("Stop() failed")
	}
}

func setupIpTest(t *testing.T) (ctrl *gomock.Controller, teardown func()) {
	ctrl = gomock.NewController(t)
	psLog.DisableOutput()
	reset := func() {
		psLog.EnableOutput()
	}
	teardown = func() {
		ctrl.Finish()
		reset()
	}
	return
}
