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
	Transmit(dst EthAddr, payload []byte, typ EthType) psErr.E
	Up()
	Down()
	Equal(dev IDevice) bool
	IsUp() bool
	Type() DevType
	Name() string
	Addr() EthAddr
	Broadcast() EthAddr
	Peer() EthAddr
	Flag() DevFlag
	HeaderLen() uint16
	MTU() uint16
	Priv() Privilege
}

type Device struct {
	_Type      DevType
	_Name      string
	_Addr      EthAddr
	_Broadcast EthAddr
	_Peer      EthAddr
	_Flag      DevFlag
	_HeaderLen uint16
	_MTU       uint16
	_Priv      Privilege
}

type Privilege struct {
	Name string
	FD   int
}

func (p *Device) Up() {
	p._Flag |= DevFlagUp
}

func (p *Device) Down() {
	p._Flag &= ^DevFlagUp
}

func (p *Device) Equal(pp IDevice) bool {
	return p._Name == pp.Name()
}

func (p *Device) IsUp() bool {
	return p._Flag&DevFlagUp == 1
}

func (p *Device) Type() DevType {
	return p._Type
}

func (p *Device) Name() string {
	return p._Name
}

func (p *Device) Addr() EthAddr {
	return p._Addr
}

func (p *Device) Broadcast() EthAddr {
	return p._Broadcast
}

func (p *Device) Peer() EthAddr {
	return p._Peer
}

func (p *Device) Flag() DevFlag {
	return p._Flag
}

func (p *Device) HeaderLen() uint16 {
	return p._HeaderLen
}

func (p *Device) MTU() uint16 {
	return p._MTU
}

func (p *Device) Priv() Privilege {
	return p._Priv
}
