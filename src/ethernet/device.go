//go:generate mockgen -source=device.go -destination=device_mock.go -package=$GOPACKAGE -self_package=github.com/42milez/ProtocolStack/src/$GOPACKAGE

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
	Open() psErr.E
	Close() psErr.E
	Poll(terminate bool) psErr.E
	Transmit(dest EthAddr, payload []byte, typ EthType) psErr.E
	Up()
	Down()
	Equal(dev IDevice) bool
	IsUp() bool
	EthAddrs() (addr EthAddr, broadcast EthAddr, peer EthAddr)
	Names() (name string, privName string)
	Typ() (typ DevType)
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

func (p *Device) Up() {
	p.FLAG |= DevFlagUp
}

func (p *Device) Down() {
	p.FLAG &= ^DevFlagUp
}

func (p *Device) Equal(pp IDevice) bool {
	name, _ := pp.Names()
	return p.Name == name
}

func (p *Device) IsUp() bool {
	return p.FLAG&DevFlagUp == 1
}

func (p *Device) EthAddrs() (addr EthAddr, broadcast EthAddr, peer EthAddr) {
	addr = p.Addr
	broadcast = p.Broadcast
	peer = p.Peer
	return
}

func (p *Device) Names() (name string, privName string) {
	name = p.Name
	privName = p.Priv.Name
	return
}

func (p *Device) Typ() (typ DevType) {
	typ = p.Type
	return
}
