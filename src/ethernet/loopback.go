package ethernet

import (
	psErr "github.com/42milez/ProtocolStack/src/error"
	s "github.com/42milez/ProtocolStack/src/syscall"
	"math"
)

const LoopbackMTU = math.MaxUint16
const LoopbackIpAddr = "127.0.0.1"
const LoopbackNetmask = "255.0.0.0"

func loopbackTransmit(dev *Device, sc s.ISyscall) psErr.Error {
	return psErr.Error{Code: psErr.OK}
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
