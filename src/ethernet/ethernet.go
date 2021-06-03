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
const EthFrameLenMax = 1514
const EthFrameLenMin = 60
const EthHdrLen = 14
const EthPayloadLenMax = EthFrameLenMax - EthHdrLen
const EthPayloadLenMin = EthFrameLenMin - EthHdrLen
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

type EthHdr struct {
	Dst  EthAddr
	Src  EthAddr
	Type EthType
}
type EthType uint16

func (v EthType) String() string {
	return ethTypes[v]
}

func EthFrameDump(hdr *EthHdr, payload []byte) {
	psLog.I(fmt.Sprintf("\ttype:    0x%04x (%s)", uint16(hdr.Type), hdr.Type))
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

func ReadEthFrame(fd int, addr EthAddr) (*Packet, psErr.E) {
	flen, err := psSyscall.Syscall.Read(fd, rxBuf)
	if err != nil {
		return nil, psErr.Error
	}

	if flen < EthHdrLen {
		psLog.E(fmt.Sprintf("Ethernet header length is too short: %d bytes", flen))
		return nil, psErr.Error
	}

	psLog.I(fmt.Sprintf("Ethernet frame was received: %d bytes", flen))

	buf := bytes.NewBuffer(rxBuf)
	hdr := EthHdr{}
	if err := binary.Read(buf, binary.BigEndian, &hdr); err != nil {
		return nil, psErr.Error
	}

	if !hdr.Dst.Equal(addr) {
		if !hdr.Dst.Equal(EthAddrBroadcast) {
			return nil, psErr.NoDataToRead
		}
	}

	payload := make([]byte, flen)
	if err := binary.Read(buf, binary.BigEndian, &payload); err != nil {
		return nil, psErr.Error
	}

	psLog.I("Incoming ethernet frame")
	EthFrameDump(&hdr, payload)

	return &Packet{
		Type:    hdr.Type,
		Content: payload,
	}, psErr.OK
}

func WriteEthFrame(fd int, dst EthAddr, src EthAddr, typ EthType, payload []byte) psErr.E {
	hdr := EthHdr{
		Dst:  dst,
		Src:  src,
		Type: typ,
	}

	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, &hdr); err != nil {
		return psErr.Error
	}
	if err := binary.Write(buf, binary.BigEndian, &payload); err != nil {
		return psErr.Error
	}
	if err := pad(buf); err != psErr.OK {
		return psErr.Error
	}
	frame := buf.Bytes()

	psLog.I("Outgoing Ethernet frame")
	EthFrameDump(&hdr, payload)

	if n, err := psSyscall.Syscall.Write(fd, frame); err != nil {
		return psErr.SyscallError
	} else {
		psLog.I(fmt.Sprintf("Ethernet frame was sent: %d bytes (payload: %d bytes)", n, len(payload)))
	}

	return psErr.OK
}

func pad(buf *bytes.Buffer) psErr.E {
	if flen := buf.Len(); flen < EthFrameLenMin {
		padLen := EthFrameLenMin - flen
		pad := make([]byte, padLen)
		if err := binary.Write(buf, binary.BigEndian, &pad); err != nil {
			return psErr.Error
		}
	}
	return psErr.OK
}

// Ethertypes
// https://www.iana.org/assignments/ieee-802-numbers/ieee-802-numbers.xhtml#ieee-802-numbers-1

var ethTypes = map[EthType]string{
	0x0800: "IPv4",
	0x0806: "ARP",
	0x86dd: "IPv6",
}

var rxBuf []byte

func init() {
	rxBuf = make([]byte, EthFrameLenMax)
}
