package middleware

import (
	e "github.com/42milez/ProtocolStack/src/error"
	"github.com/42milez/ProtocolStack/src/ethernet"
	l "github.com/42milez/ProtocolStack/src/logger"
	"github.com/42milez/ProtocolStack/src/network"
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

type Handler func(data []uint8, dev ethernet.Device)

var devices []*ethernet.Device
var interfaces []*Iface

func init() {
	devices = make([]*ethernet.Device, 0)
	interfaces = make([]*Iface, 0)
}

func register(protocolType ProtocolType, handler Handler) e.Error {
	for _, v := range protocols {
		if v.Type == protocolType {
			l.W("protocol is already registered")
			l.W("\ttype: %v ", protocolType.String())
			return e.Error{Code: e.CantRegister}
		}
	}

	p := Protocol{
		Type:    protocolType,
		Handler: handler,
	}

	protocols = append(protocols, p)

	l.I("protocol registered")
	l.I("\ttype: %v ", protocolType.String())

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
				l.I("terminating worker...")
				terminate = true
			default:
				l.I("worker is running...")
			}
			for _, dev := range devices {
				if dev.FLAG&ethernet.DevFlagUp == 0 {
					continue
				}
				if err := dev.Op.Poll(dev, terminate); err.Code != e.OK {
					l.E("error: %v ", err)
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

func RegisterDevice(dev *ethernet.Device) {
	dev.Name = "net" + strconv.Itoa(len(devices))
	devices = append(devices, dev)
	l.I("device registered")
	l.I("\tname: %v (%v) ", dev.Name, dev.Priv.Name)
}

func RegisterInterface(iface *Iface, dev *ethernet.Device) e.Error {
	for _, v := range interfaces {
		if v.Dev == dev && v.Family == iface.Family {
			l.W("interface is already registered: %v ", v.Family.String())
			return e.Error{Code: e.CantRegister}
		}
	}

	interfaces = append(interfaces, iface)
	iface.Dev = dev

	l.I("interface attached")
	l.I("\tip:     %v ", iface.Unicast.String())
	l.I("\tdevice: %v (%v) ", iface.Dev.Name, iface.Dev.Priv.Name)

	return e.Error{Code: e.OK}
}
