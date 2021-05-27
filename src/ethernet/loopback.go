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

func (dev *LoopbackDevice) Open() psErr.E {
	return psErr.OK
}

func (dev *LoopbackDevice) Close() psErr.E {
	return psErr.OK
}

func (dev *LoopbackDevice) Poll(terminate bool) psErr.E {
	return psErr.OK
}

func (dev *LoopbackDevice) Transmit(dest EthAddr, payload []byte, typ EthType) psErr.E {
	return psErr.OK
}
