package ethernet

// GenLoopbackDevice generates loopback device object.
func GenLoopbackDevice(name string) *LoopbackDevice {
	return &LoopbackDevice{
		Device: Device{
			Type_: DevTypeLoopback,
			Name_: name,
			Addr_: EthAddrAny,
			Flag_: DevFlagLoopback,
			MTU_:  LoopbackMTU,
		},
	}
}

// GenTapDevice generates TAP device object.
func GenTapDevice(devName string, privName string, addr EthAddr) *TapDevice {
	return &TapDevice{
		Device: Device{
			Type_: DevTypeEthernet,
			Name_: devName,
			Addr_: addr,
			Flag_: DevFlagBroadcast | DevFlagNeedArp,
			MTU_:  EthPayloadLenMax,
			Priv_: Privilege{
				FD:   -1,
				Name: privName,
			},
		},
	}
}
