//go:generate mockgen -source=repository.go -destination=repository_mock.go -package=$GOPACKAGE -self_package=github.com/42milez/ProtocolStack/src/$GOPACKAGE

package repo

import (
	"fmt"
	psErr "github.com/42milez/ProtocolStack/src/error"
	psLog "github.com/42milez/ProtocolStack/src/log"
	"github.com/42milez/ProtocolStack/src/monitor"
	"github.com/42milez/ProtocolStack/src/mw"
	"github.com/42milez/ProtocolStack/src/worker"
	"sync"
)

const xChBufSize = 5

var monCh chan *worker.Message
var sigCh chan *worker.Message
var watcherId uint32

var DeviceRepo IDeviceRepo
var IfaceRepo IIfaceRepo
var RouteRepo IRouteRepo

type Handler func(data []byte, dev mw.IDevice) psErr.E

type Route struct {
	Network mw.IP
	Netmask mw.IP
	NextHop mw.IP
	Iface   *mw.Iface
}

type IDeviceRepo interface {
	NextNumber() int
	Poll() psErr.E
	Register(dev mw.IDevice) psErr.E
	Up() psErr.E
}

type deviceRepo struct {
	devices []mw.IDevice
}

func (p *deviceRepo) NextNumber() int {
	return len(p.devices)
}

func (p *deviceRepo) Poll() psErr.E {
	for _, dev := range p.devices {
		if !dev.IsUp() {
			continue
		}
		if err := dev.Poll(); err != psErr.OK {
			if err == psErr.Interrupted {
				return psErr.OK
			}
			return psErr.Error
		}
	}
	return psErr.OK
}

func (p *deviceRepo) Register(dev mw.IDevice) psErr.E {
	for _, d := range p.devices {
		if d.Equal(dev) {
			psLog.W("Device is already registered")
			psLog.W(fmt.Sprintf("\ttype: %s", d.Type()))
			psLog.W(fmt.Sprintf("\tname: %s (%s)", d.Name(), d.Priv().Name))
			return psErr.Error
		}
	}
	p.devices = append(p.devices, dev)
	psLog.I("Device was registered")
	psLog.I(fmt.Sprintf("\ttype:      %s", dev.Type()))
	psLog.I(fmt.Sprintf("\tname:      %s (%s)", dev.Name(), dev.Priv().Name))
	psLog.I(fmt.Sprintf("\taddr:      %s", dev.Addr()))
	return psErr.OK
}

func (p *deviceRepo) Up() psErr.E {
	for _, dev := range p.devices {
		if dev.IsUp() {
			psLog.W("Device is already up")
			psLog.W(fmt.Sprintf("\ttype: %s", dev.Type()))
			psLog.W(fmt.Sprintf("\tname: %s (%s)", dev.Name(), dev.Priv().Name))
			return psErr.Error
		}
		if err := dev.Open(); err != psErr.OK {
			psLog.E(fmt.Sprintf("Can't open device: %s", err))
			psLog.E(fmt.Sprintf("\ttype: %s", dev.Type()))
			psLog.E(fmt.Sprintf("\tname: %s (%s)", dev.Name(), dev.Priv().Name))
			return psErr.Error
		}
		dev.Up()
		psLog.I("Device was opened")
		psLog.I(fmt.Sprintf("\ttype: %s", dev.Type()))
		psLog.I(fmt.Sprintf("\tname: %s (%s)", dev.Name(), dev.Priv().Name))
	}
	return psErr.OK
}

type IIfaceRepo interface {
	Get(unicast mw.IP) *mw.Iface
	Lookup(dev mw.IDevice, family mw.AddrFamily) *mw.Iface
	Register(iface *mw.Iface, dev mw.IDevice) psErr.E
}

type ifaceRepo struct {
	ifaces []*mw.Iface
}

func (p *ifaceRepo) Get(unicast mw.IP) *mw.Iface {
	for _, v := range p.ifaces {
		if v.Unicast.Equal(unicast) {
			return v
		}
	}
	return nil
}

func (p *ifaceRepo) Lookup(dev mw.IDevice, family mw.AddrFamily) *mw.Iface {
	for _, v := range p.ifaces {
		if v.Dev.Equal(dev) && v.Family == family {
			return v
		}
	}
	return nil
}

