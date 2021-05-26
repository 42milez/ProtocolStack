package network

import (
	psErr "github.com/42milez/ProtocolStack/src/error"
	"github.com/42milez/ProtocolStack/src/ethernet"
	psLog "github.com/42milez/ProtocolStack/src/log"
	"syscall"
)

var DeviceRepo *deviceRepo
var IfaceRepo *ifaceRepo

type deviceRepo struct {
	devices []ethernet.IDevice
}

func (p *deviceRepo) NextNumber() int {
	return len(p.devices)
}

func (p *deviceRepo) Poll(terminate bool) psErr.Error {
	for _, dev := range p.devices {
		if !dev.IsUp() {
			continue
		}
		if err := dev.Poll(terminate); err.Code != psErr.OK {
			return err
		}
	}
	return psErr.Error{Code: psErr.OK}
}

func (p *deviceRepo) Register(dev ethernet.IDevice) psErr.Error {
	for _, d := range p.devices {
		if d.Equal(dev) {
			typ := d.Typ()
			name, privName := d.Names()
			psLog.W("device is already registered")
			psLog.W("\ttype:        %v", typ)
			psLog.W("\tname:        %v", name)
			psLog.W("\tname (priv): %v", privName)
			return psErr.Error{Code: psErr.CantRegister}
		}
	}
	p.devices = append(p.devices, dev)
	typ := dev.Typ()
	name, privName := dev.Names()
	psLog.I("device registered")
	psLog.I("\ttype:        %v", typ)
	psLog.I("\tname:        %v", name)
	psLog.I("\tname (priv): %v", privName)
	return psErr.Error{Code: psErr.OK}
}

func (p *deviceRepo) Up() psErr.Error {
	for _, dev := range p.devices {
		typ := dev.Typ()
		name, privName := dev.Names()
		if dev.IsUp() {
			psLog.W("device is already opened")
			psLog.W("\tname: %v (%v) ", name, privName)
			return psErr.Error{Code: psErr.AlreadyOpened}
		}
		if err := dev.Open(); err.Code != psErr.OK {
			psLog.E("can't open a device")
			psLog.E("\tname: %v (%v) ", name, privName)
			psLog.E("\ttype: %v ", typ)
			return psErr.Error{Code: psErr.CantOpen}
		}
		dev.Up()
		psLog.I("device opened")
		psLog.I("\tname: %v (%v) ", name, privName)
	}
	return psErr.Error{Code: psErr.OK}
}

type ifaceRepo struct {
	ifaces []*Iface
}

func (p *ifaceRepo) Register(iface *Iface, dev ethernet.IDevice) psErr.Error {
	for _, i := range p.ifaces {
		if i.Dev.Equal(dev) && i.Family == iface.Family {
			psLog.W("interface is already registered: %v ", i.Family.String())
			return psErr.Error{Code: psErr.CantRegister}
		}
	}

	p.ifaces = append(p.ifaces, iface)
	iface.Dev = dev

	name, privName := dev.Names()
	psLog.I("interface attached")
	psLog.I("\tip:     %v ", iface.Unicast.String())
	psLog.I("\tdevice: %v (%v) ", name, privName)

	return psErr.Error{Code: psErr.OK}
}

type Handler func(data []byte, dev ethernet.IDevice) psErr.Error

type Timer struct {
	Name     string
	Interval syscall.Timeval
	Last     syscall.Timeval
	Handler  Handler
}

func init() {
	DeviceRepo = &deviceRepo{}
	IfaceRepo = &ifaceRepo{}
}
