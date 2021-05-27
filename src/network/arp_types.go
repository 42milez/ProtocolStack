package network

// Address Resolution Protocol (ARP) Parameters
// https://www.iana.org/assignments/arp-parameters/arp-parameters.xhtml

// EtherType
// https://en.wikipedia.org/wiki/EtherType#Examples

// Notes:
//  - Protocol Type is same as EtherType.

const ArpHwTypeEthernet ArpHwType = 1

const ArpOpRequest ArpOpcode = 1
const ArpOpReply ArpOpcode = 2

const (
	ArpCacheStateFree ArpCacheState = iota
	ArpCacheStateIncomplete
	ArpCacheStateResolved
	ArpCacheStateStatic
)

type ArpHwType uint16

func (v ArpHwType) String() string {
	switch v {
	case ArpHwTypeEthernet:
		return "Ethernet"
	default:
		return "Unknown"
	}
}

type ArpOpcode uint16

func (v ArpOpcode) String() string {
	switch v {
	case ArpOpRequest:
		return "REQUEST"
	case ArpOpReply:
		return "REPLY"
	default:
		return "UNKNOWN"
	}
}

type ArpCacheState uint8
