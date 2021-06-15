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

	rcvStatus := <-rcvMonCh
	sndStatus := <-sndMonCh

	if rcvStatus.Current != worker.Running || sndStatus.Current != worker.Running {
		t.Errorf("Start() failed")
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
