package ethernet

import (
	"fmt"
	e "github.com/42milez/ProtocolStack/src/error"
	l "github.com/42milez/ProtocolStack/src/log"
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

var EthAddrAny = MAC("00:00:00:00:00:00")
var EthAddrBroadcast = MAC("FF:FF:FF:FF:FF:FF")

var endian int

type EthAddr [EthAddrLen]byte

func (v EthAddr) Equal(vv EthAddr) bool {
	return v == vv
}

func (v EthAddr) EqualByte(vv []byte) bool {
	if len(v) != len(vv) {
		return false
	}
	return v[0] == vv[0] && v[1] == vv[1] && v[2] == vv[2] && v[3] == vv[3] && v[4] == vv[4] && v[5] == vv[5]
}

func (v EthAddr) String() string {
	return fmt.Sprintf("%v:%v:%v:%v:%v:%v", v[0], v[1], v[2], v[3], v[4], v[5])
}

type EthType uint16

func (v EthType) String() string {
	switch ntoh16(uint16(v)) {
	case EthTypeArp:
		return "ARP"
	case EthTypeIpv4:
		return "IPv4"
	case EthTypeIpv6:
		return "IPv6"
	default:
		return "UNKNOWN"
	}
}

type EthHeader struct {
	Dst  EthAddr
	Src  EthAddr
	Type EthType
}

func EtherDump(hdr *EthHeader) {
	l.I("\tmac (dst): %d:%d:%d:%d:%d:%d", hdr.Dst[0], hdr.Dst[1], hdr.Dst[2], hdr.Dst[3], hdr.Dst[4], hdr.Dst[5])
	l.I("\tmac (src): %d:%d:%d:%d:%d:%d", hdr.Src[0], hdr.Src[1], hdr.Src[2], hdr.Src[3], hdr.Src[4], hdr.Src[5])
	l.I("\teth_type:  0x%04x (%s)", hdr.Type.String(), hdr.Type.String())
}

func ReadFrame(dev *Device) e.Error {
	var buf [EthFrameSizeMax]byte

	flen, _, errno := syscall.Syscall(
		syscall.SYS_READ,
		uintptr(dev.Priv.FD),
		uintptr(unsafe.Pointer(&buf)),
		uintptr(EthFrameSizeMax))

	if errno != 0 {
		l.E("SYS_READ failed: %v ", errno)
		return e.Error{Code: e.CantRead}
	}

	if flen < EthHeaderSize*8 {
		l.E("the length of ethernet header is too short")
		return e.Error{Code: e.InvalidHeader}
	}

	hdr := (*EthHeader)(unsafe.Pointer(&buf))
	if !hdr.Dst.EqualByte(dev.Addr.Byte()) {
		if !hdr.Dst.EqualByte(EthAddrBroadcast.Byte()) {
			return e.Error{Code: e.NoDataToRead}
		}
	}

	l.I("received an ethernet frame")
	l.I("\tlength: %v ", flen)
	EtherDump(hdr)

	l.I("\tdevice:       %v (%v) ", dev.Name, dev.Priv.Name)
	l.I("\teth_type:     %v (0x%04x) ", hdr.Type.String(), hdr.Type)
	l.I("\tframe length: %v ", flen)

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
