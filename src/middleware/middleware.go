package middleware

import (
	"errors"
	"fmt"
	"github.com/42milez/ProtocolStack/src/device"
	"github.com/42milez/ProtocolStack/src/network"
	"log"
	"strconv"
	"sync"
	"syscall"
	"time"
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
			return errors.New(fmt.Sprintf("%s is already registered", protocolType.String()))
		}
	}

	p := Protocol{
		Type: protocolType,
		Handler: handler,
	}

	protocols = append(protocols, p)

	log.Printf("%s is registered.\n", protocolType.String())

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

func Start(wg *sync.WaitGroup) error {
	var err error

	for _, dev := range devices {
		if err = dev.Open(); err != nil {
			return err
		}
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Println("running...")
		time.Sleep(time.Second)
	}()

	log.Println("started.")

	return nil
}

func RegisterDevice(dev *device.Device) {
	dev.Name = "net" + strconv.Itoa(len(devices))
	devices = append(devices, dev)
	log.Printf("device registered: dev=%s\n", dev.Name)
}

func RegisterInterface(iface *Iface, dev *device.Device) error {
	for _, v := range interfaces {
		if v.Dev == dev && v.Family == iface.Family {
			return errors.New(fmt.Sprintf("%s is already exists", v.Family.String()))
		}
	}
	interfaces = append(interfaces, iface)
	iface.Dev = dev
	log.Printf("iface attached: iface=%s, dev=%s", iface.Unicast.String(), iface.Dev.Name)
	return nil
}
