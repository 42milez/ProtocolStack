package ethernet

import (
	"github.com/42milez/ProtocolStack/src/device"
)

func tapOpen(dev *device.Device) error {
	return nil
}

func tapClose(dev *device.Device) error {
	return nil
}

func tapTransmit(dev *device.Device) error {
	return nil
}

func tapPoll(dev *device.Device) error {
	return nil
}

// GenTapDevice generates TAP device object.
func GenTapDevice(name string, mac MAC) (*device.Device, error) {
	dev := &device.Device{
		Type:      device.DevTypeEthernet,
		MTU:       EthPayloadSizeMax,
		FLAG:      device.DevFlagBroadcast | device.DevFlagNeedArp,
		HeaderLen: EthHeaderSize,
		AddrLen:   EthAddrLen,
		Broadcast: EthAddrBroadcast,
		Op: device.Operation{
			Open:     tapOpen,
			Close:    tapClose,
			Transmit: tapTransmit,
			Poll:     tapPoll,
		},
		Priv: device.Privilege{Name: name, FD: -1},
	}

	if addr, err := mac.Byte(); err != nil {
		return nil, err
	} else {
		dev.Addr = addr
	}

	return dev, nil
}
