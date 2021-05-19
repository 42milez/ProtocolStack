package ethernet

import (
	psErr "github.com/42milez/ProtocolStack/src/error"
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestGenTapDevice_A(t *testing.T) {
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

func TestGenTapDevice_B(t *testing.T) {
	devName := "tap000000000000000"
	devEthAddr := EthAddr{11, 12, 13, 14, 15, 16}
	want1 := (*Device)(nil)
	want2 := psErr.Error{
		Code: psErr.CantCreate,
		Msg: "device name must be less than or equal to 16 characters",
	}
	got1, got2 := GenTapDevice(devName, devEthAddr)
	if ! cmp.Equal(got1, want1) {
		t.Errorf("GenTapDevice() = %v; want %v", got1, want1)
	}
	if d := cmp.Diff(got2, want2); d != "" {
		t.Errorf("GenTapDevice() differs: (-got +want)\n%s", d)
	}
}
