package mw

import (
	"bytes"
	"encoding/binary"
	"fmt"
	psErr "github.com/42milez/ProtocolStack/src/error"
	psLog "github.com/42milez/ProtocolStack/src/log"
	psSyscall "github.com/42milez/ProtocolStack/src/syscall"
)

const EthHdrLen = 14
const EthAddrLen = 6
const EthFrameLenMax = 1514
const EthFrameLenMin = 60
const EthPayloadLenMax = EthFrameLenMax - EthHdrLen
const EthPayloadLenMin = EthFrameLenMin - EthHdrLen

var EthAny = EthAddr{0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
var EthBroadcast = EthAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}

var rxBuf []byte

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

func ReadFrame(fd int, addr EthAddr) (*EthMessage, psErr.E) {
	flen, err := psSyscall.Syscall.Read(fd, rxBuf)
	if err != nil {
		return nil, psErr.Error
	}

	if flen < EthHdrLen {
		psLog.E(fmt.Sprintf("ethernet header length is too short: %d bytes", flen))
		return nil, psErr.Error
	}

	buf := bytes.NewBuffer(rxBuf)
	hdr := EthHdr{}
	if err := binary.Read(buf, binary.BigEndian, &hdr); err != nil {
		return nil, psErr.ReadFromBufError
	}

	if !hdr.Dst.Equal(addr) {
		if !hdr.Dst.Equal(EthBroadcast) {
			return nil, psErr.NoDataToRead
		}
	}

	payload := make([]byte, flen-EthHdrLen)
	if err := binary.Read(buf, binary.BigEndian, &payload); err != nil {
		return nil, psErr.ReadFromBufError
	}

	psLog.D(fmt.Sprintf("incoming ethernet frame (%d bytes)", flen), dump(&hdr, payload)...)

	return &EthMessage{
		Type:    hdr.Type,
		Content: payload,
	}, psErr.OK
}

func WriteFrame(fd int, dst EthAddr, src EthAddr, typ EthType, payload []byte) psErr.E {
	hdr := EthHdr{
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

	psLog.D(fmt.Sprintf("outgoing ethernet frame (%d bytes)", EthHdrLen+len(payload)), dump(&hdr, payload)...)

	if _, err := psSyscall.Syscall.Write(fd, frame); err != nil {
		return psErr.SyscallError
	}

	return psErr.OK
}

func dump(hdr *EthHdr, payload []byte) (ret []string) {
	ret = append(ret, fmt.Sprintf("type:    %s (0x%04x)", hdr.Type, uint16(hdr.Type)))
	ret = append(ret, fmt.Sprintf("dst:     %s", hdr.Dst))
	ret = append(ret, fmt.Sprintf("src:     %s", hdr.Src))
	s := "payload: "
	for i, v := range payload {
		s += fmt.Sprintf("%02x ", v)
		if (i+1)%20 == 0 {
			s += "\n                                 "
		}
	}
	ret = append(ret, s)
	return
}

func pad(buf *bytes.Buffer) psErr.E {
	if flen := buf.Len(); flen < EthFrameLenMin {
		padLen := EthFrameLenMin - flen
		pad := make([]byte, padLen)
		if err := binary.Write(buf, binary.BigEndian, &pad); err != nil {
			return psErr.WriteToBufError
		}
	}
	return psErr.OK
}

func init() {
	rxBuf = make([]byte, EthFrameLenMax)
}
