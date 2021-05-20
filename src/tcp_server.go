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
	//var iface *middleware.Iface
	var err psErr.Error

	// Create a loopback device and its iface, then link them.
	psLog.I("--------------------------------------------------")
	psLog.I(" INITIALIZE DEVICES                               ")
	psLog.I("--------------------------------------------------")
	loopbackDev := middleware.GenLoopbackDevice()

	if err = middleware.RegisterDevice(loopbackDev); err.Code != psErr.OK {
		return psErr.Error{Code: psErr.Failed}
	}

	iface1 := middleware.GenIface(ethernet.LoopbackIpAddr, ethernet.LoopbackNetmask)
	if err = middleware.RegisterInterface(iface1, loopbackDev); err.Code != psErr.OK {
		return psErr.Error{Code: psErr.Failed}
	}

	route.Register(iface1, network.V4Zero)

	// Create a TAP device and its iface, then link them.
	tapDev := middleware.GenTapDevice(0, ethernet.EthAddr{11, 22, 33, 44, 55, 66})
	if err.Code != psErr.OK {
		return psErr.Error{Code: psErr.Failed}
	}

	middleware.RegisterDevice(tapDev)

	iface2 := middleware.GenIface("192.0.2.2", "255.255.255.0")
	if err = middleware.RegisterInterface(iface2, tapDev); err.Code != psErr.OK {
		return psErr.Error{Code: psErr.Failed}
	}
	route.Register(iface2, network.V4Zero)

	// Register the iface of the TAP device as the default gateway.
	route.RegisterDefaultGateway(iface2, network.ParseIP("192.0.2.1"))

	// Create sub-thread for polling.
	psLog.I("--------------------------------------------------")
	psLog.I(" START WORKERS                                    ")
	psLog.I("--------------------------------------------------")
	if err = start(netSigCh, &wg); err.Code != psErr.OK {
		return psErr.Error{Code: psErr.Failed}
	}

	return psErr.Error{Code: psErr.OK}
}

func start(netSigCh <-chan os.Signal, wg *sync.WaitGroup) psErr.Error {
	middleware.Up()
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
			if err := middleware.Poll(terminate); err.Code != psErr.OK {
				// TODO: notify error to main goroutine
				// ...
				psLog.E("this is error message...")
			}
			if terminate {
				return
			}
		}
	}()

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
