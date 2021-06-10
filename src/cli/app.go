package cli

import (
	psErr "github.com/42milez/ProtocolStack/src/error"
	psLog "github.com/42milez/ProtocolStack/src/log"
	"github.com/42milez/ProtocolStack/src/monitor"
	"github.com/42milez/ProtocolStack/src/mw"
	"github.com/42milez/ProtocolStack/src/net"
	"github.com/42milez/ProtocolStack/src/net/arp"
	"github.com/42milez/ProtocolStack/src/net/eth"
	"github.com/42milez/ProtocolStack/src/net/icmp"
	"github.com/42milez/ProtocolStack/src/net/ip"
	"github.com/42milez/ProtocolStack/src/net/tcp"
	"github.com/42milez/ProtocolStack/src/repo"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

const serviceTimeout = 3 * time.Second

var sigCh chan os.Signal

var arpWg sync.WaitGroup
var ethWg sync.WaitGroup
var icmpWg sync.WaitGroup
var ipWg sync.WaitGroup
var monitorWg sync.WaitGroup
var repoWg sync.WaitGroup
var tcpWg sync.WaitGroup

func setup() psErr.E {
	psLog.D(
		"-------------------------------------------------------",
		"          I N I T I A L I Z E   D E V I C E S          ",
		"-------------------------------------------------------")

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

	psLog.D(
		"-------------------------------------------------------",
		"              S T A R T   S E R V I C E S              ",
		"-------------------------------------------------------")

	if err := startServices(); err != psErr.OK {
		return psErr.Error
	}

	psLog.D(
		"///////////////////////////////////////////////////////",
		"         A P P L I C A T I O N   S T A R T E D         ",
		"///////////////////////////////////////////////////////")

	return psErr.OK
}

func startServices() psErr.E {
	if err := arp.Start(&arpWg); err != psErr.OK {
		return psErr.Error
	}
	if err := eth.Start(&ethWg); err != psErr.OK {
		return psErr.Error
	}
	if err := icmp.Start(&icmpWg); err != psErr.OK {
		return psErr.Error
	}
	if err := ip.Start(&ipWg); err != psErr.OK {
		return psErr.Error
	}
	if err := monitor.Start(&monitorWg); err != psErr.OK {
		return psErr.Error
	}
	if err := repo.Start(&repoWg); err != psErr.OK {
		return psErr.Error
	}
	if err := tcp.Start(&tcpWg); err != psErr.OK {
		return psErr.Error
	}

	var zero = time.Now()
	for {
		if monitor.Status() == monitor.Green {
			break
		}
		if time.Since(zero) > serviceTimeout {
			psLog.E("some services didn't ready within the time")
			return psErr.Error
		}
	}

	psLog.D("all service workers started")

	return psErr.OK
}

func stopServices() {
	arp.Stop()
	eth.Stop()
	icmp.Stop()
	ip.Stop()
	monitor.Stop()
	repo.Stop()
	tcp.Stop()

	arpWg.Wait()
	ethWg.Wait()
	icmpWg.Wait()
	ipWg.Wait()
	monitorWg.Wait()
	repoWg.Wait()
	tcpWg.Wait()
}

func init() {
	sigCh = make(chan os.Signal, 1)
	// https://pkg.go.dev/os/signal#Notify
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
}
