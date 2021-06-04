package net

import (
	"fmt"
	psErr "github.com/42milez/ProtocolStack/src/error"
	"github.com/42milez/ProtocolStack/src/eth"
	psLog "github.com/42milez/ProtocolStack/src/log"
)

func InputHandler(packet *eth.Packet) psErr.E {
	switch packet.Type {
	case eth.ARP:
		if err := ARP.Receive(packet.Content, packet.Dev); err != psErr.OK {
			return psErr.Error
		}
	case eth.EthTypeIpv4:
		if err := IpReceive(packet.Content, packet.Dev); err != psErr.OK {
			return psErr.Error
		}
		return psErr.OK
	default:
		psLog.E(fmt.Sprintf("Unknown ether type: 0x%04x", uint16(packet.Type)))
		return psErr.Error
	}
	return psErr.OK
}

func OutputHandler(packet *eth.Packet) psErr.E {
	return psErr.OK
}
