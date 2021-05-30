package ethernet

// GenLoopbackDevice generates loopback device object.
func GenLoopbackDevice(name string) *LoopbackDevice {
	return &LoopbackDevice{
		Device: Device{
			Type_:      DevTypeLoopback,
			Name_:      name,
			Flag_:      DevFlagLoopback,
			HeaderLen_: 0,
			MTU_:       LoopbackMTU,
		},
	}
}

// GenTapDevice generates TAP device object.
func GenTapDevice(devName string, privName string, addr EthAddr) *TapDevice {
	return &TapDevice{
		Device: Device{
			Type_:      DevTypeEthernet,
			Name_:      devName,
			Addr_:      addr,
			Broadcast_: EthAddrBroadcast,
			Flag_:      DevFlagBroadcast | DevFlagNeedArp,
			HeaderLen_: EthHeaderSize,
			MTU_:       EthPayloadSizeMax,
			Priv_: Privilege{
				FD:   -1,
				Name: privName,
			},
		},
	}
}
