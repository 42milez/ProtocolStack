package ethernet

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestGenTapDevice_1(t *testing.T) {
	devName := "tap0"
	devEthAddr := EthAddr{11, 12, 13, 14, 15, 16}
	want := &Device{
		Type:      DevTypeEthernet,
		MTU:       EthPayloadSizeMax,
		FLAG:      DevFlagBroadcast | DevFlagNeedArp,
		HeaderLen: EthHeaderSize,
		Addr:      devEthAddr,
		AddrLen:   EthAddrLen,
		Broadcast: EthAddrBroadcast,
		Op:        TapOperation{},
		Priv:      Privilege{FD: -1, Name: devName},
	}
	got, _ := GenTapDevice(devName, devEthAddr)
	if d := cmp.Diff(got, want); d != "" {
		t.Errorf("GenTapDevice() differs: (-got +want)\n%s", d)
	}
}
