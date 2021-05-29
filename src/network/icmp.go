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
const IcmpTypeEchoReply = 0
const IcmpTypeEcho = 8

// ICMP Type Numbers
// https://www.iana.org/assignments/icmp-parameters/icmp-parameters.xhtml#icmp-parameters-types

var icmpTypes = map[IcmpType]string{
	0: "Echo Reply",
	// 1-2: Unassigned
	3: "Destination Unreachable",
	4: "Source Quench (Deprecated)",
	5: "Redirect",
	6: "Alternate Host Address (Deprecated)",
	// 7: Unassigned
	8: "Echo",
	9: "Router Advertisement",
	10: "Router Solicitation",
	11: "Time Exceeded",
	12: "Parameter Problem",
	13: "Timestamp",
	14: "Timestamp Reply",
	15: "Information Request (Deprecated)",
	16: "Information Reply (Deprecated)",
	17: "Address Mask Request (Deprecated)",
	18: "Address Mask Reply (Deprecated)",
	19: "Reserved (for Security)",
	// 20-29: Reserved (for Robustness Experiment)
	30: "Traceroute (Deprecated)",
	31: "Datagram Conversion Error (Deprecated)",
	32: "Mobile Host Redirect (Deprecated)",
	33: "IPv6 Where-Are-You (Deprecated)",
	34: "IPv6 I-Am-Here (Deprecated)",
	35: "Mobile Registration Request (Deprecated)",
	36: "Mobile Registration Reply (Deprecated)",
	37: "Domain Name Request (Deprecated)",
	38: "Domain Name Reply (Deprecated)",
	39: "SKIP (Deprecated)",
	40: "Photuris",
	41: "ICMP messages utilized by experimental mobility protocols such as Seamoby",
	42: "Extended Echo Request",
	43: "Extended Echo Reply",
	// 44-252: Unassigned
	253: "RFC3692-style Experiment 1",
	254: "RFC3692-style Experiment 2",
	// 255: Reserved
}

type IcmpType uint8

type IcmpHeader struct {
	Type     IcmpType
	Code     uint8
	Checksum uint16
	Content  uint32
}

func IcmpInputHandler(payload []byte, dev ethernet.IDevice) psErr.E {
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
	payload[2] = 0 // assign 0 to Checksum field (16bit)
	payload[3] = 0
	if cs2 := Checksum(payload); cs2 != cs1 {
		psLog.E(fmt.Sprintf("Checksum mismatch: Expect = 0x%04x, Actual = 0x%04x", cs1, cs2))
		return psErr.ChecksumMismatch
	}

	psLog.I("Incoming ICMP packet")
	icmpDump(&hdr)

	switch hdr.Type {
	case IcmpTypeEcho:
		iface := IfaceRepo.Get(dev, V4AddrFamily)
		if iface.Unicast.EqualV4()
	}

	return psErr.OK
}

func icmpDump(hdr *IcmpHeader) {
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
