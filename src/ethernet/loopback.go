package ethernet

import (
	psErr "github.com/42milez/ProtocolStack/src/error"
	"github.com/42milez/ProtocolStack/src/middleware"
	psSyscall "github.com/42milez/ProtocolStack/src/syscall"
	"math"
	"strconv"
)

const LoopbackMTU = math.MaxUint16
const LoopbackIpAddr = "127.0.0.1"
const LoopbackNetmask = "255.0.0.0"

type LoopbackDevice struct{
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

// GenLoopbackDevice generates loopback device object.
func GenLoopbackDevice() (*LoopbackDevice, psErr.Error) {
	dev := &LoopbackDevice{
		Device{
			Name:      "net" + strconv.Itoa(middleware.NextDeviceIndex()),
			Type:      DevTypeLoopback,
			MTU:       LoopbackMTU,
			HeaderLen: 0,
			FLAG:      DevFlagLoopback,
			Syscall:   &psSyscall.Syscall{},
		},
	}
	return dev, psErr.Error{Code: psErr.OK}
}
