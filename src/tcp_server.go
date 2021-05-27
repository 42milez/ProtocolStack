package main

import (
	psErr "github.com/42milez/ProtocolStack/src/error"
	"github.com/42milez/ProtocolStack/src/ethernet"
	psLog "github.com/42milez/ProtocolStack/src/log"
	"github.com/42milez/ProtocolStack/src/network"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var wg sync.WaitGroup

var ethSigCh chan os.Signal
var mainSigCh chan os.Signal
var netSigCh chan os.Signal
var sigCh chan os.Signal

func setup() psErr.E {
	psLog.I("--------------------------------------------------")
	psLog.I(" INITIALIZE DEVICES                               ")
	psLog.I("--------------------------------------------------")

	// Create a loopback device and its iface, then link them.
	loopbackDev := network.GenLoopbackDevice()

	if err := network.DeviceRepo.Register(loopbackDev); err != psErr.OK {
		psLog.E("network.DeviceRepo.Register() failed: %s", err)
		return psErr.Error
	}

	iface1 := network.GenIface(ethernet.LoopbackIpAddr, ethernet.LoopbackNetmask, ethernet.LoopbackBroadcast)
	if err := network.IfaceRepo.Register(iface1, loopbackDev); err != psErr.OK {
		psLog.E("network.IfaceRepo.Register() failed: %s", err)
		return psErr.Error
	}

	network.RouteRepo.Register(network.ParseIP(ethernet.LoopbackNetwork), network.V4Zero, iface1)

	// Create a TAP device and its interface, then link them.
	tapDev := network.GenTapDevice(0, ethernet.EthAddr{11, 22, 33, 44, 55, 66})
	network.DeviceRepo.Register(tapDev)

	iface2 := network.GenIface("192.0.2.2", "255.255.255.0", "192.0.2.255")
	if err := network.IfaceRepo.Register(iface2, tapDev); err != psErr.OK {
		psLog.E("network.IfaceRepo.Register() failed: %s", err)
		return psErr.Error
	}

	network.RouteRepo.Register(network.ParseIP("192.0.0.0"), network.V4Zero, iface2)
	network.RouteRepo.RegisterDefaultGateway(iface2, network.ParseIP("192.0.2.1"))

	psLog.I("--------------------------------------------------")
	psLog.I(" START WORKERS                                    ")
	psLog.I("--------------------------------------------------")

	if err := start(&wg); err != psErr.OK {
		psLog.E("start() failed: %s", err)
		return psErr.Error
	}

	return psErr.OK
}

func start(wg *sync.WaitGroup) psErr.E {
	network.DeviceRepo.Up()

	// worker for polling incoming packets
	wg.Add(1)
	go func() {
		defer wg.Done()
		var terminate = false
		psLog.I("Eth worker started")
		for {
			select {
			case <-ethSigCh:
				psLog.I("Terminating Eth worker...")
				terminate = true
			default:
				if err := network.DeviceRepo.Poll(terminate); err != psErr.OK {
					// TODO: notify error to main goroutine
					// ...
					psLog.F("network.DeviceRepo.Poll() failed: %s", err)
				}
			}
			if terminate {
				return
			}
		}
	}()

	// worker for handling incoming/outgoing packets
	wg.Add(1)
	go func() {
		defer wg.Done()
		psLog.I("Net worker started")
		for {
			select {
			case <-netSigCh:
				psLog.I("Terminating Net worker...")
				return
			case packet := <-ethernet.RxCh:
				if err := network.InputHandler(packet); err != psErr.OK {
					psLog.F("network.InputHandler() failed: %s", err)
					// TODO: notify error to main goroutine
					// ...
				}
			case packet := <-ethernet.TxCh:
				if err := network.OutputHandler(packet); err != psErr.OK {
					psLog.F("network.OutputHandler() failed: %s", err)
					// TODO: notify error to main goroutine
					// ...
				}
			}
		}
	}()

	return psErr.OK
}

func init() {
	mainSigCh = make(chan os.Signal)
	ethSigCh = make(chan os.Signal)
	netSigCh = make(chan os.Signal)
	sigCh = make(chan os.Signal, 1)
	// https://pkg.go.dev/os/signal#Notify
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
}

func handleSignal(sigCh <-chan os.Signal, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		sig := <-sigCh
		psLog.I("Signal received: %s", sig)
		ethSigCh <- syscall.SIGUSR1
		mainSigCh <- syscall.SIGUSR1
		netSigCh <- syscall.SIGUSR1
	}()
}

func main() {
	if err := setup(); err != psErr.OK {
		psLog.F("Setup failed")
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
				psLog.I("Shutting down server...")
				return
			default:
				time.Sleep(time.Second * 1)
			}
		}
	}()

	wg.Wait()

	psLog.I("Server stopped")
}
