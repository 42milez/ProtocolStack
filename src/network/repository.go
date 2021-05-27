package network

import (
	"fmt"
	psErr "github.com/42milez/ProtocolStack/src/error"
	"github.com/42milez/ProtocolStack/src/ethernet"
	psLog "github.com/42milez/ProtocolStack/src/log"
	"syscall"
)

var DeviceRepo *deviceRepo
var IfaceRepo *ifaceRepo
var RouteRepo *routeRepo

type Handler func(data []byte, dev ethernet.IDevice) psErr.E

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

func (p *deviceRepo) Poll(terminate bool) psErr.E {
	for _, dev := range p.devices {
		if !dev.IsUp() {
			continue
		}
		if err := dev.Poll(terminate); err != psErr.OK {
			if err == psErr.Interrupted {
				return psErr.OK
			}
			psLog.E(fmt.Sprintf("IDevice.Poll() failed: %s", err))
			return psErr.Error
		}
	}
	return psErr.OK
}

func (p *deviceRepo) Register(dev ethernet.IDevice) psErr.E {
	for _, d := range p.devices {
		if d.Equal(dev) {
			typ := d.Typ()
			name, privName := d.Names()
			psLog.W("Device is already registered")
			psLog.W(fmt.Sprintf("\ttype: %s", typ))
			psLog.W(fmt.Sprintf("\tname: %s (%s)", name, privName))
			return psErr.Error
		}
	}
	p.devices = append(p.devices, dev)
	typ := dev.Typ()
	name, privName := dev.Names()
	addr, broadcast, peer := dev.EthAddrs()
	psLog.I("Device was registered")
	psLog.I(fmt.Sprintf("\ttype:      %s", typ))
	psLog.I(fmt.Sprintf("\tname:      %s (%s)", name, privName))
	psLog.I(fmt.Sprintf("\taddr:      %s", addr))
	psLog.I(fmt.Sprintf("\tbroadcast: %s", broadcast))
	psLog.I(fmt.Sprintf("\tpeer:      %s", peer))
	return psErr.OK
}

func (p *deviceRepo) Up() psErr.E {
	for _, dev := range p.devices {
		typ := dev.Typ()
		name, privName := dev.Names()
		if dev.IsUp() {
			psLog.W("Device is already opened")
			psLog.W(fmt.Sprintf("\ttype: %s", typ))
			psLog.W(fmt.Sprintf("\tname: %s (%s)", name, privName))
			return psErr.Error
		}
		if err := dev.Open(); err != psErr.OK {
			psLog.E(fmt.Sprintf("IDevice.Open() failed: %s", err))
			psLog.E(fmt.Sprintf("\ttype: %s", typ))
			psLog.E(fmt.Sprintf("\tname: %s (%s)", name, privName))
			return psErr.Error
		}
		dev.Up()
		psLog.I("Device was opened")
		psLog.I(fmt.Sprintf("\ttype: %s", typ))
		psLog.I(fmt.Sprintf("\tname: %s (%s)", name, privName))
	}
	return psErr.OK
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

func (p *ifaceRepo) Register(iface *Iface, dev ethernet.IDevice) psErr.E {
	for _, i := range p.ifaces {
		if i.Dev.Equal(dev) && i.Family == iface.Family {
			psLog.W(fmt.Sprintf("Interface is already registered: %s", i.Family))
			return psErr.Error
		}
	}

	p.ifaces = append(p.ifaces, iface)
	iface.Dev = dev

	name, privName := dev.Names()
	psLog.I("Interface was attached")
	psLog.I(fmt.Sprintf("\tip:     %s", iface.Unicast))
	psLog.I(fmt.Sprintf("\tdevice: %s (%s)", name, privName))

	return psErr.OK
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
	psLog.I("Route was registered")
	psLog.I(fmt.Sprintf("\tnetwork:  %s", route.Network))
	psLog.I(fmt.Sprintf("\tnetmask:  %s", route.Netmask))
	psLog.I(fmt.Sprintf("\tunicast:  %s", iface.Unicast))
	psLog.I(fmt.Sprintf("\tnext hop: %s", nextHop))
	psLog.I(fmt.Sprintf("\tdevice:   %s (%s)", name, privName))
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
	psLog.I("Default gateway was registered")
	psLog.I(fmt.Sprintf("\tnetwork:  %s", route.Network))
	psLog.I(fmt.Sprintf("\tnetmask:  %s", route.Netmask))
	psLog.I(fmt.Sprintf("\tunicast:  %s", iface.Unicast))
	psLog.I(fmt.Sprintf("\tnext hop: %s", nextHop))
	psLog.I(fmt.Sprintf("\tdevice:   %s (%s)", name, privName))
}

func init() {
	DeviceRepo = &deviceRepo{}
	IfaceRepo = &ifaceRepo{}
	RouteRepo = &routeRepo{}
}
