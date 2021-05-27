package network

import (
	psErr "github.com/42milez/ProtocolStack/src/error"
	"github.com/42milez/ProtocolStack/src/ethernet"
	psLog "github.com/42milez/ProtocolStack/src/log"
)

func InputHandler(packet *ethernet.Packet) psErr.E {
	switch packet.Type {
	case ethernet.EthTypeArp:
		if err := ArpInputHandler(packet.Payload, packet.Dev); err != psErr.OK {
			psLog.E("ArpInputHandler() failed: %s", err)
			return psErr.Error
		}
	case ethernet.EthTypeIpv4:
		psLog.I("icmp packet received")
		return psErr.OK
	default:
		psLog.E("Unknown ether type: 0x%04x", uint16(packet.Type))
		return psErr.Error
	}
	return psErr.OK
}

func OutputHandler(packet *ethernet.Packet) psErr.E {
	return psErr.OK
}
