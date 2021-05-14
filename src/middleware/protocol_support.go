package middleware

import "sync"

type ProtocolType uint16

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

// TODO: delete this comment later
//struct net_protocol {
//    struct net_protocol *next;
//    char name[16];
//    uint16_t type;
//    pthread_mutex_t mutex;   // mutex for input queue
//    struct queue_head queue; // input queue
//    void (*handler)(const uint8_t *data, size_t len, struct net_device *dev);
//};

type Protocol struct {
	Type    ProtocolType
	Mutex   *sync.Mutex
	Handler Handler
}
