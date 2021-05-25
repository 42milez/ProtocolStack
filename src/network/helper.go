package network

import (
	"github.com/42milez/ProtocolStack/src/ethernet"
	psSyscall "github.com/42milez/ProtocolStack/src/syscall"
	"strconv"
)

// GenIface generates Iface.
func GenIface(unicast string, netmask string, broadcast string) *Iface {
	iface := &Iface{
		Family:    FamilyV4,
		Unicast:   ParseIP(unicast),
		Netmask:   ParseIP(netmask),
		Broadcast: ParseIP(broadcast),
	}
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
			Name:      "net" + strconv.Itoa(len(devices)),
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
