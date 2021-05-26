package network

import (
	psErr "github.com/42milez/ProtocolStack/src/error"
	"github.com/42milez/ProtocolStack/src/ethernet"
	psLog "github.com/42milez/ProtocolStack/src/log"
	"syscall"
)

var DeviceRepo *deviceRepo
var IfaceRepo *ifaceRepo
var RouteRepo *routeRepo

type Handler func(data []byte, dev ethernet.IDevice) psErr.Error

type Timer struct {
	Name     string
	Interval syscall.Timeval
	Last     syscall.Timeval
	Handler  Handler
}

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
			psLog.W("device already registered")
			psLog.W("\ttype: %s", typ)
			psLog.W("\tname: %s (%s)", name, privName)
			return psErr.Error{Code: psErr.CantRegister}
		}
	}
	p.devices = append(p.devices, dev)
	typ := dev.Typ()
	name, privName := dev.Names()
	addr, broadcast, peer := dev.EthAddrs()
	psLog.I("device registered")
	psLog.I("\ttype:      %s", typ)
	psLog.I("\tname:      %s (%s)", name, privName)
	psLog.I("\taddr:      %s", addr)
	psLog.I("\tbroadcast: %s", broadcast)
	psLog.I("\tpeer:      %s", peer)
	return psErr.Error{Code: psErr.OK}
}

func (p *deviceRepo) Up() psErr.Error {
	for _, dev := range p.devices {
		typ := dev.Typ()
		name, privName := dev.Names()
		if dev.IsUp() {
			psLog.W("device already opened")
			psLog.W("\ttype: %s ", typ)
			psLog.W("\tname: %s (%s) ", name, privName)
			return psErr.Error{Code: psErr.AlreadyOpened}
		}
		if err := dev.Open(); err.Code != psErr.OK {
			psLog.E("can't open a device")
			psLog.E("\ttype: %s ", typ)
			psLog.E("\tname: %s (%s) ", name, privName)
			return psErr.Error{Code: psErr.CantOpen}
		}
		dev.Up()
		psLog.I("device opened")
		psLog.I("\ttype: %s ", typ)
		psLog.I("\tname: %s (%s) ", name, privName)
	}
	return psErr.Error{Code: psErr.OK}
}

type ifaceRepo struct {
	ifaces []*Iface
}

func (p *ifaceRepo) Get(dev ethernet.IDevice, family AddrFamily) *Iface {
	for _, v := range p.ifaces {
		if v.Dev.Equal(dev) && v.Family == family {
			return v
		}
	}
	return nil
}

func (p *ifaceRepo) Register(iface *Iface, dev ethernet.IDevice) psErr.Error {
	for _, i := range p.ifaces {
		if i.Dev.Equal(dev) && i.Family == iface.Family {
			psLog.W("interface already registered: %v ", i.Family.String())
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

type route struct {
	Network IP
	Netmask IP
	NextHop IP
	iface   *Iface
}

type routeRepo struct {
	routes []*route
}

func (p *routeRepo) Register(network IP, nextHop IP, iface *Iface) {
	route := &route{
		Network: network,
		Netmask: iface.Netmask,
		NextHop: nextHop,
		iface:   iface,
	}
	p.routes = append(p.routes, route)
	name, privName := iface.Dev.Names()
	psLog.I("route registered")
	psLog.I("\tnetwork:  %s", route.Network)
	psLog.I("\tnetmask:  %s", route.Netmask)
	psLog.I("\tunicast:  %s", iface.Unicast)
	psLog.I("\tnext hop: %s", nextHop)
	psLog.I("\tdevice:   %s (%s) ", name, privName)
}

func (p *routeRepo) RegisterDefaultGateway(iface *Iface, nextHop IP) {
	route := &route{
		Network: V4Zero,
		Netmask: V4Zero,
		NextHop: nextHop,
		iface:   iface,
	}
	p.routes = append(p.routes, route)
	name, privName := iface.Dev.Names()
	psLog.I("default gateway registered")
	psLog.I("\tnetwork:  %s", route.Network)
	psLog.I("\tnetmask:  %s", route.Netmask)
	psLog.I("\tunicast:  %s", iface.Unicast)
	psLog.I("\tnext hop: %s", nextHop)
	psLog.I("\tdevice:   %s (%s) ", name, privName)
}

func init() {
	DeviceRepo = &deviceRepo{}
	IfaceRepo = &ifaceRepo{}
	RouteRepo = &routeRepo{}
}
