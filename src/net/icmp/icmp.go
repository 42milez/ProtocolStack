package icmp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	psErr "github.com/42milez/ProtocolStack/src/error"
	psLog "github.com/42milez/ProtocolStack/src/log"
	"github.com/42milez/ProtocolStack/src/mw"
	"github.com/42milez/ProtocolStack/src/net/ip"
	"github.com/42milez/ProtocolStack/src/repo"
)

const IcmpHdrLen = 8 // byte
const IcmpTypeEchoReply = 0x00
const IcmpTypeEcho = 0x08

func IcmpReceive(payload []byte, dst [mw.V4AddrLen]byte, src [mw.V4AddrLen]byte, dev mw.IDevice) psErr.E {
	if len(payload) < IcmpHdrLen {
		psLog.E(fmt.Sprintf("ICMP header length is too short: %d bytes", len(payload)))
		return psErr.InvalidPacket
	}

	buf := bytes.NewBuffer(payload)
	hdr := IcmpHdr{}
	if err := binary.Read(buf, binary.BigEndian, &hdr); err != nil {
		return psErr.ReadFromBufError
	}

	checksum1 := uint16(payload[2])<<8 | uint16(payload[3])
	payload[2] = 0x00 // assign 0 to Checksum field (16bit)
	payload[3] = 0x00
	if checksum2 := mw.Checksum(payload); checksum2 != checksum1 {
		psLog.E(fmt.Sprintf("Checksum mismatch: Expect = 0x%04x, Actual = 0x%04x", checksum1, checksum2))
		return psErr.ChecksumMismatch
	}

	psLog.I("Incoming ICMP packet")
	dumpIcmpPacket(&hdr, payload[IcmpHdrLen:])

	switch hdr.Type {
	case IcmpTypeEcho:
		s := mw.IP(src[:])
		d := mw.IP(dst[:])
		iface := repo.IfaceRepo.Lookup(dev, mw.V4AddrFamily)
		if !iface.Unicast.EqualV4(dst) {
			d = iface.Unicast
		}
		msg := &mw.IcmpTxMessage{
			Type:    IcmpTypeEchoReply,
			Code:    hdr.Code,
			Content: hdr.Content,
			Payload: payload[IcmpHdrLen:],
			Src:     d,
			Dst:     s,
		}
		mw.IcmpTxCh <- msg
		//if err := IcmpSend(IcmpTypeEchoReply, hdr.Code, hdr.Content, payload[IcmpHdrLen:], d, s); err != psErr.OK {
		//	return psErr.Error
		//}
	}

	return psErr.OK
}

func IcmpSend(typ uint8, code uint8, content uint32, payload []byte, src mw.IP, dst mw.IP) psErr.E {
	hdr := IcmpHdr{
		Type:    typ,
		Code:    code,
		Content: content,
	}

	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, &hdr); err != nil {
		return psErr.WriteToBufError
	}
	if err := binary.Write(buf, binary.BigEndian, &payload); err != nil {
		return psErr.WriteToBufError
	}

	packet := buf.Bytes()
	hdr.Checksum = mw.Checksum(packet)
	packet[2] = uint8((hdr.Checksum & 0xff00) >> 8)
	packet[3] = uint8(hdr.Checksum & 0x00ff)

	psLog.I("Outgoing ICMP packet")
	dumpIcmpPacket(&hdr, payload)

	mw.IpTxCh <- &mw.IpMessage{
		ProtoNum: ip.ICMP,
		Packet:   packet,
		Dst:      dst,
		Src:      src,
	}

	return psErr.OK
}

func dumpIcmpPacket(hdr *IcmpHdr, payload []byte) {
	psLog.I(fmt.Sprintf("\ttype:     %s (%d)", icmpTypes[hdr.Type], hdr.Type))
	psLog.I(fmt.Sprintf("\tcode:     %d", hdr.Code))
	psLog.I(fmt.Sprintf("\tchecksum: 0x%04x", hdr.Checksum))

	switch hdr.Type {
	case IcmpTypeEchoReply:
	case IcmpTypeEcho:
		psLog.I(fmt.Sprintf("\tid:       %d", (hdr.Content&0xffff0000)>>16))
		psLog.I(fmt.Sprintf("\tseq:      %d", hdr.Content&0x0000ffff))
	default:
		psLog.I(fmt.Sprintf("\tcontent:  %02x %02x %02x %02x",
			uint8((hdr.Content&0xf000)>>24),
			uint8((hdr.Content&0x0f00)>>16),
			uint8((hdr.Content&0x00f0)>>8),
			uint8(hdr.Content&0x000f)))
	}

	s := "\tpayload:  "
	for i, v := range payload {
		s += fmt.Sprintf("%02x ", v)
		if (i+1)%20 == 0 {
			psLog.I(s)
			s = "\t\t  "
		}
	}
}

func StartService() {
	go func() {
		for {
			msg := <-mw.IcmpRxCh
			if err := IcmpReceive(msg.Payload, msg.Dst, msg.Src, msg.Dev); err != psErr.OK {
				return
			}
		}
	}()
	go func() {
		for {
			msg := <-mw.IcmpTxCh
			if err := IcmpSend(msg.Type, msg.Code, msg.Content, msg.Payload, msg.Src, msg.Dst); err != psErr.OK {
				return
			}
		}
	}()
}
