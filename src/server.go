// +build server

package main

import (
	"fmt"
	psErr "github.com/42milez/ProtocolStack/src/error"
	psLog "github.com/42milez/ProtocolStack/src/log"
	"github.com/42milez/ProtocolStack/src/mw"
	"github.com/42milez/ProtocolStack/src/net"
	"github.com/42milez/ProtocolStack/src/net/arp"
	"github.com/42milez/ProtocolStack/src/net/eth"
	"github.com/42milez/ProtocolStack/src/net/icmp"
	"github.com/42milez/ProtocolStack/src/net/ip"
	"github.com/42milez/ProtocolStack/src/repo"
	"github.com/42milez/ProtocolStack/src/worker"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

// arp: 3
// eth: 2
// icmp: 2
// ip: 2
const nServiceWorkers = 9

var wg sync.WaitGroup
var arpWg sync.WaitGroup
var ethWg sync.WaitGroup
var ipWg sync.WaitGroup
var icmpWg sync.WaitGroup
var repoWg sync.WaitGroup

var rxCh chan os.Signal
var sigCh chan os.Signal

func handleSignal(sigCh <-chan os.Signal, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		sig := <-sigCh
		psLog.I(fmt.Sprintf("Signal: %s", sig))
		rxCh <-syscall.SIGUSR1
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

	repo.RouteRepo.Register(mw.ParseIP(mw.LoopbackNetwork), mw.V4Any, iface1)

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

	repo.RouteRepo.Register(mw.ParseIP("192.0.0.0"), mw.V4Any, iface2)

	repo.RouteRepo.RegisterDefaultGateway(iface2, mw.ParseIP("192.0.2.1"))

	if err := start(); err != psErr.OK {
		return psErr.Error
	}

	return psErr.OK
}

func start() psErr.E {
	if err := arp.StartService(&arpWg); err != psErr.OK {
		return psErr.Error
	}
	if err := eth.StartService(&ethWg); err != psErr.OK {
		return psErr.Error
	}
	if err := ip.StartService(&ipWg); err != psErr.OK {
		return psErr.Error
	}
	if err := icmp.StartService(&icmpWg); err != psErr.OK {
		return psErr.Error
	}
	if err := repo.StartService(&repoWg); err != psErr.OK {
		return psErr.Error
	}

	var nWorkers int
	var zero = time.Now()
	var timeout = 10*time.Second
	for {
		select {
		case msg := <-arp.MonitorCh:
			if msg.Current == worker.Running {
				psLog.I(fmt.Sprintf("arp worker started: %d", int(msg.ID)))
				nWorkers += 1
			}
		case msg := <-eth.RcvTxCh:
			if msg.Current == worker.Running {
				psLog.I(fmt.Sprintf("eth worker started: %d", int(msg.ID)))
				nWorkers += 1
			}
		case msg := <- icmp.RcvTxCh:
			if msg.Current == worker.Running {
				psLog.I(fmt.Sprintf("icmp worker started: %d", int(msg.ID)))
				nWorkers += 1
			}
		case msg := <-ip.RcvTxCh:
			if msg.Current == worker.Running {
				psLog.I(fmt.Sprintf("ip worker started: %d", int(msg.ID)))
				nWorkers += 1
			}
		default:
			time.Sleep(100*time.Millisecond)
		}
		if nWorkers == nServiceWorkers {
			break
		}
		if time.Now().Sub(zero) > timeout {
			return psErr.Error
		}
	}

	psLog.I("//////////////////////////////////////////////////")

	return psErr.OK
}

func init() {
	rxCh = make(chan os.Signal)
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
			case <-rxCh:
				return
			default:
				time.Sleep(time.Second * 1)
			}
		}
	}()

	wg.Wait()

	psLog.I("Server stopped")
}
