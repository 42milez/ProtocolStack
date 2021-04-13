package network

import (
	"github.com/42milez/ProtocolStack/src/device"
)

const IpHeaderSizeMin = 20
const IpHeaderSizeMax = 60

// AddrFamily is IP address family.
type AddrFamily int

const (
	FamilyV4 AddrFamily = iota
	FamilyV6
)

func (f AddrFamily) String() string {
	switch f {
	case FamilyV4:
		return "FAMILY_V4"
	case FamilyV6:
		return "FAMILY_V6"
	default:
		return "UNKNOWN"
	}
}

// IpInputHandler handles incoming datagram.
func IpInputHandler(data []uint8, dev device.Device) {}
