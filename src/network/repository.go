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
			psLog.W("Device is already registered")
			psLog.W(fmt.Sprintf("\ttype: %s", d.Type()))
			psLog.W(fmt.Sprintf("\tname: %s (%s)", d.Name(), d.Priv().Name))
			return psErr.Error
		}
	}
	p.devices = append(p.devices, dev)
	psLog.I("Device was registered")
	psLog.I(fmt.Sprintf("\ttype:      %s", dev.Type()))
	psLog.I(fmt.Sprintf("\tname:      %s (%s)", dev.Name(), dev.Priv().Name))
	psLog.I(fmt.Sprintf("\taddr:      %s", dev.Addr()))
	psLog.I(fmt.Sprintf("\tbroadcast: %s", dev.Broadcast()))
	psLog.I(fmt.Sprintf("\tpeer:      %s", dev.Peer()))
	return psErr.OK
}

func (p *deviceRepo) Up() psErr.E {
	for _, dev := range p.devices {
		if dev.IsUp() {
			psLog.W("Device is already opened")
			psLog.W(fmt.Sprintf("\ttype: %s", dev.Type()))
			psLog.W(fmt.Sprintf("\tname: %s (%s)", dev.Name(), dev.Priv().Name))
			return psErr.Error
		}
		if err := dev.Open(); err != psErr.OK {
			psLog.E(fmt.Sprintf("IDevice.Open() failed: %s", err))
			psLog.E(fmt.Sprintf("\ttype: %s", dev.Type()))
			psLog.E(fmt.Sprintf("\tname: %s (%s)", dev.Name(), dev.Priv().Name))
			return psErr.Error
		}
		dev.Up()
		psLog.I("Device was opened")
		psLog.I(fmt.Sprintf("\ttype: %s", dev.Type()))
		psLog.I(fmt.Sprintf("\tname: %s (%s)", dev.Name(), dev.Priv().Name))
	}
	return psErr.OK
}

type ifaceRepo struct {
	ifaces []*Iface
}

func (p *ifaceRepo) Get(ip IP) *Iface {
	for _, v := range p.ifaces {
		if v.Unicast.Equal(ip) {
			return v
		}
	}
	return nil
}

func (p *ifaceRepo) Lookup(dev ethernet.IDevice, family AddrFamily) *Iface {
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

	psLog.I("Interface was attached")
	psLog.I(fmt.Sprintf("\tip:     %s", iface.Unicast))
	psLog.I(fmt.Sprintf("\tdevice: %s (%s)", dev.Name(), dev.Priv().Name))

	return psErr.OK
}

type Route struct {
	Network IP
	Netmask IP
	NextHop IP
	Iface   *Iface
}

type routeRepo struct {
	routes []*Route
}

func (p *routeRepo) Get(ip IP) *Route {
	var ret *Route
	for _, route := range p.routes {
		if ip.Mask(route.Netmask).Equal(route.Network) {
			// Longest prefix match
			// https://en.wikipedia.org/wiki/Longest_prefix_match
			if ret == nil || longestIP(ret.Netmask, route.Netmask).Equal(route.Netmask) {
				ret = route
			}
		}
	}
	return ret
}

func (p *routeRepo) Register(network IP, nextHop IP, iface *Iface) {
	route := &Route{
		Network: network,
		Netmask: iface.Netmask,
		NextHop: nextHop,
		Iface:   iface,
	}
	p.routes = append(p.routes, route)
	psLog.I("Route was registered")
	psLog.I(fmt.Sprintf("\tnetwork:  %s", route.Network))
	psLog.I(fmt.Sprintf("\tnetmask:  %s", route.Netmask))
	psLog.I(fmt.Sprintf("\tunicast:  %s", iface.Unicast))
	psLog.I(fmt.Sprintf("\tnext hop: %s", nextHop))
	psLog.I(fmt.Sprintf("\tdevice:   %s (%s)", iface.Dev.Name(), iface.Dev.Priv().Name))
}

func (p *routeRepo) RegisterDefaultGateway(iface *Iface, nextHop IP) {
	route := &Route{
		Network: V4Zero,
		Netmask: V4Zero,
		NextHop: nextHop,
		Iface:   iface,
	}
	p.routes = append(p.routes, route)
	psLog.I("Default gateway was registered")
	psLog.I(fmt.Sprintf("\tnetwork:  %s", route.Network))
	psLog.I(fmt.Sprintf("\tnetmask:  %s", route.Netmask))
	psLog.I(fmt.Sprintf("\tunicast:  %s", iface.Unicast))
	psLog.I(fmt.Sprintf("\tnext hop: %s", nextHop))
	psLog.I(fmt.Sprintf("\tdevice:   %s (%s)", iface.Dev.Name(), iface.Dev.Priv().Name))
}

func init() {
	DeviceRepo = &deviceRepo{}
	IfaceRepo = &ifaceRepo{}
	RouteRepo = &routeRepo{}
}
