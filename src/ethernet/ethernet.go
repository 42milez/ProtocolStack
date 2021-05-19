package ethernet

import (
	"bytes"
	"encoding/binary"
	"fmt"
	psBinary "github.com/42milez/ProtocolStack/src/binary"
	psErr "github.com/42milez/ProtocolStack/src/error"
	psLog "github.com/42milez/ProtocolStack/src/log"
	psSyscall "github.com/42milez/ProtocolStack/src/syscall"
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

var EthAddrAny = EthAddr{0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
var EthAddrBroadcast = EthAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}

var endian int

type EthAddr [EthAddrLen]byte

func (v EthAddr) Equal(vv EthAddr) bool {
	return v == vv
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

func EthDump(hdr *EthHeader) {
	psLog.I("\tmac (dst): %d:%d:%d:%d:%d:%d", hdr.Dst[0], hdr.Dst[1], hdr.Dst[2], hdr.Dst[3], hdr.Dst[4], hdr.Dst[5])
	psLog.I("\tmac (src): %d:%d:%d:%d:%d:%d", hdr.Src[0], hdr.Src[1], hdr.Src[2], hdr.Src[3], hdr.Src[4], hdr.Src[5])
	psLog.I("\teth_type:  0x%04x (%s)", hdr.Type.String(), hdr.Type.String())
}

func ReadFrame(dev *Device, sc psSyscall.ISyscall) psErr.Error {
	// TODO: make buf static variable to reuse
	buf1 := make([]byte, EthFrameSizeMax)

	flen, _, errno := sc.Read(dev.Priv.FD, unsafe.Pointer(&buf1), EthFrameSizeMax)
	if errno != 0 {
		psLog.E("SYS_READ failed: %v ", errno)
		return psErr.Error{Code: psErr.CantRead}
	}

	if flen < EthHeaderSize*8 {
		psLog.E("the length of ethernet header is too short")
		psLog.E("\tflen: %v ", errno)
		return psErr.Error{Code: psErr.InvalidHeader}
	}

	hdr := EthHeader{}
	buf2 := bytes.NewBuffer(buf1)
	if err := binary.Read(buf2, binary.BigEndian, &hdr); err != nil {
		return psErr.Error{Code: psErr.CantConvert, Msg: err.Error()}
	}

	if !hdr.Dst.Equal(dev.Addr) {
		if !hdr.Dst.Equal(EthAddrBroadcast) {
			return psErr.Error{Code: psErr.NoDataToRead}
		}
	}

	psLog.I("received an ethernet frame")
	psLog.I("\tdevice:    %v (%v) ", dev.Name, dev.Priv.Name)
	psLog.I("\tlength:    %v ", flen)
	EthDump(&hdr)

	return psErr.Error{Code: psErr.OK}
}

func ntoh16(n uint16) uint16 {
	if endian == psBinary.LittleEndian {
		return swap16(n)
	} else {
		return n
	}
}

func swap16(v uint16) uint16 {
	return (v << 8) | (v >> 8)
}

func init() {
	endian = psBinary.ByteOrder()
}
