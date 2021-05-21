package ethernet

import (
	psErr "github.com/42milez/ProtocolStack/src/error"
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
	_, name, _ := v.Info()
	return dev.Name == name
}

func (dev *Device) Info() (string, string, string) {
	return dev.Type.String(), dev.Name, dev.Priv.Name
}

func (dev *Device) IsUp() bool {
	return dev.FLAG&DevFlagUp == 1
}
