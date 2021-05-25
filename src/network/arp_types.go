package network

// https://www.iana.org/assignments/arp-parameters/arp-parameters.xhtml

const ArpHwTypeEthernet = 1

type ArpHwType uint16

func (v ArpHwType) String() string {
	switch v {
	case ArpHwTypeEthernet:
		return "Ethernet"
	default:
		return "UNKNOWN"
	}
}

// https://www.iana.org/assignments/arp-parameters/arp-parameters.xhtml#arp-parameters-1

const ArpOpRequest = 1
const ArpOpReply = 2

type ArpOp uint16

func (v ArpOp) String() string {
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

const (
	ArpCacheStateFree ArpCacheState = iota
	ArpCacheStateIncomplete
	ArpCacheStateResolved
	ArpCacheStateStatic
)
