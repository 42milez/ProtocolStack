package eth

import (
	psErr "github.com/42milez/ProtocolStack/src/error"
	psLog "github.com/42milez/ProtocolStack/src/log"
	"github.com/golang/mock/gomock"
	"sync"
	"testing"
)

func TestStart(t *testing.T) {
	_, teardown := setup(t)
	defer teardown()

	var wg sync.WaitGroup

	want := psErr.OK
	got := Start(&wg)
	if got != want {
		t.Errorf("Start() = %s; want %s", got, want)
	}
}

func setup(t *testing.T) (ctrl *gomock.Controller, teardown func()) {
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
