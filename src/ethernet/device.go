package ethernet

import (
	psErr "github.com/42milez/ProtocolStack/src/error"
	psLog "github.com/42milez/ProtocolStack/src/log"
	psSyscall "github.com/42milez/ProtocolStack/src/syscall"
	"github.com/google/go-cmp/cmp"
)

type DevType int

const (
	DevTypeEthernet DevType = iota
	DevTypeLoopback
	DevTypeNull
)

func (t DevType) String() string {
	switch t {
	case DevTypeEthernet:
		return "DEVICE_TYPE_ETHERNET"
	case DevTypeLoopback:
		return "DEVICE_TYPE_LOOPBACK"
	case DevTypeNull:
		return "DEVICE_TYPE_NULL"
	default:
		return "UNKNOWN"
	}
}

type DevFlag uint16

const DevFlagUp DevFlag = 0x0001
const DevFlagLoopback DevFlag = 0x0010
const DevFlagBroadcast DevFlag = 0x0020
const DevFlagP2P DevFlag = 0x0040
const DevFlagNeedArp DevFlag = 0x0100

type IDevice interface {
	Open() psErr.Error
	Close() psErr.Error
	Poll(terminate bool) psErr.Error
	Transmit() psErr.Error
	Enable()
	Disable()
	Equal(dev IDevice) bool
	Info() (string, string, string)
	IsUp() bool
}

type Device struct {
	Type      DevType
	Name      string
	Addr      EthAddr
	Broadcast EthAddr
	Peer      EthAddr
	FLAG      DevFlag
	HeaderLen uint16
	MTU       uint16
	Priv      Privilege
	Syscall   psSyscall.ISyscall
}

type Privilege struct {
	Name string
	FD   int
}

func (dev *Device) Enable() {
	dev.FLAG |= DevFlagUp
}

func (dev *Device) Disable() {
	dev.FLAG &= ^DevFlagUp
}

func (dev *Device) Equal(v IDevice) bool {
	return cmp.Equal(dev, v)
}

func (dev *Device) Info() (string, string, string) {
	return dev.Type.String(), dev.Name, dev.Priv.Name
}

func (dev *Device) IsUp() bool {
	if dev.FLAG&DevFlagUp == 0 {
		return false
	}
	return true
}

func Up(dev IDevice) psErr.Error {
	typ, name1, name2 := dev.Info()

	if dev.IsUp() {
		psLog.W("device is already opened")
		psLog.W("\tname: %v (%v) ", name1, name2)
		return psErr.Error{Code: psErr.AlreadyOpened}
	}

	if err := dev.Open(); err.Code != psErr.OK {
		psLog.E("can't open a device")
		psLog.E("\tname: %v (%v) ", name1, name2)
		psLog.E("\ttype: %v ", typ)
		return psErr.Error{Code: psErr.CantOpen}
	}

	dev.Enable()

	psLog.I("device opened")
	psLog.I("\tname: %v (%v) ", name1, name2)

	return psErr.Error{Code: psErr.OK}
}
