package middleware

import (
	"github.com/42milez/ProtocolStack/src/ethernet"
	"github.com/42milez/ProtocolStack/src/network"
)

// An Iface is a single iface.
type Iface struct {
	Dev       ethernet.IDevice
	Family    network.AddrFamily
	Unicast   network.IP
	Netmask   network.IP
	Network   network.IP
	Broadcast network.IP
}

// GenIF generates Iface.
func GenIF(unicast string, netmask string) *Iface {
	iface := &Iface{
		Family:    network.FamilyV4,
		Unicast:   network.ParseIP(unicast),
		Netmask:   network.ParseIP(netmask),
		Broadcast: make(network.IP, 0),
		Network:   make(network.IP, 0),
	}

	//unicastUint32 := binary.BigEndian.Uint32(iface.Unicast)
	//netmaskUint32 := binary.BigEndian.Uint32(iface.Netmask)
	//broadcastUint32 := (unicastUint32 & netmaskUint32) | ^netmaskUint32

	//binary.BigEndian.PutUint32(iface.Broadcast, broadcastUint32)
	//binary.BigEndian.PutUint32(iface.Network, unicastUint32 & netmaskUint32)

	return iface
}
