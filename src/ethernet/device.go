package ethernet

import (
	psErr "github.com/42milez/ProtocolStack/src/error"
	psLog "github.com/42milez/ProtocolStack/src/log"
	psSyscall "github.com/42milez/ProtocolStack/src/syscall"
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

type Operation interface {
	Open(dev *Device, sc psSyscall.ISyscall) psErr.Error
	Close(dev *Device, sc psSyscall.ISyscall) psErr.Error
	Transmit(dev *Device, sc psSyscall.ISyscall) psErr.Error
	Poll(dev *Device, sc psSyscall.ISyscall, terminate bool) psErr.Error
}

type Privilege struct {
	Name string
	FD   int
}

type Device struct {
	Type      DevType
	Name      string
	Addr      EthAddr
	AddrLen   uint16
	Broadcast EthAddr
	Peer      EthAddr
	FLAG      DevFlag
	HeaderLen uint16
	MTU       uint16
	Op        Operation
	Priv      Privilege
}

func (dev *Device) Open() psErr.Error {
	if (dev.FLAG & DevFlagUp) != 0 {
		psLog.W("device is already opened")
		psLog.W("\tname: %v (%v) ", dev.Name, dev.Priv.Name)
		return psErr.Error{Code: psErr.AlreadyOpened}
	}

	if err := dev.Op.Open(dev, &psSyscall.Syscall{}); err.Code != psErr.OK {
		psLog.E("can't open a device")
		psLog.E("\tname: %v (%v) ", dev.Name, dev.Priv.Name)
		psLog.E("\ttype: %v ", dev.Type)
		return psErr.Error{Code: psErr.CantOpen}
	}

	dev.FLAG |= DevFlagUp

	psLog.I("device opened")
	psLog.I("\tname: %v (%v) ", dev.Name, dev.Priv.Name)

	return psErr.Error{Code: psErr.OK}
}
