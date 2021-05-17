package ethernet

import (
	e "github.com/42milez/ProtocolStack/src/error"
	s "github.com/42milez/ProtocolStack/src/sys"
	"math"
)

const LoopbackMTU = math.MaxUint16
const LoopbackIpAddr = "127.0.0.1"
const LoopbackNetmask = "255.0.0.0"

func loopbackTransmit(dev *Device, sc *s.Syscall) e.Error {
	return e.Error{Code: e.OK}
}

// GenLoopbackDevice generates loopback device object.
func GenLoopbackDevice() *Device {
	dev := &Device{
		Type:      DevTypeLoopback,
		MTU:       LoopbackMTU,
		HeaderLen: 0,
		AddrLen:   0,
		FLAG:      DevFlagLoopback,
		Op: Operation{
			Open:     nil,
			Close:    nil,
			Transmit: loopbackTransmit,
			Poll:     nil,
		},
	}
	return dev
}
