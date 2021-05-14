package main

import (
	"github.com/42milez/ProtocolStack/src/device"
	e "github.com/42milez/ProtocolStack/src/error"
	"github.com/42milez/ProtocolStack/src/ethernet"
	"github.com/42milez/ProtocolStack/src/middleware"
	"github.com/42milez/ProtocolStack/src/network"
	"github.com/42milez/ProtocolStack/src/route"
	"log"
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
	var dev *device.Device
	var iface *middleware.Iface
	var err e.Error

	log.Println("--------------------------------------------------")
	log.Println(" PROTOCOL INITIALIZATION")
	log.Println("--------------------------------------------------")
	if err = middleware.Setup(); err.Code != e.OK {
		return e.Error{Code: e.Failed}
	}

	// Create a loopback device and its iface, then link them.
	log.Println("--------------------------------------------------")
	log.Println(" DEVICE INITIALIZATION")
	log.Println("--------------------------------------------------")
	dev = ethernet.GenLoopbackDevice()
	middleware.RegisterDevice(dev)
	iface = middleware.GenIF(ethernet.LoopbackIpAddr, ethernet.LoopbackNetmask)
	if err = middleware.RegisterInterface(iface, dev); err.Code != e.OK {
		return e.Error{Code: e.Failed}
	}
	route.Register(iface, network.V4Zero)

	// Create a TAP device and its iface, then link them.
	if dev, err = ethernet.GenTapDevice("tap0", "00:00:5e:00:53:01"); err.Code != e.OK {
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
	log.Println("--------------------------------------------------")
	log.Println(" START WORKERS")
	log.Println("--------------------------------------------------")
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
		log.Printf("signal received: %s\n", sig)
		mainSigCh <- syscall.SIGUSR1
		netSigCh <- syscall.SIGUSR1
	}()
}

func main() {
	if err := setup(); err.Code != e.OK {
		log.Fatal("setup failed.")
	}

	handleSignal(sigCh, &wg)

	log.Println("--------------------------------------------------")
	log.Println(" SERVER STARTED")
	log.Println("--------------------------------------------------")

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-mainSigCh:
				log.Println("shutting down server...")
				return
			default:
				log.Println("server is running...")
				time.Sleep(time.Second * 3)
			}
		}
	}()

	wg.Wait()

	log.Println("server stopped.")
}
