package main

import (
	"github.com/42milez/ProtocolStack/src/device"
	"github.com/42milez/ProtocolStack/src/ethernet"
	"github.com/42milez/ProtocolStack/src/ipv4"
	"github.com/42milez/ProtocolStack/src/middleware"
	"log"
)

func setup() error {
	var err error

	if err = middleware.Setup(); err != nil {
		return err
	}

	// 1. Create a loopback object.
	// 2. Register the object.
	// 3. Create an IF object.
	// 4. Attach the IF object to the loopback object.
	dev := ethernet.GenLoopbackDevice()
	device.RegisterDevice(dev)
	iface := ipv4.GenIF(ethernet.LoopbackIpAddr, ethernet.LoopbackNetmask)
	ipv4.RegisterIF(iface, dev)

	// TAPデバイスオブジェクト生成
	// ...

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
