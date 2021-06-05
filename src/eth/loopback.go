package eth

import (
	psErr "github.com/42milez/ProtocolStack/src/error"
	"github.com/42milez/ProtocolStack/src/mw"
	"math"
)

const LoopbackMTU = math.MaxUint16

type LoopbackDevice struct {
	mw.Device
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

func (p *LoopbackDevice) Transmit(dst mw.Addr, payload []byte, typ mw.EthType) psErr.E {
	return psErr.OK
}
