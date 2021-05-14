package route

import (
	"github.com/42milez/ProtocolStack/src/middleware"
	"github.com/42milez/ProtocolStack/src/network"
	"log"
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
	log.Printf(
		"route added: network=%s, netmask=%s, iface=%s, nextHop: %s, dev=%s",
		route.Network.String(),
		route.Netmask.String(),
		iface.Unicast.String(),
		nextHop.String(),
		iface.Dev.Name,
	)
}

func RegisterDefaultGateway(iface *middleware.Iface, nextHop network.IP) {
	route := &Route{
		Network: network.V4Zero,
		Netmask: network.V4Zero,
		NextHop: nextHop,
		iface:   iface,
	}
	routes = append(routes, route)
	log.Printf(
		"gateway added: network=%s, netmask=%s, iface=%s, nextHop: %s, dev=%s",
		route.Network.String(),
		route.Netmask.String(),
		iface.Unicast.String(),
		nextHop.String(),
		iface.Dev.Name,
	)
}
