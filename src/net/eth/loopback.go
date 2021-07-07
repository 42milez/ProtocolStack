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

func (p *LoopbackDevice) Open() error {
	return psErr.OK
}

func (p *LoopbackDevice) Close() error {
	return psErr.OK
}

func (p *LoopbackDevice) Poll() error {
	return psErr.OK
}

func (p *LoopbackDevice) Transmit(dst mw.EthAddr, payload []byte, typ mw.EthType) error {
	return psErr.OK
}
