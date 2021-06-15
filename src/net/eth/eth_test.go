package eth

import (
	psLog "github.com/42milez/ProtocolStack/src/log"
	"github.com/42milez/ProtocolStack/src/worker"
	"github.com/golang/mock/gomock"
	"sync"
	"testing"
)

func TestStart(t *testing.T) {
	_, teardown := setupEthTest(t)
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
	_, teardown := setupEthTest(t)
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

func setupEthTest(t *testing.T) (ctrl *gomock.Controller, teardown func()) {
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
