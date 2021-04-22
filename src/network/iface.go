package network

import (
	"github.com/42milez/ProtocolStack/src/device"
)

// An Iface is a single interface.
type Iface struct {
	Dev *device.Device
	Family AddrFamily
	Unicast IP
	Netmask IP
	Network IP
	Broadcast IP
}

// GenIF generates Iface.
func GenIF(unicast string, netmask string) *Iface {
	iface := &Iface{
		Family:  FamilyV4,
		Unicast: ParseIP(unicast),
		Netmask: ParseIP(netmask),
		Broadcast: make(IP, 0),
		Network: make(IP, 0),
	}

	//unicastUint32 := binary.BigEndian.Uint32(iface.Unicast)
	//netmaskUint32 := binary.BigEndian.Uint32(iface.Netmask)
	//broadcastUint32 := (unicastUint32 & netmaskUint32) | ^netmaskUint32

	//binary.BigEndian.PutUint32(iface.Broadcast, broadcastUint32)
	//binary.BigEndian.PutUint32(iface.Network, unicastUint32 & netmaskUint32)

	return iface
}
