package network

import (
	"github.com/42milez/ProtocolStack/src/ethernet"
)

// An Iface is a single iface.
type Iface struct {
	Family    AddrFamily
	Netmask   IP
	Broadcast IP
	Unicast   IP
	Dev       ethernet.IDevice
}
