package eth

import (
	psErr "github.com/42milez/ProtocolStack/src/error"
	"math"
)

const LoopbackMTU = math.MaxUint16

type LoopbackDevice struct {
	Device
}

func (p *LoopbackDevice) Open() psErr.E {
	return psErr.OK
}

func (p *LoopbackDevice) Close() psErr.E {
	return psErr.OK
}

func (p *LoopbackDevice) Poll(terminate bool) psErr.E {
	return psErr.OK
}

func (p *LoopbackDevice) Transmit(dst Addr, payload []byte, typ EthType) psErr.E {
	return psErr.OK
}
