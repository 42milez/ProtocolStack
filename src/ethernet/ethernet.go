package ethernet

import (
	"bytes"
	"encoding/binary"
	"fmt"
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

const EthTypeArp EthType = 0x0806
const EthTypeIpv4 EthType = 0x0800
const EthTypeIpv6 EthType = 0x86dd

var EthAddrAny = EthAddr{0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
var EthAddrBroadcast = EthAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}

type EthAddr [EthAddrLen]byte

func (v EthAddr) Equal(vv EthAddr) bool {
	return v == vv
}

func (v EthAddr) String() string {
	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", v[0], v[1], v[2], v[3], v[4], v[5])
}

type EthType uint16

func (v EthType) String() string {
	switch v {
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
	psLog.I("\teth_type:  0x%04x (%s)", uint16(hdr.Type), hdr.Type.String())
}

func ReadFrame(fd int, addr EthAddr, sc psSyscall.ISyscall) (*Packet, psErr.Error) {
	// TODO: make buf static variable to reuse
	buf := make([]byte, EthFrameSizeMax)

	fsize, err := sc.Read(fd, buf)
	if err != nil {
		psLog.E("SYS_READ failed: %v ", err)
		return nil, psErr.Error{Code: psErr.CantRead}
	}

	if fsize < EthHeaderSize {
		psLog.E("▶ Ethernet header length too short")
		psLog.E("\tlength: %v bytes", fsize)
		return nil, psErr.Error{Code: psErr.InvalidHeader}
	}

	hdr := EthHeader{}
	if err := binary.Read(bytes.NewBuffer(buf), binary.BigEndian, &hdr); err != nil {
		return nil, psErr.Error{Code: psErr.CantConvert, Msg: err.Error()}
	}

	if !hdr.Dst.Equal(addr) {
		if !hdr.Dst.Equal(EthAddrBroadcast) {
			return nil, psErr.Error{Code: psErr.NoDataToRead}
		}
	}

	psLog.I("▶ Ethernet frame received")
	psLog.I("\tlength:     %v bytes", fsize)
	EthDump(&hdr)

	packet := &Packet{
		Type:    hdr.Type,
		Payload: buf[EthHeaderSize:],
	}

	return packet, psErr.Error{Code: psErr.OK}
}
