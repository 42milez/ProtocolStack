package network

import (
	psLog "github.com/42milez/ProtocolStack/src/log"
)

type Route struct {
	Network IP
	Netmask IP
	NextHop IP
	iface   *Iface
}

var routes []*Route

func init() {
	routes = make([]*Route, 0)
}

func RegisterRoute(network IP, nextHop IP, iface *Iface) {
	route := &Route{
		Network: network,
		Netmask: iface.Netmask,
		NextHop: nextHop,
		iface:   iface,
	}
	routes = append(routes, route)
	name, privName := iface.Dev.Names()
	psLog.I("route registered")
	psLog.I("\tnetwork:  %s", route.Network)
	psLog.I("\tnetmask:  %s", route.Netmask)
	psLog.I("\tunicast:  %s", iface.Unicast)
	psLog.I("\tnext hop: %s", nextHop)
	psLog.I("\tdevice:   %s (%s) ", name, privName)
}

func RegisterDefaultGateway(iface *Iface, nextHop IP) {
	route := &Route{
		Network: V4Zero,
		Netmask: V4Zero,
		NextHop: nextHop,
		iface:   iface,
	}
	routes = append(routes, route)
	name, privName := iface.Dev.Names()
	psLog.I("default gateway registered")
	psLog.I("\tnetwork:  %s", route.Network)
	psLog.I("\tnetmask:  %s", route.Netmask)
	psLog.I("\tunicast:  %s", iface.Unicast)
	psLog.I("\tnext hop: %s", nextHop)
	psLog.I("\tdevice:   %s (%s) ", name, privName)
}
