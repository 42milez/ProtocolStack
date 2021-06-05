package net

import (
	"fmt"
	psErr "github.com/42milez/ProtocolStack/src/error"
	psLog "github.com/42milez/ProtocolStack/src/log"
	"github.com/42milez/ProtocolStack/src/mw"
)

func Transmit(dst mw.EthAddr, payload []byte, typ mw.EthType, iface *mw.Iface) psErr.E {
	if !iface.Dev.IsUp() {
		psLog.E(fmt.Sprintf("Device %s is down", iface.Dev.Name()))
		return psErr.DeviceNotOpened
	}

	if len(payload) > int(iface.Dev.MTU()) {
		psLog.E(fmt.Sprintf("EthMessage is too long: mtu = %d, actual = %d", iface.Dev.MTU(), len(payload)))
		return psErr.PacketTooLong
	}

	if err := iface.Dev.Transmit(dst, payload, typ); err != psErr.OK {
		return psErr.Error
	}

	return psErr.OK
}
