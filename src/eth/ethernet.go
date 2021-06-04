package eth

import (
	"bytes"
	"encoding/binary"
	"fmt"
	psErr "github.com/42milez/ProtocolStack/src/error"
	psLog "github.com/42milez/ProtocolStack/src/log"
	psSyscall "github.com/42milez/ProtocolStack/src/syscall"
)

const AddrLen = 6
const FrameLenMax = 1514
const FrameLenMin = 60
const HdrLen = 14
const PayloadLenMax = FrameLenMax - HdrLen
const PayloadLenMin = FrameLenMin - HdrLen
const ARP Type = 0x0806
const IPv4 Type = 0x0800
const IPv6 Type = 0x86dd

var Any = Addr{0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
var Broadcast = Addr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}

// Ethertypes
// https://www.iana.org/assignments/ieee-802-numbers/ieee-802-numbers.xhtml#ieee-802-numbers-1

var ethTypes = map[Type]string{
	0x0800: "IPv4",
	0x0806: "ARP",
	0x86dd: "IPv6",
}

var rxBuf []byte

type Addr [AddrLen]byte

func (v Addr) Equal(vv Addr) bool {
	return v == vv
}

func (v Addr) String() string {
	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", v[0], v[1], v[2], v[3], v[4], v[5])
}

type Hdr struct {
	Dst  Addr
	Src  Addr
	Type Type
}
type Type uint16

func (v Type) String() string {
	return ethTypes[v]
}

func ReadEthFrame(fd int, addr Addr) (*Packet, psErr.E) {
	flen, err := psSyscall.Syscall.Read(fd, rxBuf)
	if err != nil {
		return nil, psErr.Error
	}

	if flen < HdrLen {
		psLog.E(fmt.Sprintf("Ethernet header length is too short: %d bytes", flen))
		return nil, psErr.Error
	}

	psLog.I(fmt.Sprintf("Ethernet frame arrived: %d bytes", flen))

	buf := bytes.NewBuffer(rxBuf)
	hdr := Hdr{}
	if err := binary.Read(buf, binary.BigEndian, &hdr); err != nil {
		return nil, psErr.ReadFromBufError
	}

	if !hdr.Dst.Equal(addr) {
		if !hdr.Dst.Equal(Broadcast) {
			return nil, psErr.NoDataToRead
		}
	}

	payload := make([]byte, flen)
	if err := binary.Read(buf, binary.BigEndian, &payload); err != nil {
		return nil, psErr.ReadFromBufError
	}

	psLog.I("Incoming eth frame")
	dumpEthFrame(&hdr, payload)

	return &Packet{
		Type:    hdr.Type,
		Content: payload,
	}, psErr.OK
}

func WriteEthFrame(fd int, dst Addr, src Addr, typ Type, payload []byte) psErr.E {
	hdr := Hdr{
		Dst:  dst,
		Src:  src,
		Type: typ,
	}

	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, &hdr); err != nil {
		return psErr.WriteToBufError
	}
	if err := binary.Write(buf, binary.BigEndian, &payload); err != nil {
		return psErr.WriteToBufError
	}
	if err := pad(buf); err != psErr.OK {
		return psErr.Error
	}
	frame := buf.Bytes()

	psLog.I("Outgoing Ethernet frame")
	dumpEthFrame(&hdr, payload)

	if n, err := psSyscall.Syscall.Write(fd, frame); err != nil {
		return psErr.SyscallError
	} else {
		psLog.I(fmt.Sprintf("Ethernet frame was sent: %d bytes (payload: %d bytes)", n, len(payload)))
	}

	return psErr.OK
}

func dumpEthFrame(hdr *Hdr, payload []byte) {
	psLog.I(fmt.Sprintf("\ttype:    %s (0x%04x)", hdr.Type, uint16(hdr.Type)))
	psLog.I(fmt.Sprintf("\tdst:     %s", hdr.Dst))
	psLog.I(fmt.Sprintf("\tsrc:     %s", hdr.Src))

	s := "\tpayload: "
	for i, v := range payload {
		s += fmt.Sprintf("%02x ", v)
		if (i+1)%20 == 0 {
			psLog.I(s)
			s = "\t\t "
		}
	}
}

func pad(buf *bytes.Buffer) psErr.E {
	if flen := buf.Len(); flen < FrameLenMin {
		padLen := FrameLenMin - flen
		pad := make([]byte, padLen)
		if err := binary.Write(buf, binary.BigEndian, &pad); err != nil {
			return psErr.WriteToBufError
		}
	}
	return psErr.OK
}

func init() {
	rxBuf = make([]byte, FrameLenMax)
}
