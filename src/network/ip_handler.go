package network

import (
	"bytes"
	"encoding/binary"
	"fmt"
	psErr "github.com/42milez/ProtocolStack/src/error"
	"github.com/42milez/ProtocolStack/src/ethernet"
	psLog "github.com/42milez/ProtocolStack/src/log"
)

// INTERNET PROTOCOL
// https://datatracker.ietf.org/doc/html/rfc791#page-13
// The number 576 is selected to allow a reasonable sized data block to be transmitted in addition to the required
// header information. For example, this size allows a data block of 512 octets plus 64 header octets to fit in a
// datagram. The maximal internet header is 60 octets, and a typical internet header is 20 octets, allowing a margin for
// headers of higher level protocols.

const IpHeaderSizeMin = 20 // bytes
const IpHeaderSizeMax = 60 // bytes

const ProtoNumICMP = 1
const ProtoNumTCP = 6
const ProtoNumUDP = 17

// ASSIGNED INTERNET PROTOCOL NUMBERS
// https://datatracker.ietf.org/doc/html/rfc790#page-6

var protocolNumbers = map[ProtocolNumber]string{
	// 0: Reserved
	1:  "ICMP",
	3:  "Gateway-to-Gateway",
	4:  "CMCC Gateway Monitoring Message",
	5:  "ST",
	6:  "TCP",
	7:  "UCL",
	9:  "Secure",
	10: "BBN RCC Monitoring",
	11: "NVP",
	12: "PUP",
	13: "Pluribus",
	14: "Telenet",
	15: "XNET",
	16: "Chaos",
	17: "User Datagram",
	18: "Multiplexing",
	19: "DCN",
	20: "TAC Monitoring",
	// 21-62: Unassigned
	63: "any local network",
	64: "SATNET and Backroom EXPAK",
	65: "MIT Subnet Support",
	// 66-68: Unassigned
	69: "SATNET Monitoring",
	71: "Internet Packet Core Utility",
	// 72-75: Unassigned
	76: "Backroom SATNET Monitoring",
	78: "WIDEBAND Monitoring",
	79: "WIDEBAND EXPAK",
	// 80-254: Unassigned
	// 255: Reserved
}

// Internet Header Format
// https://datatracker.ietf.org/doc/html/rfc791#section-3.1

type ProtocolNumber uint8

type IpHeader struct {
	VHL      uint8
	TOS      uint8
	TotalLen uint16
	ID       uint16
	Offset   uint16
	TTL      uint8
	Protocol ProtocolNumber
	Checksum uint16
	Src      [V4AddrLen]byte
	Dst      [V4AddrLen]byte
	Options  [0]byte
}

func IpInputHandler(payload []byte, dev ethernet.IDevice) psErr.E {
	packetLen := len(payload)

	if packetLen < IpHeaderSizeMin {
		psLog.E(fmt.Sprintf("IP packet length is too short: %d bytes", packetLen))
		return psErr.InvalidPacket
	}

	buf := bytes.NewBuffer(payload)
	hdr := IpHeader{}
	if err := binary.Read(buf, binary.BigEndian, &hdr); err != nil {
		psLog.E(fmt.Sprintf("binary.Read() failed: %s", err))
		return psErr.Error
	}

	if version := hdr.VHL >> 4; version != 4 {
		psLog.E(fmt.Sprintf("IP version %d is not supported", version))
		return psErr.InvalidProtocolVersion
	}

	hdrLen := int(hdr.VHL&0x0f) * 4
	if packetLen < hdrLen {
		psLog.E(fmt.Sprintf("IP packet length is too short: IHL = %d, Actual Packet Size = %d", hdrLen, packetLen))
		return psErr.InvalidPacket
	}

	if totalLen := int(hdr.TotalLen); packetLen < totalLen {
		psLog.E(fmt.Sprintf("IP packet length is too short: Total Length = %d, Actual Length = %d", totalLen, packetLen))
		return psErr.InvalidPacket
	}

	if hdr.TTL == 0 {
		psLog.E("TTL expired")
		return psErr.TtlExpired
	}

	cs1 := uint16(payload[10])<<8 | uint16(payload[11])
	payload[10] = 0 // assign 0 to Header Checksum field (16bit)
	payload[11] = 0
	if cs2 := Checksum(payload); cs2 != cs1 {
		psLog.E(fmt.Sprintf("Checksum mismatch: Expect = 0x%04x, Actual = 0x%04x", cs1, cs2))
		return psErr.ChecksumMismatch
	}

	iface := IfaceRepo.Get(dev, V4AddrFamily)
	if iface == nil {
		psLog.E(fmt.Sprintf("Interface for %s is not registered", dev.DevName()))
		return psErr.InterfaceNotFound
	}

	if !iface.Unicast.EqualV4(hdr.Dst) {
		if !iface.Broadcast.EqualV4(hdr.Dst) && V4Broadcast.EqualV4(hdr.Dst) {
			psLog.I("Ignored IP packet (It was sent to different address)")
			return psErr.OK
		}
	}

	psLog.I("Incoming IP packet")
	ipDump(&hdr)

	switch hdr.Protocol {
	case ProtoNumICMP:
		if err := IcmpReceive(payload[hdrLen:], hdr.Src, hdr.Dst, dev); err != psErr.OK {
			psLog.E(fmt.Sprintf("IcmpInputHandler() failed: %s", err))
			return psErr.Error
		}
	case ProtoNumTCP:
		psLog.E("Currently NOT support TCP")
		return psErr.Error
	case ProtoNumUDP:
		psLog.E("Currently NOT support UDP")
		return psErr.Error
	default:
		psLog.E(fmt.Sprintf("Unsupported protocol: %d", hdr.Protocol))
		return psErr.UnsupportedProtocol
	}

	return psErr.OK
}

func ipDump(hdr *IpHeader) {
	psLog.I(fmt.Sprintf("\tversion:             IPv%d", hdr.VHL>>4))
	psLog.I(fmt.Sprintf("\tihl:                 %d", hdr.VHL&0x0f))
	psLog.I(fmt.Sprintf("\ttype of service:     0b%08b", hdr.TOS))
	psLog.I(fmt.Sprintf("\ttotal length:        %d bytes (payload: %d bytes)", hdr.TotalLen, hdr.TotalLen-uint16(4*(hdr.VHL&0x0f))))
	psLog.I(fmt.Sprintf("\tid:                  %d", hdr.ID))
	psLog.I(fmt.Sprintf("\tflags:               0b%03b", (hdr.Offset&0xefff)>>13))
	psLog.I(fmt.Sprintf("\tfragment offset:     %d", hdr.Offset&0x1fff))
	psLog.I(fmt.Sprintf("\tttl:                 %d", hdr.TTL))
	psLog.I(fmt.Sprintf("\tprotocol:            %d (%s)", hdr.Protocol, protocolNumbers[hdr.Protocol]))
	psLog.I(fmt.Sprintf("\tchecksum:            0x%04x", hdr.Checksum))
	psLog.I(fmt.Sprintf("\tsource address:      %d.%d.%d.%d", hdr.Src[0], hdr.Src[1], hdr.Src[2], hdr.Src[3]))
	psLog.I(fmt.Sprintf("\tdestination address: %d.%d.%d.%d", hdr.Dst[0], hdr.Dst[1], hdr.Dst[2], hdr.Dst[3]))
}
