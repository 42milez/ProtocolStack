package network

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

func (f AddrFamily) String() string {
	switch f {
	case V4:
		return "FAMILY_V4"
	case V6:
		return "FAMILY_V6"
	default:
		return "UNKNOWN"
	}
}

// IpInputHandler handles incoming datagram.
func IpInputHandler(data []uint8, dev device.Device) {}
