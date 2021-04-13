package ipv4

import (
	"github.com/42milez/ProtocolStack/src/device"
)

const IpHeaderSizeMin = 20
const IpHeaderSizeMax = 60

// AddrFamily is IP address family.
type AddrFamily int

const (
	V4 AddrFamily = iota
	V6
)

func to32(bytes...) uint32 {
	binary.
}

// IpInputHandler handles incoming datagram.
func IpInputHandler(data []uint8, dev device.Device) {}
