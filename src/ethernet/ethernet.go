package ethernet

import (
	"fmt"
	"github.com/42milez/ProtocolStack/src/device"
	e "github.com/42milez/ProtocolStack/src/error"
	"log"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

const EthAddrLen = 6
const EthHeaderSize = 14
const EthFrameSizeMin = 60
const EthFrameSizeMax = 1514
const EthPayloadSizeMin = EthFrameSizeMin - EthHeaderSize
const EthPayloadSizeMax = EthFrameSizeMax - EthHeaderSize

const EthTypeArp = 0x0806
const EthTypeIpv4 = 0x0800
const EthTypeIpv6 = 0x86dd

var EthAddrAny = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
var EthAddrBroadcast = []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}

var endian int

type MAC string

func (mac MAC) Byte() ([]byte, error) {
	t := strings.Split(string(mac), ":")
	p := make([]byte, EthAddrLen)
	for i := 0; i < EthAddrLen; i++ {
		var n uint64
		var err error
		if n, err = strconv.ParseUint(t[i], 16, 8); err != nil {
			return nil, err
		}
		if n > 0xff {
			return nil, fmt.Errorf("invalid MAC address")
		}
		p[i] = byte(n)
	}
	return p, nil
}

type EthHeader struct {
	Dst  [EthAddrLen]byte
	Src  [EthAddrLen]byte
	Type uint16
}

func (h EthHeader) EqualDstAddr(v []byte) bool {
	return h.Dst[0] != v[0] &&
		h.Dst[1] != v[1] &&
		h.Dst[2] != v[2] &&
		h.Dst[3] != v[3] &&
		h.Dst[4] != v[4] &&
		h.Dst[5] != v[5]
}

func (h EthHeader) TypeAsString() string {
	switch ntoh16(h.Type) {
	case EthTypeArp:
		return "ARP"
	case EthTypeIpv4:
		return "IPv4"
	case EthTypeIpv6:
		return "IPv6"
	default:
		return "Unknown Type"
	}
}

func EtherDump(hdr *EthHeader) {
	log.Printf("\tmac (dst):   %d:%d:%d:%d:%d:%d", hdr.Dst[0], hdr.Dst[1], hdr.Dst[2], hdr.Dst[3], hdr.Dst[4], hdr.Dst[5])
	log.Printf("\tmac (src):   %d:%d:%d:%d:%d:%d", hdr.Src[0], hdr.Src[1], hdr.Src[2], hdr.Src[3], hdr.Src[4], hdr.Src[5])
	log.Printf("\tethertype: 0x%04x (%s)", ntoh16(hdr.Type), hdr.TypeAsString())
}

func ReadFrame(dev *device.Device) e.Error {
	var buf [EthFrameSizeMax]byte

	flen, _, errno := syscall.Syscall(
		syscall.SYS_READ,
		uintptr(dev.Priv.FD),
		uintptr(unsafe.Pointer(&buf)),
		uintptr(EthFrameSizeMax))

	if errno != 0 {
		log.Printf("SYS_READ failed: %v\n", errno)
		return e.Error{Code: e.CantRead}
	}

	if flen < EthHeaderSize*8 {
		log.Println("the length of ethernet header is too short")
		return e.Error{Code: e.InvalidHeader}
	}

	hdr := (*EthHeader)(unsafe.Pointer(&buf))
	if !hdr.EqualDstAddr(dev.Addr) {
		if !hdr.EqualDstAddr(EthAddrBroadcast) {
			return e.Error{Code: e.NoDataToRead}
		}
	}

	log.Println("received an ethernet frame")
	log.Printf("\tlength: %v\n", flen)
	EtherDump(hdr)

	log.Printf("\tdevice:       %v (%v)\n", dev.Name, dev.Priv.Name)
	log.Printf("\tethertype:    %v (0x%04x)\n", hdr.TypeAsString(), ntoh16(hdr.Type))
	log.Printf("\tframe length: %v\n", flen)

	return e.Error{Code: e.OK}
}

const bigEndian = 4321
const littleEndian = 1234

func byteOrder() int {
	x := 0x0100
	p := unsafe.Pointer(&x)
	if 0x01 == *(*byte)(p) {
		return bigEndian
	} else {
		return littleEndian
	}
}

func ntoh16(n uint16) uint16 {
	if endian == littleEndian {
		return swap16(n)
	} else {
		return n
	}
}

func swap16(v uint16) uint16 {
	return (v << 8) | (v >> 8)
}

func init() {
	endian = byteOrder()
}
