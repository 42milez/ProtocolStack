package route

import (
	psLog "github.com/42milez/ProtocolStack/src/log"
	"github.com/42milez/ProtocolStack/src/network"
)

type Route struct {
	Network network.IP
	Netmask network.IP
	NextHop network.IP
	iface   *network.Iface
}

var routes []*Route

func init() {
	routes = make([]*Route, 0)
}

func RegisterRoute(network network.IP, nextHop network.IP, iface *network.Iface) {
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

func RegisterDefaultGateway(iface *network.Iface, nextHop network.IP) {
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
	_, name1, name2 := iface.Dev.Info()
	psLog.I("\tdevice:   %v (%v) ", name1, name2)
}
