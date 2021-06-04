//go:generate mockgen -source=device.go -destination=device_mock.go -package=$GOPACKAGE -self_package=github.com/42milez/ProtocolStack/src/$GOPACKAGE

package eth

import (
	psErr "github.com/42milez/ProtocolStack/src/error"
)

const DevFlagUp DevFlag = 0x0001
const DevFlagLoopback DevFlag = 0x0010
const DevFlagBroadcast DevFlag = 0x0020
const DevFlagNeedArp DevFlag = 0x0100

const (
	DevTypeEthernet DevType = iota
	DevTypeLoopback
	DevTypeNull
)

type DevFlag uint16
type DevType int

func (v DevType) String() string {
	return devTypes[v]
}

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
	Flag() DevFlag
	MTU() uint16
	Priv() Privilege
}

type Device struct {
	Type_ DevType
	Name_ string
	Addr_ EthAddr
	Flag_ DevFlag
	MTU_  uint16
	Priv_ Privilege
}

func (p *Device) Up() {
	p.Flag_ |= DevFlagUp
}

func (p *Device) Down() {
	p.Flag_ &= ^DevFlagUp
}

func (p *Device) Equal(pp IDevice) bool {
	return p.Name_ == pp.Name()
}

func (p *Device) IsUp() bool {
	return p.Flag_&DevFlagUp == 1
}

func (p *Device) Type() DevType {
	return p.Type_
}

func (p *Device) Name() string {
	return p.Name_
}

func (p *Device) Addr() EthAddr {
	return p.Addr_
}

func (p *Device) Flag() DevFlag {
	return p.Flag_
}

func (p *Device) MTU() uint16 {
	return p.MTU_
}

func (p *Device) Priv() Privilege {
	return p.Priv_
}

type Privilege struct {
	Name string
	FD   int
}

var devTypes = [...]string{
	0: "Ethernet",
	1: "Loopback",
	2: "Null",
}
