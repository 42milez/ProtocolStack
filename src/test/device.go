package test

import (
	"github.com/42milez/ProtocolStack/src/mw"
	"github.com/42milez/ProtocolStack/src/net/eth"
)

var DeviceBuilder deviceBuilder

type deviceBuilder struct{}

func (deviceBuilder) Default() *eth.TapDevice {
	return &eth.TapDevice{
		Device: mw.Device{
			Type_: mw.EthernetDevice,
			Name_: "net0",
			MTU_:  mw.EthPayloadLenMax,
			Flag_: mw.BroadcastFlag | mw.NeedArpFlag,
			Addr_: mw.EthAddr{11, 12, 13, 14, 15, 16},
			Priv_: mw.Privilege{
				FD:   3,
				Name: "tap0",
			},
		},
	}
}

func init() {
	DeviceBuilder = deviceBuilder{}
}
