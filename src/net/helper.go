package net

import "github.com/42milez/ProtocolStack/src/mw"

// GenIface generates Iface.
func GenIface(unicast string, netmask string, broadcast string) *mw.Iface {
	iface := &mw.Iface{
		Family:    mw.V4AddrFamily,
		Unicast:   mw.ParseIP(unicast),
		Netmask:   mw.ParseIP(netmask),
		Broadcast: mw.ParseIP(broadcast),
	}
	return iface
}
