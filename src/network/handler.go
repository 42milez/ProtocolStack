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
	default:
		psLog.E("unknown ether type")
		return psErr.Error{Code: psErr.CantProcess}
	}
	return psErr.Error{Code: psErr.OK}
}

func OutputHandler(packet *ethernet.Packet) psErr.Error {
	return psErr.Error{Code: psErr.OK}
}
