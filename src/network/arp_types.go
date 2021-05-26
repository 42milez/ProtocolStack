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

type ArpOpcode uint16

type ArpCacheState uint8
