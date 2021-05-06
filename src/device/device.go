package device

import (
	"errors"
	"fmt"
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

const (
	DevFlagUp        DevFlag = 0x0001
	DevFlagLoopback          = 0x0010
	DevFlagBroadcast         = 0x0020
	DevFlagP2P               = 0x0040
	DevFlagNeedArp           = 0x0100
)

type Operation struct {
	Open     func(dev *Device) error
	Close    func(dev *Device) error
	Transmit func(dev *Device) error
	Poll     func(dev *Device) error
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

func (dev *Device) Open() error {
	if dev.Op.Open != nil {
		if (dev.FLAG & DevFlagUp) != 0 {
			return errors.New(fmt.Sprintf("%s already opend", dev.Name))
		}
		if err := dev.Op.Open(dev); err != nil {
			return err
		}
		dev.FLAG |= DevFlagUp
		log.Printf("%s opened", dev.Name)
	}
	return nil
}
