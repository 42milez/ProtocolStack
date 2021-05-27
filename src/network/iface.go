package network

import (
	"github.com/42milez/ProtocolStack/src/ethernet"
)

// An Iface is a single iface.
type Iface struct {
	Family    AddrFamily
	Unicast   IP
	Netmask   IP
	Broadcast IP
	Dev       ethernet.IDevice
}
