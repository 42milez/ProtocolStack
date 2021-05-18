package main

import (
	e "github.com/42milez/ProtocolStack/src/error"
	"github.com/42milez/ProtocolStack/src/ethernet"
	l "github.com/42milez/ProtocolStack/src/log"
	"github.com/42milez/ProtocolStack/src/middleware"
	"github.com/42milez/ProtocolStack/src/network"
	"github.com/42milez/ProtocolStack/src/route"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var wg sync.WaitGroup

var mainSigCh chan os.Signal
var netSigCh chan os.Signal
var sigCh chan os.Signal

func setup() e.Error {
	var dev *ethernet.Device
	var iface *middleware.Iface
	var err e.Error

	l.I("--------------------------------------------------")
	l.I(" INITIALIZE PROTOCOLS                             ")
	l.I("--------------------------------------------------")
	if err = middleware.Setup(); err.Code != e.OK {
		return e.Error{Code: e.Failed}
	}

	// Create a loopback device and its iface, then link them.
	l.I("--------------------------------------------------")
	l.I(" INITIALIZE DEVICES                               ")
	l.I("--------------------------------------------------")
	dev = ethernet.GenLoopbackDevice()
	middleware.RegisterDevice(dev)
	iface = middleware.GenIF(ethernet.LoopbackIpAddr, ethernet.LoopbackNetmask)
	if err = middleware.RegisterInterface(iface, dev); err.Code != e.OK {
		return e.Error{Code: e.Failed}
	}
	route.Register(iface, network.V4Zero)

	// Create a TAP device and its iface, then link them.
	if dev, err = ethernet.GenTapDevice("tap0", ethernet.EthAddr{11, 22, 33, 44, 55, 66}); err.Code != e.OK {
		return e.Error{Code: e.Failed}
	}
	middleware.RegisterDevice(dev)
	iface = middleware.GenIF("192.0.2.2", "255.255.255.0")
	if err = middleware.RegisterInterface(iface, dev); err.Code != e.OK {
		return e.Error{Code: e.Failed}
	}
	route.Register(iface, network.V4Zero)

	// Register the iface of the TAP device as the default gateway.
	route.RegisterDefaultGateway(iface, network.ParseIP("192.0.2.1"))

	// Create sub-thread for polling.
	l.I("--------------------------------------------------")
	l.I(" START WORKERS                                    ")
	l.I("--------------------------------------------------")
	if err = middleware.Start(netSigCh, &wg); err.Code != e.OK {
		return e.Error{Code: e.Failed}
	}

	return e.Error{Code: e.OK}
}

func init() {
	mainSigCh = make(chan os.Signal)
	netSigCh = make(chan os.Signal)
	sigCh = make(chan os.Signal)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
}

func handleSignal(sigCh <-chan os.Signal, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		sig := <-sigCh
		l.I("signal received: %s ", sig)
		mainSigCh <- syscall.SIGUSR1
		netSigCh <- syscall.SIGUSR1
	}()
}

func main() {
	if err := setup(); err.Code != e.OK {
		l.F("setup failed.")
	}

	handleSignal(sigCh, &wg)

	l.I("                                                  ")
	l.I("//////////////////////////////////////////////////")
	l.I("           S E R V E R    S T A R T E D           ")
	l.I("//////////////////////////////////////////////////")
	l.I("                                                  ")

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-mainSigCh:
				l.I("shutting down server...")
				return
			default:
				l.I("server is running...")
				time.Sleep(time.Second * 3)
			}
		}
	}()

	wg.Wait()

	l.I("server stopped.")
}
