package middleware

import (
	psErr "github.com/42milez/ProtocolStack/src/error"
	"github.com/42milez/ProtocolStack/src/ethernet"
	psLog "github.com/42milez/ProtocolStack/src/log"
	"github.com/42milez/ProtocolStack/src/network"
	s "github.com/42milez/ProtocolStack/src/syscall"
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

func register(protocolType ProtocolType, handler Handler) psErr.Error {
	for _, v := range protocols {
		if v.Type == protocolType {
			psLog.W("protocol is already registered")
			psLog.W("\ttype: %v ", protocolType.String())
			return psErr.Error{Code: psErr.CantRegister}
		}
	}

	p := Protocol{
		Type:    protocolType,
		Handler: handler,
	}

	protocols = append(protocols, p)

	psLog.I("protocol registered")
	psLog.I("\ttype: %v ", protocolType.String())

	return psErr.Error{Code: psErr.OK}
}

// TODO: delete this comment later
//int net_init(void)

func Setup() psErr.Error {
	// ARP
	// ...

	// ICMP
	// ...

	// IP
	// ...
	if err := register(ProtocolTypeIp, network.IpInputHandler); err.Code != psErr.OK {
		return psErr.Error{Code: psErr.Failed}
	}

	// TCP
	// ...

	// UDP
	// ...

	return psErr.Error{Code: psErr.OK}
}

func Start(netSigCh <-chan os.Signal, wg *sync.WaitGroup) psErr.Error {
	for _, dev := range devices {
		if err := dev.Open(); err.Code != psErr.OK {
			return psErr.Error{Code: psErr.Failed}
		}
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		var terminate = false
		for {
			select {
			case <-netSigCh:
				psLog.I("terminating worker...")
				terminate = true
			default:
				psLog.I("worker is running...")
			}
			for _, dev := range devices {
				if dev.FLAG&ethernet.DevFlagUp == 0 {
					continue
				}
				if err := dev.Op.Poll(dev, &s.Syscall{}, terminate); err.Code != psErr.OK {
					psLog.E("error: %v ", err)
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

	return psErr.Error{Code: psErr.OK}
}

func RegisterDevice(dev *ethernet.Device) {
	dev.Name = "net" + strconv.Itoa(len(devices))
	devices = append(devices, dev)
	psLog.I("device registered")
	psLog.I("\tname: %v (%v) ", dev.Name, dev.Priv.Name)
}

func RegisterInterface(iface *Iface, dev *ethernet.Device) psErr.Error {
	for _, v := range interfaces {
		if v.Dev == dev && v.Family == iface.Family {
			psLog.W("interface is already registered: %v ", v.Family.String())
			return psErr.Error{Code: psErr.CantRegister}
		}
	}

	interfaces = append(interfaces, iface)
	iface.Dev = dev

	psLog.I("interface attached")
	psLog.I("\tip:     %v ", iface.Unicast.String())
	psLog.I("\tdevice: %v (%v) ", iface.Dev.Name, iface.Dev.Priv.Name)

	return psErr.Error{Code: psErr.OK}
}
