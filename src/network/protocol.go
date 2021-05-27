package network

type ProtocolType uint8

const (
	ProtocolTypeArp ProtocolType = iota
	ProtocolTypeIcmp
	ProtocolTypeIp
	ProtocolTypeTcp
	ProtocolTypeUdp
)

func (t ProtocolType) String() string {
	switch t {
	case ProtocolTypeArp:
		return "ARP"
	case ProtocolTypeIcmp:
		return "ICMP"
	case ProtocolTypeIp:
		return "IP"
	case ProtocolTypeTcp:
		return "TCP"
	case ProtocolTypeUdp:
		return "UDP"
	default:
		return "UNKNOWN"
	}
}
