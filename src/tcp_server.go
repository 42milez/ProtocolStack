package main

import (
	"fmt"
	"github.com/42milez/ProtocolStack/src/device"
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

var mainCh chan bool
var netCh chan bool
var sigCh chan os.Signal

func setup() error {
	var dev *device.Device
	var iface *middleware.Iface
	var err error

	if err = middleware.Setup(); err != nil {
		return err
	}

	// Create a loopback device and its iface, then link them.
	dev = ethernet.GenLoopbackDevice()
	middleware.RegisterDevice(dev)
	iface = middleware.GenIF(ethernet.LoopbackIpAddr, ethernet.LoopbackNetmask)
	if err = middleware.RegisterInterface(iface, dev); err != nil {
		return err
	}
	route.Register(iface, network.V4Zero)

	// Create a TAP device and its iface, then link them.
	if dev, err = ethernet.GenTapDevice("tap0", "00:00:5e:00:53:01"); err != nil {
		return err
	}
	middleware.RegisterDevice(dev)
	iface = middleware.GenIF("192.0.2.2", "255.255.255.0")
	if err = middleware.RegisterInterface(iface, dev); err != nil {
		return err
	}
	route.Register(iface, network.V4Zero)

	// Register the iface of the TAP device as the default gateway.
	route.RegisterDefaultGateway(iface, network.ParseIP("192.0.2.1"))

	// Create sub-thread for polling.
	if err = middleware.Start(netCh, &wg); err != nil {
		return err
	}

	return nil
}

func init() {
	mainCh = make(chan bool)
	netCh = make(chan bool)
	sigCh = make(chan os.Signal)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
}

func handleSignal(sigCh <-chan os.Signal, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		sig := <-sigCh
		fmt.Printf("signal received: %s\n", sig)
		mainCh <- true
		netCh <- true
	}()
}

func main() {
	if err := setup(); err != nil {
		log.Println(err.Error())
		log.Fatal("setup failed.")
	}
	log.Println("Hello, TCP server!")

	handleSignal(sigCh, &wg)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-mainCh:
				log.Println("main: shutting down...")
				return
			default:
				log.Println("main: running...")
				time.Sleep(time.Second)
			}
		}
	}()

	wg.Wait()

	log.Println("server stopped.")
}
