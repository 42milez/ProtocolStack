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
	psLog.I(fmt.Sprintf("\tdst:  %s", hdr.Dst))
	psLog.I(fmt.Sprintf("\tsrc:  %s", hdr.Src))
	psLog.I(fmt.Sprintf("\ttype: 0x%04x (%s)", uint16(hdr.Type), hdr.Type))
}

func ReadFrame(fd int, addr EthAddr, sc psSyscall.ISyscall) (*Packet, psErr.E) {
	// TODO: make buf static variable to reuse
	buf := make([]byte, EthFrameSizeMax)

	fsize, err := sc.Read(fd, buf)
	if err != nil {
		psLog.E(fmt.Sprintf("syscall.Read() failed: %s", err))
		return nil, psErr.Error
	}

	if fsize < EthHeaderSize {
		psLog.E("Ethernet header length is too short")
		psLog.E(fmt.Sprintf("\tlength: %v bytes", fsize))
		return nil, psErr.Error
	}

	hdr := EthHeader{}
	if err := binary.Read(bytes.NewBuffer(buf), binary.BigEndian, &hdr); err != nil {
		psLog.E(fmt.Sprintf("binary.Read() failed: %s", err))
		return nil, psErr.Error
	}

	if !hdr.Dst.Equal(addr) {
		if !hdr.Dst.Equal(EthAddrBroadcast) {
			return nil, psErr.NoDataToRead
		}
	}

	psLog.I("Incoming ethernet frame")
	EthDump(&hdr)

	packet := &Packet{
		Type:    hdr.Type,
		Payload: buf[EthHeaderSize:],
	}

	return packet, psErr.OK
}

func WriteFrame(fd int, dst EthAddr, src EthAddr, typ EthType, payload []byte) psErr.E {
	hdr := EthHeader{
		Dst:  dst,
		Src:  src,
		Type: typ,
	}

	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, &hdr); err != nil {
		psLog.E(fmt.Sprintf("binary.Write() failed: %s", err))
		return psErr.Error
	}
	if err := binary.Write(buf, binary.BigEndian, &payload); err != nil {
		psLog.E(fmt.Sprintf("binary.Write() failed: %s", err))
		return psErr.Error
	}

	if fsize := buf.Len(); fsize < EthFrameSizeMin {
		pad := make([]byte, EthFrameSizeMin-fsize)
		if err := binary.Write(buf, binary.BigEndian, &pad); err != nil {
			psLog.E(fmt.Sprintf("binary.Write() failed: %s", err))
			return psErr.Error
		}
	}

	psLog.I("Outgoing Ethernet frame")
	psLog.I(fmt.Sprintf("\tdst:     %s", hdr.Dst))
	psLog.I(fmt.Sprintf("\tsrc:     %s", hdr.Src))
	psLog.I(fmt.Sprintf("\ttype:    %s", hdr.Type))
	s := "\tpayload: "
	for i, v := range payload {
		s += fmt.Sprintf("%02x ", v)
		if (i+1)%10 == 0 {
			psLog.I(s)
			s = "\t\t "
		}
	}

	if n, err := psSyscall.Syscall.Write(fd, buf.Bytes()); err != nil {
		psLog.E(fmt.Sprintf("syscall.Write() failed: %s", err))
		return psErr.Error
	} else {
		psLog.I(fmt.Sprintf("Ethernet frame has been written: %d bytes", n))
	}

	return psErr.OK
}
