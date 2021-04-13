package network

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/42milez/ProtocolStack/src/device"
	"log"
)

// An Iface is a single interface.
type Iface struct {
	Dev *device.Device
	Family AddrFamily
	Unicast IP
	Netmask IP
	Broadcast IP
}

// GenIF generates Iface.
func GenIF(unicast string, netmask string) *Iface {
	iface := &Iface{
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

// AttachIF attaches an Iface to device.Device.
func AttachIF(iface *Iface, dev *device.Device) error {
	iface.Dev = dev
	for _, v := range dev.Ifaces {
		if v.Family == iface.Family {
			return errors.New(fmt.Sprintf("%s is already exists", v.Family.String()))
		}
	}
	log.Printf("interface attached: iface=%s, dev=%s", iface.Unicast.String(), iface.Dev.Name)
	return nil
}
