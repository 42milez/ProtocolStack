package ethernet

import (
	"github.com/42milez/ProtocolStack/src/device"
	"math"
)

const LoopbackMTU = math.MaxUint16

const LoopbackIpAddr = "127.0.0.1"
const LoopbackNetmask = "255.0.0.0"

func transmit(dev *device.Device) int {
	return 0
}

func GenLoopbackDevice() *device.Device {
	return &device.Device{
		Type: device.DevTypeLoopback,
		MTU: LoopbackMTU,
		HeaderLen: 0,
		AddrLen: 0,
		FLAG: device.DevFlagLoopback,
		Op: device.Operation{
			Open: nil,
			Close: nil,
			Transmit: transmit,
			Poll: nil,
		},
	}
}
