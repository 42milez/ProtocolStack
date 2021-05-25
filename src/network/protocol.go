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
		return "PROTOCOL_TYPE_ARP"
	case ProtocolTypeIcmp:
		return "PROTOCOL_TYPE_ICMP"
	case ProtocolTypeIp:
		return "PROTOCOL_TYPE_IP"
	case ProtocolTypeTcp:
		return "PROTOCOL_TYPE_TCP"
	case ProtocolTypeUdp:
		return "PROTOCOL_TYPE_UDP"
	default:
		return "UNKNOWN"
	}
}
