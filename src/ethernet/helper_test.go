package ethernet

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestGenLoopbackDevice(t *testing.T) {
	want := &LoopbackDevice{
		Device: Device{
			Type_:   DevTypeLoopback,
			Name_:   "net0",
			MTU_:    LoopbackMTU,
			HdrLen_: 0,
			Flag_:   DevFlagLoopback,
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
	devEthAddr := EthAddr{11, 12, 13, 14, 15, 16}
	want := &TapDevice{
		Device: Device{
			Type_:   DevTypeEthernet,
			Name_:   devName,
			MTU_:    EthPayloadSizeMax,
			Flag_:   DevFlagBroadcast | DevFlagNeedArp,
			HdrLen_: EthHdrSize,
			Addr_:   devEthAddr,
			Priv_: Privilege{
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
