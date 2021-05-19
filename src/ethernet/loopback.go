package ethernet

import (
	psErr "github.com/42milez/ProtocolStack/src/error"
	psSyscall "github.com/42milez/ProtocolStack/src/syscall"
	"math"
)

const LoopbackMTU = math.MaxUint16
const LoopbackIpAddr = "127.0.0.1"
const LoopbackNetmask = "255.0.0.0"

type LoopbackOperation struct{}

func (v LoopbackOperation) Open(dev *Device, sc psSyscall.ISyscall) psErr.Error {
	return psErr.Error{Code: psErr.OK}
}

func (v LoopbackOperation) Close(dev *Device, sc psSyscall.ISyscall) psErr.Error {
	return psErr.Error{Code: psErr.OK}
}

func (v LoopbackOperation) Transmit(dev *Device, sc psSyscall.ISyscall) psErr.Error {
	return psErr.Error{Code: psErr.OK}
}

func (v LoopbackOperation) Poll(dev *Device, sc psSyscall.ISyscall, terminate bool) psErr.Error {
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
		Op:        LoopbackOperation{},
	}
	return dev
}
