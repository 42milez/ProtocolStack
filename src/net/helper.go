package net

// GenIface generates Iface.
func GenIface(unicast string, netmask string, broadcast string) *Iface {
	iface := &Iface{
		Family:    V4AddrFamily,
		Unicast:   ParseIP(unicast),
		Netmask:   ParseIP(netmask),
		Broadcast: ParseIP(broadcast),
	}
	return iface
}
