package main

import (
	"github.com/42milez/ProtocolStack/src/device"
	"github.com/42milez/ProtocolStack/src/ethernet"
	"github.com/42milez/ProtocolStack/src/middleware"
	"github.com/42milez/ProtocolStack/src/network"
	"github.com/42milez/ProtocolStack/src/route"
	"log"
)

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
	if err = middleware.Start(); err != nil {
		return err
	}

	return nil
}

func main() {
	if err := setup(); err != nil {
		log.Println(err.Error())
		log.Fatal("Setup failed.")
	}

	log.Printf("Hello, TCP server!")
}
