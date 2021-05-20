package middleware

import (
	psErr "github.com/42milez/ProtocolStack/src/error"
	psLog "github.com/42milez/ProtocolStack/src/log"
	"testing"
)

func TestSetup(t *testing.T) {
	psLog.DisableOutput()
	defer psLog.EnableOutput()

	got := Setup()
	if got.Code != psErr.OK {
		t.Errorf("Setup() = %v; want %v", got, psErr.OK)
	}
}
