package route

import (
	l "github.com/42milez/ProtocolStack/src/logger"
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
	l.I("route registered")
	l.I("\tnetwork:  %v ", route.Network.String())
	l.I("\tnetmask:  %v ", route.Netmask.String())
	l.I("\tunicast:  %v ", iface.Unicast.String())
	l.I("\tnext hop: %v ", nextHop.String())
	l.I("\tdevice:   %v (%v) ", iface.Dev.Name, iface.Dev.Priv.Name)
}

func RegisterDefaultGateway(iface *middleware.Iface, nextHop network.IP) {
	route := &Route{
		Network: network.V4Zero,
		Netmask: network.V4Zero,
		NextHop: nextHop,
		iface:   iface,
	}
	routes = append(routes, route)
	l.I("default gateway registered")
	l.I("\tnetwork:  %v ", route.Network.String())
	l.I("\tnetmask:  %v ", route.Netmask.String())
	l.I("\tunicast:  %v ", iface.Unicast.String())
	l.I("\tnext hop: %v ", nextHop.String())
	l.I("\tdevice:   %v (%v) ", iface.Dev.Name, iface.Dev.Priv.Name)
}
