package network

import (
	psErr "github.com/42milez/ProtocolStack/src/error"
	"github.com/42milez/ProtocolStack/src/ethernet"
	psLog "github.com/42milez/ProtocolStack/src/log"
)

func InputHandler(packet *ethernet.Packet) psErr.Error {
	switch packet.Type {
	case ethernet.EthTypeArp:
		if err := ArpInputHandler(packet.Payload, packet.Dev); err.Code != psErr.OK {
			psLog.E("can't process arp packet: %s", err.Error())
			return psErr.Error{Code: psErr.CantProcess}
		}
	case ethernet.EthTypeIpv4:
		psLog.I("icmp packet received")
		return psErr.Error{Code: psErr.OK}
	default:
		psLog.E("▶ Unknown ether type: 0x%04x", uint16(packet.Type))
		return psErr.Error{Code: psErr.CantProcess}
	}
	return psErr.Error{Code: psErr.OK}
}

func OutputHandler(packet *ethernet.Packet) psErr.Error {
	return psErr.Error{Code: psErr.OK}
}
