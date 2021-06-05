package eth

import "github.com/42milez/ProtocolStack/src/mw"

// GenLoopbackDevice generates loopback device object.
func GenLoopbackDevice(name string) *LoopbackDevice {
	return &LoopbackDevice{
		Device: mw.Device{
			Type_: mw.DevTypeLoopback,
			Name_: name,
			Addr_: mw.EthAny,
			Flag_: mw.DevFlagLoopback,
			MTU_:  LoopbackMTU,
		},
	}
}

// GenTapDevice generates TAP device object.
func GenTapDevice(devName string, privName string, addr mw.EthAddr) *TapDevice {
	return &TapDevice{
		Device: mw.Device{
			Type_: mw.DevTypeEthernet,
			Name_: devName,
			Addr_: addr,
			Flag_: mw.DevFlagBroadcast | mw.DevFlagNeedArp,
			MTU_:  mw.EthPayloadLenMax,
			Priv_: mw.Privilege{
				FD:   -1,
				Name: privName,
			},
		},
	}
}
