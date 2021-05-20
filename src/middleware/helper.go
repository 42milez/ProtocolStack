package middleware

import (
	"github.com/42milez/ProtocolStack/src/ethernet"
	"github.com/42milez/ProtocolStack/src/network"
	psSyscall "github.com/42milez/ProtocolStack/src/syscall"
	"strconv"
)

// GenIface generates Iface.
func GenIface(unicast string, netmask string) *Iface {
	iface := &Iface{
		Family:    network.FamilyV4,
		Unicast:   network.ParseIP(unicast),
		Netmask:   network.ParseIP(netmask),
		Broadcast: make(network.IP, 0),
		Network:   make(network.IP, 0),
	}

	//unicastUint32 := binary.BigEndian.Uint32(iface.Unicast)
	//netmaskUint32 := binary.BigEndian.Uint32(iface.Netmask)
	//broadcastUint32 := (unicastUint32 & netmaskUint32) | ^netmaskUint32

	//binary.BigEndian.PutUint32(iface.Broadcast, broadcastUint32)
	//binary.BigEndian.PutUint32(iface.Network, unicastUint32 & netmaskUint32)

	return iface
}

// GenLoopbackDevice generates loopback device object.
func GenLoopbackDevice() *ethernet.LoopbackDevice {
	dev := &ethernet.LoopbackDevice{
		Device: ethernet.Device{
			Name:      "net" + strconv.Itoa(len(devices)),
			Type:      ethernet.DevTypeLoopback,
			MTU:       ethernet.LoopbackMTU,
			HeaderLen: 0,
			FLAG:      ethernet.DevFlagLoopback,
			Syscall:   &psSyscall.Syscall{},
		},
	}
	return dev
}

// GenTapDevice generates TAP device object.
func GenTapDevice(index uint8, addr ethernet.EthAddr) *ethernet.TapDevice {
	return &ethernet.TapDevice{
		Device: ethernet.Device{
			Type:      ethernet.DevTypeEthernet,
			MTU:       ethernet.EthPayloadSizeMax,
			FLAG:      ethernet.DevFlagBroadcast | ethernet.DevFlagNeedArp,
			HeaderLen: ethernet.EthHeaderSize,
			Addr:      addr,
			Broadcast: ethernet.EthAddrBroadcast,
			Priv: ethernet.Privilege{
				FD:   -1,
				Name: "tap" + strconv.Itoa(int(index)),
			},
			Syscall: &psSyscall.Syscall{},
		},
	}
}
