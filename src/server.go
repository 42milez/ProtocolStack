// +build server

package main

import (
	"fmt"
	psErr "github.com/42milez/ProtocolStack/src/error"
	"github.com/42milez/ProtocolStack/src/eth"
	psLog "github.com/42milez/ProtocolStack/src/log"
	"github.com/42milez/ProtocolStack/src/mw"
	"github.com/42milez/ProtocolStack/src/net"
	"github.com/42milez/ProtocolStack/src/net/arp"
	"github.com/42milez/ProtocolStack/src/net/icmp"
	"github.com/42milez/ProtocolStack/src/net/ip"
	"github.com/42milez/ProtocolStack/src/repo"
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
		ethSigCh <-syscall.SIGUSR1
		mainSigCh <-syscall.SIGUSR1
		netSigCh <-syscall.SIGUSR1
	}()
}

func setup() psErr.E {
	psLog.I("--------------------------------------------------")
	psLog.I(" INITIALIZE DEVICES                               ")
	psLog.I("--------------------------------------------------")

	// Create a loopback device and its iface, then link them.
	loopbackDev := eth.GenLoopbackDevice("net" + strconv.Itoa(repo.DeviceRepo.NextNumber()))
	if err := repo.DeviceRepo.Register(loopbackDev); err != psErr.OK {
		return psErr.Error
	}

	iface1 := net.GenIface(mw.LoopbackIpAddr, mw.LoopbackNetmask, mw.LoopbackBroadcast)
	if err := repo.IfaceRepo.Register(iface1, loopbackDev); err != psErr.OK {
		return psErr.Error
	}

	repo.RouteRepo.Register(mw.ParseIP(mw.LoopbackNetwork), mw.V4Zero, iface1)

	// Create a TAP device and its interface, then link them.
	tapDev := eth.GenTapDevice(
		"net"+strconv.Itoa(repo.DeviceRepo.NextNumber()),
		"tap0",
		mw.EthAddr{11, 22, 33, 44, 55, 66})
	if err := repo.DeviceRepo.Register(tapDev); err != psErr.OK {
		return psErr.Error
	}

	iface2 := net.GenIface("192.0.2.2", "255.255.255.0", "192.0.2.255")
	if err := repo.IfaceRepo.Register(iface2, tapDev); err != psErr.OK {
		return psErr.Error
	}

	repo.RouteRepo.Register(mw.ParseIP("192.0.0.0"), mw.V4Zero, iface2)

	repo.RouteRepo.RegisterDefaultGateway(iface2, mw.ParseIP("192.0.2.1"))

	psLog.I("--------------------------------------------------")
	psLog.I(" START SERVER                                     ")
	psLog.I("--------------------------------------------------")

	if err := start(&wg); err != psErr.OK {
		return psErr.Error
	}

	return psErr.OK
}

func start(wg *sync.WaitGroup) psErr.E {
	if err := repo.DeviceRepo.Up(); err != psErr.OK {
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
				if err := repo.DeviceRepo.Poll(terminate); err != psErr.OK {
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
			case msg := <-mw.EthRxCh:
				if err := mw.Rx(msg); err != psErr.OK {
					psLog.F(err.Error())
					// TODO: notify error to main goroutine
					// ...
				}
			case msg := <-mw.EthTxCh:
				if err := mw.Tx(msg); err != psErr.OK {
					psLog.F(err.Error())
					// TODO: notify error to main goroutine
					// ...
				}
			}
		}
	}()

	arp.StartService()
	ip.StartService()
	icmp.StartService()

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
