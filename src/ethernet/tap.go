package ethernet

import (
	"github.com/42milez/ProtocolStack/src/device"
	"log"
)

func tapOpen(dev *device.Device) int {
	return 0
}

func tapClose(dev *device.Device) int {
	return 0
}

func tapTransmit(dev *device.Device) int {
	return 0
}

func tapPoll(dev *device.Device) int {
	return 0
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

	log.Println("TAP device generated.")

	return dev, nil
}