func (p *ifaceRepo) Register(iface *mw.Iface, dev mw.IDevice) psErr.E {
	for _, i := range p.ifaces {
		if i.Dev.Equal(dev) && i.Family == iface.Family {
			psLog.W(fmt.Sprintf("Interface is already registered: %s", i.Family))
			return psErr.Error
		}
	}

	p.ifaces = append(p.ifaces, iface)
	iface.Dev = dev

	psLog.I("Interface was attached")
	psLog.I(fmt.Sprintf("\tip:     %s", iface.Unicast))
	psLog.I(fmt.Sprintf("\tdevice: %s (%s)", dev.Name(), dev.Priv().Name))

	return psErr.OK
}

type IRouteRepo interface {
	Get(ip mw.IP) *Route
	Register(network mw.IP, nextHop mw.IP, iface *mw.Iface)
	RegisterDefaultGateway(iface *mw.Iface, nextHop mw.IP)
}

type routeRepo struct {
	routes []*Route
}

func (p *routeRepo) Get(ip mw.IP) *Route {
	var ret *Route
	for _, route := range p.routes {
		if ip.Mask(route.Netmask).Equal(route.Network) {
			// Longest prefix match
			// https://en.wikipedia.org/wiki/Longest_prefix_match
			if ret == nil || mw.LongestIP(ret.Netmask, route.Netmask).Equal(route.Netmask) {
				ret = route
			}
		}
	}
	return ret
}

func (p *routeRepo) Register(network mw.IP, nextHop mw.IP, iface *mw.Iface) {
	route := &Route{
		Network: network,
		Netmask: iface.Netmask,
		NextHop: nextHop,
		Iface:   iface,
	}
	p.routes = append(p.routes, route)
	psLog.I("Route was registered")
	psLog.I(fmt.Sprintf("\tnetwork:  %s", route.Network))
	psLog.I(fmt.Sprintf("\tnetmask:  %s", route.Netmask))
	psLog.I(fmt.Sprintf("\tunicast:  %s", iface.Unicast))
	psLog.I(fmt.Sprintf("\tnext hop: %s", nextHop))
	psLog.I(fmt.Sprintf("\tdevice:   %s (%s)", iface.Dev.Name(), iface.Dev.Priv().Name))
}

func (p *routeRepo) RegisterDefaultGateway(iface *mw.Iface, nextHop mw.IP) {
	route := &Route{
		Network: mw.V4Any,
		Netmask: mw.V4Any,
		NextHop: nextHop,
		Iface:   iface,
	}
	p.routes = append(p.routes, route)
	psLog.I("Default gateway was registered")
	psLog.I(fmt.Sprintf("\tnetwork:  %s", route.Network))
	psLog.I(fmt.Sprintf("\tnetmask:  %s", route.Netmask))
	psLog.I(fmt.Sprintf("\tunicast:  %s", iface.Unicast))
	psLog.I(fmt.Sprintf("\tnext hop: %s", nextHop))
	psLog.I(fmt.Sprintf("\tdevice:   %s (%s)", iface.Dev.Name(), iface.Dev.Priv().Name))
}

func Start(wg *sync.WaitGroup) psErr.E {
	if err := DeviceRepo.Up(); err != psErr.OK {
		return psErr.Error
	}

	wg.Add(1)
	go watcher(wg)

	return psErr.OK
}

func Stop() {
	msg := &worker.Message{
		Desired: worker.Stopped,
	}
	sigCh <- msg
}

func watcher(wg *sync.WaitGroup) {
	defer wg.Done()

	monCh <- &worker.Message{
		ID:      watcherId + 999,
		Current: worker.Running,
	}

	for {
		select {
		case msg := <-sigCh:
			if msg.Desired == worker.Stopped {
				return
			}
		default:
			if err := DeviceRepo.Poll(); err != psErr.OK {
				monCh <- &worker.Message{
					Current: worker.Error,
				}
				return
			}
		}
	}
}

func init() {
	monCh = make(chan *worker.Message, xChBufSize)
	sigCh = make(chan *worker.Message, xChBufSize)
	watcherId = monitor.Register("Repository Watcher", monCh, sigCh)

	DeviceRepo = &deviceRepo{}
	IfaceRepo = &ifaceRepo{}
	RouteRepo = &routeRepo{}
}
