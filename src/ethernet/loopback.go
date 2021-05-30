package ethernet

import (
	psErr "github.com/42milez/ProtocolStack/src/error"
	"math"
)

const LoopbackMTU = math.MaxUint16

type LoopbackDevice struct {
	Device
}

func (dev *LoopbackDevice) Open() psErr.E {
	return psErr.OK
}

func (dev *LoopbackDevice) Close() psErr.E {
	return psErr.OK
}

func (dev *LoopbackDevice) Poll(terminate bool) psErr.E {
	return psErr.OK
}

func (dev *LoopbackDevice) Transmit(dst EthAddr, payload []byte, typ EthType) psErr.E {
	return psErr.OK
}
