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
	psLog.I("route registered")
	psLog.I("\tnetwork:  %v ", route.Network.String())
	psLog.I("\tnetmask:  %v ", route.Netmask.String())
	psLog.I("\tunicast:  %v ", iface.Unicast.String())
	psLog.I("\tnext hop: %v ", nextHop.String())
	_, name1, name2 := iface.Dev.Info()
	psLog.I("\tdevice:   %v (%v) ", name1, name2)
}

func RegisterDefaultGateway(iface *Iface, nextHop IP) {
	route := &Route{
		Network: V4Zero,
		Netmask: V4Zero,
		NextHop: nextHop,
		iface:   iface,
	}
	routes = append(routes, route)
	psLog.I("default gateway registered")
	psLog.I("\tnetwork:  %v ", route.Network.String())
	psLog.I("\tnetmask:  %v ", route.Netmask.String())
	psLog.I("\tunicast:  %v ", iface.Unicast.String())
	psLog.I("\tnext hop: %v ", nextHop.String())
	_, name1, name2 := iface.Dev.Info()
	psLog.I("\tdevice:   %v (%v) ", name1, name2)
}
