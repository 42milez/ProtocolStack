package network

import (
	psErr "github.com/42milez/ProtocolStack/src/error"
	"github.com/42milez/ProtocolStack/src/ethernet"
	psLog "github.com/42milez/ProtocolStack/src/log"
	"syscall"
)

var devices []ethernet.IDevice
var interfaces []*Iface
var protocols []Protocol

type Timer struct {
	Name     string
	Interval syscall.Timeval
	Last     syscall.Timeval
	Handler  Handler
}

type Handler func(data []byte, dev ethernet.IDevice) psErr.Error

func GetIface(dev ethernet.IDevice, family AddrFamily) {

}

func RegisterDevice(dev ethernet.IDevice) psErr.Error {
	for _, d := range devices {
		if d.Equal(dev) {
			typ, n1, n2 := d.Info()
			psLog.W("device is already registered")
			psLog.W("\ttype:        %v", typ)
			psLog.W("\tname:        %v", n1)
			psLog.W("\tname (priv): %v", n2)
			return psErr.Error{Code: psErr.CantRegister}
		}
	}
	devices = append(devices, dev)
	typ, n1, n2 := dev.Info()
	psLog.I("device registered")
	psLog.I("\ttype:        %v", typ)
	psLog.I("\tname:        %v", n1)
	psLog.I("\tname (priv): %v", n2)
	return psErr.Error{Code: psErr.OK}
}

func RegisterInterface(iface *Iface, dev ethernet.IDevice) psErr.Error {
	for _, i := range interfaces {
		if i.Dev.Equal(dev) && i.Family == iface.Family {
			psLog.W("interface is already registered: %v ", i.Family.String())
			return psErr.Error{Code: psErr.CantRegister}
		}
	}

	interfaces = append(interfaces, iface)
	iface.Dev = dev

	_, name1, name2 := dev.Info()
	psLog.I("interface attached")
	psLog.I("\tip:     %v ", iface.Unicast.String())
	psLog.I("\tdevice: %v (%v) ", name1, name2)

	return psErr.Error{Code: psErr.OK}
}

func Up() psErr.Error {
	for _, dev := range devices {
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
		dev.Up()
		psLog.I("device opened")
		psLog.I("\tname: %v (%v) ", name1, name2)
	}
	return psErr.Error{Code: psErr.OK}
}

func Poll(terminate bool) psErr.Error {
	for _, dev := range devices {
		if !dev.IsUp() {
			continue
		}
		if err := dev.Poll(terminate); err.Code != psErr.OK {
			return err
		}
	}
	return psErr.Error{Code: psErr.OK}
}

func init() {
	devices = make([]ethernet.IDevice, 0)
	interfaces = make([]*Iface, 0)
	protocols = make([]Protocol, 0)

	// ARP
	register(ProtocolTypeArp, ArpInputHandler)

	// ICMP
	// ...

	// IP
	register(ProtocolTypeIp, IpInputHandler)

	// TCP
	// ...

	// UDP
	// ...
}

func register(protocolType ProtocolType, handler Handler) psErr.Error {
	p := Protocol{
		Type:    protocolType,
		Handler: handler,
	}
	protocols = append(protocols, p)
	return psErr.Error{Code: psErr.OK}
}
