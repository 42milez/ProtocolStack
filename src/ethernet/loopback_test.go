package ethernet

import (
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
