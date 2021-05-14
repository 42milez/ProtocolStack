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
	log.Println("registered a route")
	log.Printf("\tnetwork:  %v\n", route.Network.String())
	log.Printf("\tnetmask:  %v\n", route.Netmask.String())
	log.Printf("\tunicast:  %v\n", iface.Unicast.String())
	log.Printf("\tnext hop: %v\n", nextHop.String())
	log.Printf("\tdevice:   %v (%v)\n", iface.Dev.Name, iface.Dev.Priv.Name)
}

func RegisterDefaultGateway(iface *middleware.Iface, nextHop network.IP) {
	route := &Route{
		Network: network.V4Zero,
		Netmask: network.V4Zero,
		NextHop: nextHop,
		iface:   iface,
	}
	routes = append(routes, route)
	log.Println("registered a default gateway")
	log.Printf("\tnetwork:  %v\n", route.Network.String())
	log.Printf("\tnetmask:  %v\n", route.Netmask.String())
	log.Printf("\tunicast:  %v\n", iface.Unicast.String())
	log.Printf("\tnext hop: %v\n", nextHop.String())
	log.Printf("\tdevice:   %v (%v)\n", iface.Dev.Name, iface.Dev.Priv.Name)
}
