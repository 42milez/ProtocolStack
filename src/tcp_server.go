package main

import (
	"github.com/42milez/ProtocolStack/src/device"
	"github.com/42milez/ProtocolStack/src/ethernet"
	"github.com/42milez/ProtocolStack/src/middleware"
	"github.com/42milez/ProtocolStack/src/network"
	"log"
)

func setup() error {
	var dev *device.Device
	var iface *network.Iface
	var err error

	if err = middleware.Setup(); err != nil {
		return err
	}

	// 1. Create a loopback device object.
	// 2. Register the object.
	// 3. Create an IF object.
	// 4. Attach the IF object to the device object.
	dev = ethernet.GenLoopbackDevice()
	device.Register(dev)
	iface = network.GenIF(ethernet.LoopbackIpAddr, ethernet.LoopbackNetmask)
	if err = network.AttachIF(iface, dev); err != nil {
		return err
	}

	// 1. Create a TAP device object.
	// 2. Register the object.
	// 3. Create an IF object.
	// 4. Attach the IF to the device object.
	if dev, err = ethernet.GenTapDevice("tap0", "00:00:5e:00:53:01"); err != nil {
		return err
	}
	device.Register(dev)
	iface = network.GenIF("192.0.2.2", "255.255.255.0")
	if err = network.AttachIF(iface, dev); err != nil {
		return err
	}

	// IFオブジェクト生成
	// ...

	// TAPデバイスオブジェクトとIFオブジェクトの紐付け
	// ...

	// TAP IFオブジェクトをデフォルトゲートウェイとして登録
	// ...

	// デバイスとサブスレッド（goroutine）を生成
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
