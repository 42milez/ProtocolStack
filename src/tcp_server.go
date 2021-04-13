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
	var iface *network.Iface
	var err error

	if err = middleware.Setup(); err != nil {
		return err
	}

	// Create a loopback device and its interface, then link them.
	dev = ethernet.GenLoopbackDevice()
	device.Register(dev)
	iface = network.GenIF(ethernet.LoopbackIpAddr, ethernet.LoopbackNetmask)
	if err = network.AttachIF(iface, dev); err != nil {
		return err
	}
	route.Register(iface, network.V4Zero)

	// Create a TAP device and its interface, then link them.
	if dev, err = ethernet.GenTapDevice("tap0", "00:00:5e:00:53:01"); err != nil {
		return err
	}
	device.Register(dev)
	iface = network.GenIF("192.0.2.2", "255.255.255.0")
	if err = network.AttachIF(iface, dev); err != nil {
		return err
	}
	route.Register(iface, network.V4Zero)

	// Register the interface of the TAP device as the default gateway.
	route.RegisterDefaultGateway(iface)

	// Open TAP device.
	// ...

	// Create sub-thread for polling.
	// ...

	return nil
}

func main() {
	if err := setup(); err != nil {
		log.Println(err.Error())
		log.Fatal("Setup failed.")
	}

	log.Printf("Hello, TCP server!")
}
