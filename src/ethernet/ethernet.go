package ethernet

import (
	"bytes"
	"encoding/binary"
	"fmt"
	psBinary "github.com/42milez/ProtocolStack/src/binary"
	psErr "github.com/42milez/ProtocolStack/src/error"
	psLog "github.com/42milez/ProtocolStack/src/log"
	psSyscall "github.com/42milez/ProtocolStack/src/syscall"
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
	psLog.I("\tmac (dst): %v", hdr.Dst.String())
	psLog.I("\tmac (src): %v", hdr.Src.String())
	psLog.I("\teth_type:  0x%04x (%s)", hdr.Type.String(), hdr.Type.String())
}

func ReadFrame(fd int, addr EthAddr, sc psSyscall.ISyscall) psErr.Error {
	// TODO: make buf static variable to reuse
	buf := make([]byte, EthFrameSizeMax)

	fsize, err := sc.Read(fd, buf)
	if err != nil {
		psLog.E("SYS_READ failed: %v ", err)
		return psErr.Error{Code: psErr.CantRead}
	}

	if fsize < EthHeaderSize {
		psLog.E("the length of ethernet header is too short")
		psLog.E("\tfsize: %v bytes", fsize)
		return psErr.Error{Code: psErr.InvalidHeader}
	}

	hdr := EthHeader{}
	if err := binary.Read(bytes.NewBuffer(buf), binary.BigEndian, &hdr); err != nil {
		return psErr.Error{Code: psErr.CantConvert, Msg: err.Error()}
	}

	if !hdr.Dst.Equal(addr) {
		if !hdr.Dst.Equal(EthAddrBroadcast) {
			return psErr.Error{Code: psErr.NoDataToRead}
		}
	}

	psLog.I("received an ethernet frame")
	psLog.I("\tfsize:     %v bytes", fsize)
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
