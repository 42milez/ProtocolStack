package route

import (
	"github.com/42milez/ProtocolStack/src/network"
	"log"
)

type Route struct {
	Network network.IP
	Netmask network.IP
	NextHop network.IP
	iface *network.Iface
}

var routes []*Route

func Register(network network.IP, netmask network.IP, nextHop network.IP, iface *network.Iface) {
	route := &Route{
		Network: network,
		Netmask: netmask,
		NextHop: nextHop,
		iface: iface,
	}
	routes = append(routes, route)

	log.Printf(
		"new route added: network=%s, netmask=%s, nextHop: %s, iface=%s, dev=%s",
		network.String(),
		netmask.String(),
		nextHop.String(),
		iface.Unicast.String(),
		iface.Dev.Name,
	)
}
