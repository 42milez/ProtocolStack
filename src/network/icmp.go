package network

import (
	"bytes"
	"encoding/binary"
	"fmt"
	psErr "github.com/42milez/ProtocolStack/src/error"
	"github.com/42milez/ProtocolStack/src/ethernet"
	psLog "github.com/42milez/ProtocolStack/src/log"
)

const IcmpHeaderSize = 8 // byte
const IcmpTypeEchoReply = 0x00
const IcmpTypeEcho = 0x08

func IcmpReceive(payload []byte, dst [V4AddrLen]byte, src [V4AddrLen]byte, dev ethernet.IDevice) psErr.E {
	if len(payload) < IcmpHeaderSize {
		psLog.E(fmt.Sprintf("ICMP header length is too short: %d bytes", len(payload)))
		return psErr.InvalidPacket
	}

	buf := bytes.NewBuffer(payload)
	hdr := IcmpHeader{}
	if err := binary.Read(buf, binary.BigEndian, &hdr); err != nil {
		psLog.E(fmt.Sprintf("binary.Read() failed: %s", err))
		return psErr.Error
	}

	cs1 := uint16(payload[2])<<8 | uint16(payload[3])
	payload[2] = 0x00 // assign 0 to Checksum field (16bit)
	payload[3] = 0x00
	if cs2 := Checksum(payload); cs2 != cs1 {
		psLog.E(fmt.Sprintf("Checksum mismatch: Expect = 0x%04x, Actual = 0x%04x", cs1, cs2))
		return psErr.ChecksumMismatch
	}

	psLog.I("Incoming ICMP packet")
	icmpHdrDump(&hdr)

	switch hdr.Type {
	case IcmpTypeEcho:
		s := IP(src[:])
		d := IP(dst[:])
		iface := IfaceRepo.Get(dev, V4AddrFamily)
		if !iface.Unicast.EqualV4(dst) {
			d = iface.Unicast
		}
		if err := IcmpSend(IcmpTypeEchoReply, hdr.Code, hdr.Content, payload[IcmpHeaderSize:], d, s); err != psErr.OK {
			psLog.E(fmt.Sprintf("IcmpSend() failed: %s", err))
			return psErr.Error
		}
	}

	return psErr.OK
}

func IcmpSend(typ IcmpType, code uint8, content uint32, payload []byte, dst IP, src IP) psErr.E {
	hdr := IcmpHeader{
		Type:    typ,
		Code:    code,
		Content: content,
	}

	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, &hdr); err != nil {
		psLog.E(fmt.Sprintf("binary.Write() failed: %s", err))
		return psErr.Error
	}

	if err := binary.Write(buf, binary.BigEndian, &payload); err != nil {
		psLog.E(fmt.Sprintf("binary.Write() failed: %s", err))
		return psErr.Error
	}

	packet := buf.Bytes()
	checksum := Checksum(packet)
	packet[2] = uint8((checksum & 0xff00) >> 8)
	packet[3] = uint8(checksum & 0x00ff)

	psLog.I("Outgoing ICMP packet")
	hdr.Checksum = checksum
	icmpHdrDump(&hdr)

	if err := IpSend(ProtoNumICMP, packet, dst, src); err != psErr.OK {
		psLog.E(fmt.Sprintf("IpSend() failed: %s", err))
		return psErr.Error
	}

	return psErr.OK
}

func icmpHdrDump(hdr *IcmpHeader) {
	psLog.I(fmt.Sprintf("\ttype:     %d (%s)", hdr.Type, icmpTypes[hdr.Type]))
	psLog.I(fmt.Sprintf("\tcode:     %d", hdr.Code))
	psLog.I(fmt.Sprintf("\tchecksum: 0x%04x", hdr.Checksum))
	switch hdr.Type {
	case IcmpTypeEchoReply:
	case IcmpTypeEcho:
		psLog.I(fmt.Sprintf("\tid:       %d", (hdr.Content&0xffff0000)>>16))
		psLog.I(fmt.Sprintf("\tseq:      %d", hdr.Content&0x0000ffff))
	default:
		psLog.I(fmt.Sprintf("\tcontent:  %02x %02x %02x %02x",
			uint8((hdr.Content&0xf000)>>24),
			uint8((hdr.Content&0x0f00)>>16),
			uint8((hdr.Content&0x00f0)>>8),
			uint8(hdr.Content&0x000f)))
	}
}
