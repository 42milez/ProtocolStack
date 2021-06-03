package main

import (
	"fmt"
	psErr "github.com/42milez/ProtocolStack/src/error"
	"github.com/42milez/ProtocolStack/src/ethernet"
	psLog "github.com/42milez/ProtocolStack/src/log"
	"github.com/42milez/ProtocolStack/src/network"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

var wg sync.WaitGroup

var ethSigCh chan os.Signal
var mainSigCh chan os.Signal
var netSigCh chan os.Signal
var sigCh chan os.Signal

var terminate bool

func handleSignal(sigCh <-chan os.Signal, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		sig := <-sigCh
		psLog.I(fmt.Sprintf("Signal: %s", sig))
		ethSigCh <- syscall.SIGUSR1
		mainSigCh <- syscall.SIGUSR1
		netSigCh <- syscall.SIGUSR1
	}()
}

func setup() psErr.E {
	psLog.I("--------------------------------------------------")
	psLog.I(" INITIALIZE DEVICES                               ")
	psLog.I("--------------------------------------------------")

	// Create a loopback device and its iface, then link them.
	loopbackDev := ethernet.GenLoopbackDevice("net" + strconv.Itoa(network.DeviceRepo.NextNumber()))
	if err := network.DeviceRepo.Register(loopbackDev); err != psErr.OK {
		return psErr.Error
	}

	iface1 := network.GenIface(network.LoopbackIpAddr, network.LoopbackNetmask, network.LoopbackBroadcast)
	if err := network.IfaceRepo.Register(iface1, loopbackDev); err != psErr.OK {
		return psErr.Error
	}

	network.RouteRepo.Register(network.ParseIP(network.LoopbackNetwork), network.V4Zero, iface1)

	// Create a TAP device and its interface, then link them.
	tapDev := ethernet.GenTapDevice(
		"net"+strconv.Itoa(network.DeviceRepo.NextNumber()),
		"tap0",
		ethernet.EthAddr{11, 22, 33, 44, 55, 66})
	if err := network.DeviceRepo.Register(tapDev); err != psErr.OK {
		return psErr.Error
	}

	iface2 := network.GenIface("192.0.2.2", "255.255.255.0", "192.0.2.255")
	if err := network.IfaceRepo.Register(iface2, tapDev); err != psErr.OK {
		return psErr.Error
	}

	network.RouteRepo.Register(network.ParseIP("192.0.0.0"), network.V4Zero, iface2)

	network.RouteRepo.RegisterDefaultGateway(iface2, network.ParseIP("192.0.2.1"))

	psLog.I("--------------------------------------------------")
	psLog.I(" START SERVER                                     ")
	psLog.I("--------------------------------------------------")

	if err := start(&wg); err != psErr.OK {
		return psErr.Error
	}

	return psErr.OK
}

func start(wg *sync.WaitGroup) psErr.E {
	if err := network.DeviceRepo.Up(); err != psErr.OK {
		return psErr.Error
	}

	// worker for watching I/O resource
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ethSigCh:
				psLog.I("Terminating Eth worker...")
				terminate = true
			default:
				if err := network.DeviceRepo.Poll(terminate); err != psErr.OK {
					// TODO: notify error to main goroutine
					// ...
					psLog.F(err.Error())
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
		for {
			select {
			case <-netSigCh:
				psLog.I("Terminating Net worker...")
				return
			case packet := <-ethernet.RxCh:
				if err := network.InputHandler(packet); err != psErr.OK {
					psLog.F(err.Error())
					// TODO: notify error to main goroutine
					// ...
				}
			case packet := <-ethernet.TxCh:
				if err := network.OutputHandler(packet); err != psErr.OK {
					psLog.F(err.Error())
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

func main() {
	if err := setup(); err != psErr.OK {
		psLog.F("Setup failed")
	}

	handleSignal(sigCh, &wg)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-mainSigCh:
				return
			default:
				time.Sleep(time.Second * 1)
			}
		}
	}()

	wg.Wait()

	psLog.I("Server stopped")
}
