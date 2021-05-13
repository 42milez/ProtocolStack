package middleware

import (
	"fmt"
	"github.com/42milez/ProtocolStack/src/device"
	"github.com/42milez/ProtocolStack/src/e"
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
	Name string
	Interval syscall.Timeval
	Last syscall.Timeval
	Handler Handler
}

type Handler func(data []uint8, dev device.Device)

var devices []*device.Device
var interfaces []*Iface

func init() {
	devices = make([]*device.Device, 0)
	interfaces = make([]*Iface, 0)
}

func register(protocolType ProtocolType, handler Handler) error {
	for _, v := range protocols {
		if v.Type == protocolType {
			fmt.Printf("protocol is already registered: %v", protocolType.String())
			return e.CantRegister
		}
	}

	p := Protocol{
		Type: protocolType,
		Handler: handler,
	}

	protocols = append(protocols, p)

	log.Printf("registered a protocol: %v\n", protocolType.String())

	return nil
}

// TODO: delete this comment later
//int net_init(void)

func Setup() error {
	// ARP
	// ...

	// ICMP
	// ...

	// IP
	// ...
	if err := register(ProtocolTypeIp, network.IpInputHandler); err != nil {
		return err
	}

	// TCP
	// ...

	// UDP
	// ...

	return nil
}

func Start(netSigCh <-chan os.Signal, wg *sync.WaitGroup) error {
	for _, dev := range devices {
		if err := dev.Open(); err != e.OK {
			return e.Fatal
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
				if dev.FLAG & device.DevFlagUp == 0 {
					continue
				}
				if err := dev.Op.Poll(dev, terminate); err != nil {
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

	log.Println("ready for processing incoming data")

	return e.OK
}

func RegisterDevice(dev *device.Device) {
	dev.Name = "net" + strconv.Itoa(len(devices))
	devices = append(devices, dev)
	log.Printf("registered a device: %v\n", dev.Name)
}

func RegisterInterface(iface *Iface, dev *device.Device) error {
	for _, v := range interfaces {
		if v.Dev == dev && v.Family == iface.Family {
			fmt.Printf("interface is already registered: %v\n", v.Family.String())
			return e.CantRegister
		}
	}

	interfaces = append(interfaces, iface)
	iface.Dev = dev

	log.Println("attached an interface")
	log.Printf("\tIP Address:  %v", iface.Unicast.String())
	log.Printf("\tDevice Name: %v", iface.Dev.Name)

	return e.OK
}
