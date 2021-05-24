package ethernet

import (
	psErr "github.com/42milez/ProtocolStack/src/error"
	"math"
)

const LoopbackMTU = math.MaxUint16
const LoopbackIpAddr = "127.0.0.1"
const LoopbackBroadcast = "127.255.255.255"
const LoopbackNetmask = "255.0.0.0"
const LoopbackNetwork = "127.0.0.0"

type LoopbackDevice struct {
	Device
}

func (dev *LoopbackDevice) Open() psErr.Error {
	return psErr.Error{Code: psErr.OK}
}

func (dev *LoopbackDevice) Close() psErr.Error {
	return psErr.Error{Code: psErr.OK}
}

func (dev *LoopbackDevice) Poll(terminate bool) psErr.Error {
	return psErr.Error{Code: psErr.OK}
}

func (dev *LoopbackDevice) Transmit() psErr.Error {
	return psErr.Error{Code: psErr.OK}
}
