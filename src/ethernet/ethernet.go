package ethernet

import (
	e "github.com/42milez/ProtocolStack/src/error"
	l "github.com/42milez/ProtocolStack/src/logger"
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
	l.I("\tmac (dst):   %d:%d:%d:%d:%d:%d", hdr.Dst[0], hdr.Dst[1], hdr.Dst[2], hdr.Dst[3], hdr.Dst[4], hdr.Dst[5])
	l.I("\tmac (src):   %d:%d:%d:%d:%d:%d", hdr.Src[0], hdr.Src[1], hdr.Src[2], hdr.Src[3], hdr.Src[4], hdr.Src[5])
	l.I("\tethertype: 0x%04x (%s)", ntoh16(hdr.Type), hdr.TypeAsString())
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
	if !hdr.EqualDstAddr(dev.Addr.Byte()) {
		if !hdr.EqualDstAddr(EthAddrBroadcast.Byte()) {
			return e.Error{Code: e.NoDataToRead}
		}
	}

	l.I("received an ethernet frame")
	l.I("\tlength: %v ", flen)
	EtherDump(hdr)

	l.I("\tdevice:       %v (%v) ", dev.Name, dev.Priv.Name)
	l.I("\tethertype:    %v (0x%04x) ", hdr.TypeAsString(), ntoh16(hdr.Type))
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
