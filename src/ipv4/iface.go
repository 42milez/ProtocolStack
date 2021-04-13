package ipv4

import (
	"encoding/binary"
	"github.com/42milez/ProtocolStack/src/device"
)

type Iface struct {
	Dev *device.Device
	Family AddrFamily
	Unicast IP
	Netmask IP
	Broadcast IP
}

var ifaces []*Iface

func GenIF(unicast string, netmask string, dev *device.Device) *Iface {
	iface := &Iface{
		Dev: dev,
		Family: V4,
		Unicast: ParseIP(unicast),
		Netmask: ParseIP(netmask),
	}

	unicastUint32 := binary.BigEndian.Uint32(iface.Unicast)
	netmaskUint32 := binary.BigEndian.Uint32(iface.Netmask)
	broadcastUint32 := (unicastUint32 & netmaskUint32) | ^netmaskUint32
	binary.BigEndian.PutUint32(iface.Broadcast, broadcastUint32)

	return iface
}
