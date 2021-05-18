package main

import (
	psErr "github.com/42milez/ProtocolStack/src/error"
	"github.com/42milez/ProtocolStack/src/ethernet"
	psLog "github.com/42milez/ProtocolStack/src/log"
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

func setup() psErr.Error {
	var dev *ethernet.Device
	var iface *middleware.Iface
	var err psErr.Error

	psLog.I("--------------------------------------------------")
	psLog.I(" INITIALIZE PROTOCOLS                             ")
	psLog.I("--------------------------------------------------")
	if err = middleware.Setup(); err.Code != psErr.OK {
		return psErr.Error{Code: psErr.Failed}
	}

	// Create a loopback device and its iface, then link them.
	psLog.I("--------------------------------------------------")
	psLog.I(" INITIALIZE DEVICES                               ")
	psLog.I("--------------------------------------------------")
	dev = ethernet.GenLoopbackDevice()
	middleware.RegisterDevice(dev)
	iface = middleware.GenIF(ethernet.LoopbackIpAddr, ethernet.LoopbackNetmask)
	if err = middleware.RegisterInterface(iface, dev); err.Code != psErr.OK {
		return psErr.Error{Code: psErr.Failed}
	}
	route.Register(iface, network.V4Zero)

	// Create a TAP device and its iface, then link them.
	if dev, err = ethernet.GenTapDevice("tap0", ethernet.EthAddr{11, 22, 33, 44, 55, 66}); err.Code != psErr.OK {
		return psErr.Error{Code: psErr.Failed}
	}
	middleware.RegisterDevice(dev)
	iface = middleware.GenIF("192.0.2.2", "255.255.255.0")
	if err = middleware.RegisterInterface(iface, dev); err.Code != psErr.OK {
		return psErr.Error{Code: psErr.Failed}
	}
	route.Register(iface, network.V4Zero)

	// Register the iface of the TAP device as the default gateway.
	route.RegisterDefaultGateway(iface, network.ParseIP("192.0.2.1"))

	// Create sub-thread for polling.
	psLog.I("--------------------------------------------------")
	psLog.I(" START WORKERS                                    ")
	psLog.I("--------------------------------------------------")
	if err = middleware.Start(netSigCh, &wg); err.Code != psErr.OK {
		return psErr.Error{Code: psErr.Failed}
	}

	return psErr.Error{Code: psErr.OK}
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
		psLog.I("signal received: %s ", sig)
		mainSigCh <- syscall.SIGUSR1
		netSigCh <- syscall.SIGUSR1
	}()
}

func main() {
	if err := setup(); err.Code != psErr.OK {
		psLog.F("setup failed.")
	}

	handleSignal(sigCh, &wg)

	psLog.I("                                                  ")
	psLog.I("//////////////////////////////////////////////////")
	psLog.I("           S E R V E R    S T A R T E D           ")
	psLog.I("//////////////////////////////////////////////////")
	psLog.I("                                                  ")

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-mainSigCh:
				psLog.I("shutting down server...")
				return
			default:
				psLog.I("server is running...")
				time.Sleep(time.Second * 3)
			}
		}
	}()

	wg.Wait()

	psLog.I("server stopped.")
}
