package net

import (
	"fmt"
	"github.com/42milez/ProtocolStack/src/eth"
	"time"
)

// Address Resolution Protocol (ARP) Parameters
// https://www.iana.org/assignments/arp-parameters/arp-parameters.xhtml

// EtherType
// https://en.wikipedia.org/wiki/EtherType#Examples

// Notes:
//  - Protocol Type is same as EtherType.

const ArpCacheSize = 32 // number of cache entries
const (
	ArpCacheStateFree ArpCacheStatus = iota
	ArpCacheStateIncomplete
	ArpCacheStateResolved
	ArpCacheStateStatic
)
const ArpHwTypeEthernet ArpHwType = 0x0001
const ArpOpRequest ArpOpcode = 0x0001
const ArpOpReply ArpOpcode = 0x0002
const ArpPacketLen = 28 // byte
const (
	ArpStatusComplete ArpStatus = iota
	ArpStatusIncomplete
	ArpStatusError
)
const arpCacheLifetime = 24 * time.Hour

// Hardware Types
// https://www.iana.org/assignments/arp-parameters/arp-parameters.xhtml#arp-parameters-2

// Operation Codes
// https://www.iana.org/assignments/arp-parameters/arp-parameters.xhtml#arp-parameters-1

var arpHwTypes = map[ArpHwType]string{
	// 0: Reserved
	1:  "Ethernet (10Mb)",
	2:  "Experimental Ethernet (3Mb)",
	3:  "Amateur Radio AX.25",
	4:  "Proteon ProNET Token Ring",
	5:  "Chaos",
	6:  "IEEE 802 Networks",
	7:  "ARCNET",
	8:  "Hyperchannel",
	9:  "Lanstar",
	10: "Autonet Short Address",
	11: "LocalTalk",
	12: "LocalNet (IBM PCNet or SYTEK LocalNET)",
	13: "Ultra link",
	14: "SMDS",
	15: "Frame Relay",
	16: "Asynchronous Transmission Mode (ATM)",
	17: "HDLC",
	18: "Fibre Channel",
	19: "Asynchronous Transmission Mode (ATM)",
	20: "Serial Line",
	21: "Asynchronous Transmission Mode (ATM)",
	22: "MIL-STD-188-220",
	23: "Metricom",
	24: "IEEE 1394.1995",
	25: "MAPOS",
	26: "Twinaxial",
	27: "EUI-64",
	28: "HIPARP",
	29: "IP and ARP over ISO 7816-3",
	30: "ARPSec",
	31: "IPsec tunnel",
	32: "InfiniBand (TM)",
	33: "TIA-102 Project 25 Common Air Interface (CAI)",
	34: "Wiegand Interface",
	35: "Pure IP",
	36: "HW_EXP1",
	37: "HFI",
	// 38-255: Unassigned
	256: "HW_EXP2",
	257: "AEthernet",
	// 258-65534: Unassigned
	// 65535: Reserved
}
var arpOpCodes = map[ArpOpcode]string{
	// 0: Reserved
	1:  "REQUEST",
	2:  "REPLY",
	3:  "request Reverse",
	4:  "reply Reverse",
	5:  "DRARP-Request",
	6:  "DRARP-Reply",
	7:  "DRARP-Error",
	8:  "InARP-Request",
	9:  "InARP-Reply",
	10: "ARP-NAK",
	11: "MARS-Request",
	12: "MARS-Multi",
	13: "MARS-MServ",
	14: "MARS-Join",
	15: "MARS-Leave",
	16: "MARS-NAK",
	17: "MARS-Unserv",
	18: "MARS-SJoin",
	19: "MARS-SLeave",
	20: "MARS-Grouplist-Request",
	21: "MARS-Grouplist-Reply",
	22: "MARS-Redirect-Map",
	23: "MAPOS-UNARP",
	24: "OP_EXP1",
	25: "OP_EXP2",
	// 26-65534: Unassigned
	// 65535: Reserved
}

type ArpCacheStatus uint8

type ArpHdr struct {
	HT     ArpHwType // hardware type
	PT     eth.Type  // protocol type
	HAL    uint8     // hardware address length
	PAL    uint8     // protocol address length
	Opcode ArpOpcode
}

type ArpHwType uint16

func (v ArpHwType) String() string {
	return arpHwTypes[v]
}

type ArpOpcode uint16

func (v ArpOpcode) String() string {
	return arpOpCodes[v]
}

type ArpPacket struct {
	ArpHdr
	SHA eth.Addr     // sender hardware address
	SPA ArpProtoAddr // sender protocol address
	THA eth.Addr     // target hardware address
	TPA ArpProtoAddr // target protocol address
}

type ArpProtoAddr [V4AddrLen]byte

func (p ArpProtoAddr) String() string {
	return fmt.Sprintf("%d.%d.%d.%d", p[0], p[1], p[2], p[3])
}

type ArpStatus int
