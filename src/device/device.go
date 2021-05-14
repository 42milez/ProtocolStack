package device

import (
	e "github.com/42milez/ProtocolStack/src/error"
	"log"
)

type DevType int

const (
	DevTypeNull DevType = iota
	DevTypeLoopback
	DevTypeEthernet
)

func (t DevType) String() string {
	switch t {
	case DevTypeNull:
		return "DEVICE_TYPE_NULL"
	case DevTypeLoopback:
		return "DEVICE_TYPE_LOOPBACK"
	case DevTypeEthernet:
		return "DEVICE_TYPE_ETHERNET"
	default:
		return "UNKNOWN"
	}
}

type DevFlag uint16

const DevFlagUp DevFlag = 0x0001
const DevFlagLoopback DevFlag = 0x0010
const DevFlagBroadcast DevFlag = 0x0020
const DevFlagP2P DevFlag = 0x0040
const DevFlagNeedArp DevFlag = 0x0100

type Operation struct {
	Open     func(dev *Device) e.Error
	Close    func(dev *Device) e.Error
	Transmit func(dev *Device) e.Error
	Poll     func(dev *Device, terminate bool) e.Error
}

type Privilege struct {
	Name [16]byte
	FD   int
}

// TODO: delete this comment later
//struct net_device {
//	struct net_device *next;
//	struct net_iface *ifaces; /* NOTE: if you want to add/delete the entries after net_run(), you need to protect ifaces with a mutex. */
//	unsigned int index;
//	char name[IFNAMSIZ];
//	uint16_t type;
//	uint16_t mtu;
//	uint16_t flags;
//	uint16_t hlen; /* header length */
//	uint16_t alen; /* address length */
//	uint8_t addr[NET_DEVICE_ADDR_LEN];
//	union {
//	uint8_t peer[NET_DEVICE_ADDR_LEN];
//	uint8_t broadcast[NET_DEVICE_ADDR_LEN];
//	};
//	struct net_device_ops *ops;
//	void *priv;
//};

type Device struct {
	Name      string
	Type      DevType
	MTU       uint16
	FLAG      DevFlag
	HeaderLen uint16
	AddrLen   uint16
	Addr      []byte
	Peer      []byte
	Broadcast []byte
	Op        Operation
	Priv      Privilege
}

func (dev *Device) Open() e.Error {
	if dev.Op.Open != nil {
		if (dev.FLAG & DevFlagUp) != 0 {
			log.Println("device is already opened")
			log.Printf("\tname: %v (%v)\n", dev.Name, dev.Priv.Name)
			return e.Error{Code: e.AlreadyOpened}
		}
		if err := dev.Op.Open(dev); err.Code != e.OK {
			log.Printf("can't open a device")
			log.Printf("\tname: %v (%v)\n", dev.Name, dev.Priv.Name)
			log.Printf("\ttype: %v\n", dev.Type)
			return e.Error{Code: e.CantOpen}
		}
		dev.FLAG |= DevFlagUp
		log.Println("successfully opened a device")
		log.Printf("\tname: %v (%v)\n", dev.Name, dev.Priv.Name)
	}
	return e.Error{Code: e.OK}
}
