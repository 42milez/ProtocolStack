package middleware

import (
	"github.com/42milez/ProtocolStack/src/ethernet"
	"github.com/42milez/ProtocolStack/src/network"
)

// An Iface is a single iface.
type Iface struct {
	Family    network.AddrFamily
	Network   network.IP
	Netmask   network.IP
	Broadcast network.IP
	Unicast   network.IP
	Dev       ethernet.IDevice
}
