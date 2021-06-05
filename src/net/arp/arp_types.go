package arp

import (
	"fmt"
	"github.com/42milez/ProtocolStack/src/mw"
)

const ArpOpRequest Opcode = 0x0001
const ArpOpReply Opcode = 0x0002
const ArpPacketLen = 28 // byte

type ArpProtoAddr [mw.V4AddrLen]byte

func (p ArpProtoAddr) String() string {
	return fmt.Sprintf("%d.%d.%d.%d", p[0], p[1], p[2], p[3])
}
