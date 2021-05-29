//go:generate mockgen -source=device.go -destination=device_mock.go -package=$GOPACKAGE -self_package=github.com/42milez/ProtocolStack/src/$GOPACKAGE

package ethernet

import (
	psErr "github.com/42milez/ProtocolStack/src/error"
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
		return "ETHERNET"
	case DevTypeLoopback:
		return "LOOPBACK"
	case DevTypeNull:
		return "NULL"
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
	DevType() DevType
	DevName() string
	PrivDevName() string
	EthAddr() EthAddr
	BroadcastEthAddr() EthAddr
	PeerEthAddr() EthAddr
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
	return p.Name == pp.DevName()
}

func (p *Device) IsUp() bool {
	return p.FLAG&DevFlagUp == 1
}

func (p *Device) DevType() DevType {
	return p.Type
}

func (p *Device) DevName() string {
	return p.Name
}

func (p *Device) PrivDevName() string {
	return p.Priv.Name
}

func (p *Device) EthAddr() EthAddr {
	return p.Addr
}

func (p *Device) BroadcastEthAddr() EthAddr {
	return p.Addr
}

func (p *Device) PeerEthAddr() EthAddr {
	return p.Addr
}
