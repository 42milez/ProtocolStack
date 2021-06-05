package eth

import (
	"github.com/42milez/ProtocolStack/src/mw"
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestGenLoopbackDevice(t *testing.T) {
	want := &LoopbackDevice{
		Device: mw.Device{
			Type_: mw.DevTypeLoopback,
			Name_: "net0",
			MTU_:  LoopbackMTU,
			Flag_: mw.DevFlagLoopback,
		},
	}
	got := GenLoopbackDevice("net0")
	if d := cmp.Diff(*got, *want); d != "" {
		t.Errorf("GenLoopbackDevice() differs: (-got +want)\n%s", d)
	}
}

func TestGenTapDevice(t *testing.T) {
	devName := "net0"
	privName := "tap0"
	devEthAddr := mw.EthAddr{11, 12, 13, 14, 15, 16}
	want := &TapDevice{
		Device: mw.Device{
			Type_: mw.DevTypeEthernet,
			Name_: devName,
			MTU_:  mw.EthPayloadLenMax,
			Flag_: mw.DevFlagBroadcast | mw.DevFlagNeedArp,
			Addr_: devEthAddr,
			Priv_: mw.Privilege{
				FD:   -1,
				Name: privName,
			},
		},
	}
	got := GenTapDevice(devName, privName, devEthAddr)
	if d := cmp.Diff(*got, *want); d != "" {
		t.Errorf("GenTapDevice() differs: (-got +want)\n%s", d)
	}
}
