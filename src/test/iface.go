package test

import "github.com/42milez/ProtocolStack/src/mw"

var IfaceBuilder ifaceBuilder

type ifaceBuilder struct{}

func (ifaceBuilder) Default() *mw.Iface {
	return &mw.Iface{
		Family:    mw.V4AddrFamily,
		Unicast:   mw.IP{192, 0, 2, 2},
		Netmask:   mw.IP{255, 255, 255, 0},
		Broadcast: mw.IP{192, 0, 2, 255},
	}
}

func init() {
	IfaceBuilder = ifaceBuilder{}
}
