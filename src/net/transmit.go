package net

import (
	"fmt"
	psErr "github.com/42milez/ProtocolStack/src/error"
	"github.com/42milez/ProtocolStack/src/eth"
	psLog "github.com/42milez/ProtocolStack/src/log"
)

func Transmit(dst eth.Addr, payload []byte, typ eth.Type, iface *Iface) psErr.E {
	if !iface.Dev.IsUp() {
		psLog.E(fmt.Sprintf("Device %s is down", iface.Dev.Name()))
		return psErr.DeviceNotOpened
	}

	if len(payload) > int(iface.Dev.MTU()) {
		psLog.E(fmt.Sprintf("Packet is too long: mtu = %d, actual = %d", iface.Dev.MTU(), len(payload)))
		return psErr.PacketTooLong
	}

	if err := iface.Dev.Transmit(dst, payload, typ); err != psErr.OK {
		return psErr.Error
	}

	return psErr.OK
}
