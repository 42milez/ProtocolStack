package route

import (
	psLog "github.com/42milez/ProtocolStack/src/log"
	"github.com/42milez/ProtocolStack/src/middleware"
	"github.com/42milez/ProtocolStack/src/network"
)

type Route struct {
	Network network.IP
	Netmask network.IP
	NextHop network.IP
	iface   *middleware.Iface
}

var routes []*Route

func init() {
	routes = make([]*Route, 0)
}

func Register(iface *middleware.Iface, nextHop network.IP) {
	route := &Route{
		Network: iface.Network,
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
	psLog.I("\tdevice:   %v (%v) ", iface.Dev.Name, iface.Dev.Priv.Name)
}

func RegisterDefaultGateway(iface *middleware.Iface, nextHop network.IP) {
	route := &Route{
		Network: network.V4Zero,
		Netmask: network.V4Zero,
		NextHop: nextHop,
		iface:   iface,
	}
	routes = append(routes, route)
	psLog.I("default gateway registered")
	psLog.I("\tnetwork:  %v ", route.Network.String())
	psLog.I("\tnetmask:  %v ", route.Netmask.String())
	psLog.I("\tunicast:  %v ", iface.Unicast.String())
	psLog.I("\tnext hop: %v ", nextHop.String())
	psLog.I("\tdevice:   %v (%v) ", iface.Dev.Name, iface.Dev.Priv.Name)
}
