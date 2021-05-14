package middleware

import (
	"github.com/42milez/ProtocolStack/src/device"
	e "github.com/42milez/ProtocolStack/src/error"
	"github.com/42milez/ProtocolStack/src/network"
	"log"
	"os"
	"strconv"
	"sync"
	"syscall"
)

var protocols []Protocol

// TODO: delete this comment later
//struct net_timer {
//	struct net_timer *next;
//	char name[16];
//	struct timeval interval;
//	struct timeval last;
//	void (*handler)(void);
//};

type Timer struct {
	Name     string
	Interval syscall.Timeval
	Last     syscall.Timeval
	Handler  Handler
}

type Handler func(data []uint8, dev device.Device)

var devices []*device.Device
var interfaces []*Iface

func init() {
	devices = make([]*device.Device, 0)
	interfaces = make([]*Iface, 0)
}

func register(protocolType ProtocolType, handler Handler) e.Error {
	for _, v := range protocols {
		if v.Type == protocolType {
			log.Println("protocol is already registered")
			log.Printf("\ttype: %v\n", protocolType.String())
			return e.Error{Code: e.CantRegister}
		}
	}

	p := Protocol{
		Type:    protocolType,
		Handler: handler,
	}

	protocols = append(protocols, p)

	log.Println("registered a protocol")
	log.Printf("\ttype: %v\n", protocolType.String())

	return e.Error{Code: e.OK}
}

// TODO: delete this comment later
//int net_init(void)

func Setup() e.Error {
	// ARP
	// ...

	// ICMP
	// ...

	// IP
	// ...
	if err := register(ProtocolTypeIp, network.IpInputHandler); err.Code != e.OK {
		return e.Error{Code: e.Failed}
	}

	// TCP
	// ...

	// UDP
	// ...

	return e.Error{Code: e.OK}
}

func Start(netSigCh <-chan os.Signal, wg *sync.WaitGroup) e.Error {
	for _, dev := range devices {
		if err := dev.Open(); err.Code != e.OK {
			return e.Error{Code: e.Failed}
		}
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		var terminate = false
		for {
			select {
			case <-netSigCh:
				log.Println("terminating receiver...")
				terminate = true
			default:
				log.Println("receiver is running...")
			}
			for _, dev := range devices {
				if dev.FLAG&device.DevFlagUp == 0 {
					continue
				}
				if err := dev.Op.Poll(dev, terminate); err.Code != e.OK {
					log.Printf("error: %v\n", err)
					// TODO: notify error to main goroutine
					// ...
					return
				}
			}
			if terminate {
				return
			}
		}
	}()

	return e.Error{Code: e.OK}
}

func RegisterDevice(dev *device.Device) {
	dev.Name = "net" + strconv.Itoa(len(devices))
	devices = append(devices, dev)
	log.Println("registered a device")
	log.Printf("\tname: %v (%v)\n", dev.Name, dev.Priv.Name)
}

func RegisterInterface(iface *Iface, dev *device.Device) e.Error {
	for _, v := range interfaces {
		if v.Dev == dev && v.Family == iface.Family {
			log.Printf("interface is already registered: %v\n", v.Family.String())
			return e.Error{Code: e.CantRegister}
		}
	}

	interfaces = append(interfaces, iface)
	iface.Dev = dev

	log.Println("attached an interface")
	log.Printf("\tip:     %v\n", iface.Unicast.String())
	log.Printf("\tdevice: %v (%v)\n", iface.Dev.Name, iface.Dev.Priv.Name)

	return e.Error{Code: e.OK}
}
