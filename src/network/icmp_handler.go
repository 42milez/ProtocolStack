package network

import (
	"bytes"
	"encoding/binary"
	"fmt"
	psErr "github.com/42milez/ProtocolStack/src/error"
	"github.com/42milez/ProtocolStack/src/ethernet"
	psLog "github.com/42milez/ProtocolStack/src/log"
)

func IcmpInputHandler(payload []byte, dev ethernet.IDevice) psErr.E {
	if len(payload) < IcmpHeaderSize {
		psLog.E(fmt.Sprintf("ICMP header length is too short: %d bytes", len(payload)))
		return psErr.InvalidPacket
	}

	buf := new(bytes.Buffer)
	packet := IcmpPacket{}
	if err := binary.Read(buf, binary.BigEndian, &packet); err != nil {
		psLog.E(fmt.Sprintf("binary.Read() failed: %s", err))
		return psErr.Error
	}

	cs1 := uint16(payload[2])<<8 | uint16(payload[3])
	payload[2] = 0 // assign 0 to Checksum field (16bit)
	payload[3] = 0
	if cs2 := Checksum(payload); cs2 != cs1 {
		psLog.E(fmt.Sprintf("Checksum mismatch: Expect = 0x%04x, Actual = 0x%04x", cs1, cs2))
		return psErr.ChecksumMismatch
	}

	psLog.I("Incoming ICMP packet")
	icmpDump(&packet)

	return psErr.OK
}

func icmpDump(packet *IcmpPacket) {
	psLog.I(fmt.Sprintf("type:     %d", packet.Type))
	psLog.I(fmt.Sprintf("code:     %d", packet.Code))
	psLog.I(fmt.Sprintf("checksum: %d", packet.Checksum))
	psLog.I(fmt.Sprintf("content:  0x%02x 0x%02x 0x%02x 0x%02x",
		uint8((packet.Content&0xf000)>>24),
		uint8((packet.Content&0x0f00)>>16),
		uint8((packet.Content&0x00f0)>>8),
		uint8(packet.Content&0x000f)))
}
