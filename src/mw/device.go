//go:generate mockgen -source=device.go -destination=device_mock.go -package=$GOPACKAGE -self_package=github.com/42milez/ProtocolStack/src/$GOPACKAGE

package mw

const UpFlag DevFlag = 0x0001
const LoopbackFlag DevFlag = 0x0010
const BroadcastFlag DevFlag = 0x0020
const NeedArpFlag DevFlag = 0x0100
const (
	EthernetDevice DevType = iota
	LoopbackDevice
	NullDevice
)

var devTypes = [...]string{
	0: "Ethernet",
	1: "Loopback",
	2: "Null",
}

type DevFlag uint16
type DevType int

func (v DevType) String() string {
	return devTypes[v]
}

type IDevice interface {
	Open() error
	Close() error
	Poll() error
	Transmit(dst EthAddr, payload []byte, typ EthType) error
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
	p.Flag_ |= UpFlag
}

func (p *Device) Down() {
	p.Flag_ &= ^UpFlag
}

func (p *Device) Equal(pp IDevice) bool {
	return p.Name_ == pp.Name()
}

func (p *Device) IsUp() bool {
	return p.Flag_&UpFlag == 1
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
