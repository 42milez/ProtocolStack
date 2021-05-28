package network

import (
	"bytes"
	"encoding/binary"
	"fmt"
	psErr "github.com/42milez/ProtocolStack/src/error"
	"github.com/42milez/ProtocolStack/src/ethernet"
	psLog "github.com/42milez/ProtocolStack/src/log"
)

const IpVersionV4 = 0x04

// INTERNET PROTOCOL
// https://datatracker.ietf.org/doc/html/rfc791#page-13
// The number 576 is selected to allow a reasonable sized data block to be transmitted in addition to the required
// header information. For example, this size allows a data block of 512 octets plus 64 header octets to fit in a
// datagram. The maximal internet header is 60 octets, and a typical internet header is 20 octets, allowing a margin for
// headers of higher level protocols.

const IpHeaderSizeMin = 20 // bytes
const IpHeaderSizeMax = 60 // bytes

// Internet Header Format
// https://datatracker.ietf.org/doc/html/rfc791#section-3.1

type IpHeader struct {
	VHL      uint8
	TOS      uint8
	TotalLen uint16
	ID       uint16
	Offset   uint16
	TTL      uint8
	Protocol uint8
	Checksum uint16
	Src      [V4AddrLen]byte
	Dst      [V4AddrLen]byte
	Options  [0]byte
}

func IpInputHandler(payload []byte, dev ethernet.IDevice) psErr.E {
	payloadLen := len(payload)

	if payloadLen < IpHeaderSizeMin {
		psLog.E(fmt.Sprintf("IP packet size is too small: %d bytes", payloadLen))
		return psErr.InvalidPacket
	}

	buf := bytes.NewBuffer(payload)
	hdr := IpHeader{}
	if err := binary.Read(buf, binary.BigEndian, &hdr); err != nil {
		psLog.E(fmt.Sprintf("binary.Read() failed: %s", err))
		return psErr.Error
	}

	if version := hdr.VHL >> 4; version != IpVersionV4 {
		psLog.E(fmt.Sprintf("IP version %d is not supported", version))
		return psErr.UnsupportedVersion
	}

	hdrLen := int(hdr.VHL & 0x0f)

	if payloadLen < hdrLen {
		psLog.E(fmt.Sprintf("IP packet length is too short: IHL = %d, Actual Packet Size = %d", hdrLen, payloadLen))
		return psErr.InvalidPacket
	}

	if totalLen := int(hdr.TotalLen); payloadLen < totalLen {
		psLog.E(fmt.Sprintf("IP packet length is too short: Total Length = %d, Actual Length = %d", totalLen, payloadLen))
		return psErr.InvalidPacket
	}

	if hdr.TTL == 0 {
		psLog.E("TTL expired")
		return psErr.TtlExpired
	}

	if sum := checksum(payload[:20]); sum != hdr.Checksum {
		psLog.E(fmt.Sprintf("Checksum mismatch: Expect = 0x%04x, Actual = 0x%04x", hdr.Checksum, sum))
		return psErr.ChecksumMismatch
	}

	iface := IfaceRepo.Get(dev, FamilyV4)
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

	return psErr.OK
}

// Computing the Internet Checksum
// https://datatracker.ietf.org/doc/html/rfc1071

func checksum(b []byte) uint16 {
	var sum uint32 = 0
	// sum up all fields of IP header by each 16bits (except Header Checksum and Options)
	for i := 0; i < len(b); i += 2 {
		// skip checksum field
		if i == 10 {
			continue
		}
		sum += uint32(uint16(b[i])<<8 | uint16(b[i+1]))
	}
	//
	sum = ((sum & 0xffff0000) >> 16) + (sum & 0x0000ffff)
	return ^(uint16(sum))
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
	psLog.I(fmt.Sprintf("\tprotocol:            %d (%s)", hdr.Protocol, protocolNumbers[int(hdr.Protocol)]))
	psLog.I(fmt.Sprintf("\theader checksum:     0x%04x", hdr.Checksum))
	psLog.I(fmt.Sprintf("\tsource address:      %d.%d.%d.%d", hdr.Src[0], hdr.Src[1], hdr.Src[2], hdr.Src[3]))
	psLog.I(fmt.Sprintf("\tdestination address: %d.%d.%d.%d", hdr.Dst[0], hdr.Dst[1], hdr.Dst[2], hdr.Dst[3]))
}

// ASSIGNED INTERNET PROTOCOL NUMBERS
// https://datatracker.ietf.org/doc/html/rfc790#page-6

var protocolNumbers = map[int]string {
	1: "ICMP",
	3: "Gateway-to-Gateway",
	4: "CMCC Gateway Monitoring Message",
	5: "ST",
	6: "TCP",
	7: "UCL",
	9: "Secure",
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
	63: "any local network",
	64: "SATNET and Backroom EXPAK",
	65: "MIT Subnet Support",
	69: "SATNET Monitoring",
	71: "Internet Packet Core Utility",
	76: "Backroom SATNET Monitoring",
	78: "WIDEBAND Monitoring",
	79: "WIDEBAND EXPAK",
}
