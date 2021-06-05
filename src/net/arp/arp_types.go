package arp

import (
	"fmt"
	"github.com/42milez/ProtocolStack/src/mw"
)

const PacketLen = 28 // byte
const Request Opcode = 0x0001
const Reply Opcode = 0x0002

type ArpProtoAddr [mw.V4AddrLen]byte

func (p ArpProtoAddr) String() string {
	return fmt.Sprintf("%d.%d.%d.%d", p[0], p[1], p[2], p[3])
}
