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
			psLog.W(fmt.Sprintf("\ttype: %s", d.DevType()))
			psLog.W(fmt.Sprintf("\tname: %s (%s)", d.DevName(), d.PrivDevName()))
			return psErr.Error
		}
	}
	p.devices = append(p.devices, dev)
	psLog.I("Device was registered")
	psLog.I(fmt.Sprintf("\ttype:      %s", dev.DevType()))
	psLog.I(fmt.Sprintf("\tname:      %s (%s)", dev.DevName(), dev.PrivDevName()))
	psLog.I(fmt.Sprintf("\taddr:      %s", dev.EthAddr()))
	psLog.I(fmt.Sprintf("\tbroadcast: %s", dev.BroadcastEthAddr()))
	psLog.I(fmt.Sprintf("\tpeer:      %s", dev.PeerEthAddr()))
	return psErr.OK
}

func (p *deviceRepo) Up() psErr.E {
	for _, dev := range p.devices {
		if dev.IsUp() {
			psLog.W("Device is already opened")
			psLog.W(fmt.Sprintf("\ttype: %s", dev.DevType()))
			psLog.W(fmt.Sprintf("\tname: %s (%s)", dev.DevName(), dev.PrivDevName()))
			return psErr.Error
		}
		if err := dev.Open(); err != psErr.OK {
			psLog.E(fmt.Sprintf("IDevice.Open() failed: %s", err))
			psLog.E(fmt.Sprintf("\ttype: %s", dev.DevType()))
			psLog.E(fmt.Sprintf("\tname: %s (%s)", dev.DevName(), dev.PrivDevName()))
			return psErr.Error
		}
		dev.Up()
		psLog.I("Device was opened")
		psLog.I(fmt.Sprintf("\ttype: %s", dev.DevType()))
		psLog.I(fmt.Sprintf("\tname: %s (%s)", dev.DevName(), dev.PrivDevName()))
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

	psLog.I("Interface was attached")
	psLog.I(fmt.Sprintf("\tip:     %s", iface.Unicast))
	psLog.I(fmt.Sprintf("\tdevice: %s (%s)", dev.DevName(), dev.PrivDevName()))

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
	psLog.I("Route was registered")
	psLog.I(fmt.Sprintf("\tnetwork:  %s", route.Network))
	psLog.I(fmt.Sprintf("\tnetmask:  %s", route.Netmask))
	psLog.I(fmt.Sprintf("\tunicast:  %s", iface.Unicast))
	psLog.I(fmt.Sprintf("\tnext hop: %s", nextHop))
	psLog.I(fmt.Sprintf("\tdevice:   %s (%s)", iface.Dev.DevName(), iface.Dev.PrivDevName()))
}

func (p *routeRepo) RegisterDefaultGateway(iface *Iface, nextHop IP) {
	route := &route{
		Network: V4Zero,
		Netmask: V4Zero,
		NextHop: nextHop,
		iface:   iface,
	}
	p.routes = append(p.routes, route)
	psLog.I("Default gateway was registered")
	psLog.I(fmt.Sprintf("\tnetwork:  %s", route.Network))
	psLog.I(fmt.Sprintf("\tnetmask:  %s", route.Netmask))
	psLog.I(fmt.Sprintf("\tunicast:  %s", iface.Unicast))
	psLog.I(fmt.Sprintf("\tnext hop: %s", nextHop))
	psLog.I(fmt.Sprintf("\tdevice:   %s (%s)", iface.Dev.DevName(), iface.Dev.PrivDevName()))
}

func init() {
	DeviceRepo = &deviceRepo{}
	IfaceRepo = &ifaceRepo{}
	RouteRepo = &routeRepo{}
}
