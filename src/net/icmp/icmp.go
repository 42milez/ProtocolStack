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
	"github.com/42milez/ProtocolStack/src/worker"
	"sync"
)

const Echo = 0x08
const EchoReply = 0x00
const HdrLen = 8 // byte
const xChBufSize = 5

var RcvRxCh chan *worker.Message
var RcvTxCh chan *worker.Message
var SndRxCh chan *worker.Message
var SndTxCh chan *worker.Message

// ICMP Type Numbers
// https://www.iana.org/assignments/icmp-parameters/icmp-parameters.xhtml#icmp-parameters-types

var types = map[uint8]string{
	0: "Echo Reply",
	// 1-2: Unassigned
	3: "Destination Unreachable",
	4: "Source Quench (Deprecated)",
	5: "Redirect",
	6: "Alternate Host Address (Deprecated)",
	// 7: Unassigned
	8:  "Echo",
	9:  "Router Advertisement",
	10: "Router Solicitation",
	11: "Time Exceeded",
	12: "Parameter Problem",
	13: "Timestamp",
	14: "Timestamp Reply",
	15: "Information Request (Deprecated)",
	16: "Information Reply (Deprecated)",
	17: "Address Mask Request (Deprecated)",
	18: "Address Mask Reply (Deprecated)",
	19: "Reserved (for Security)",
	// 20-29: Reserved (for Robustness Experiment)
	30: "Traceroute (Deprecated)",
	31: "Datagram Conversion Error (Deprecated)",
	32: "Mobile Host Redirect (Deprecated)",
	33: "IPv6 Where-Are-You (Deprecated)",
	34: "IPv6 I-Am-Here (Deprecated)",
	35: "Mobile Registration Request (Deprecated)",
	36: "Mobile Registration Reply (Deprecated)",
	37: "Domain Name Request (Deprecated)",
	38: "Domain Name Reply (Deprecated)",
	39: "SKIP (Deprecated)",
	40: "Photuris",
	41: "ICMP messages utilized by experimental mobility protocols such as Seamoby",
	42: "Extended Echo Request",
	43: "Extended Echo Reply",
	// 44-252: Unassigned
	253: "RFC3692-style Experiment 1",
	254: "RFC3692-style Experiment 2",
	// 255: Reserved
}

type Hdr struct {
	Type     uint8
	Code     uint8
	Checksum uint16
	Content  uint32
}

func Receive(payload []byte, dst [mw.V4AddrLen]byte, src [mw.V4AddrLen]byte, dev mw.IDevice) psErr.E {
	if len(payload) < HdrLen {
		psLog.E(fmt.Sprintf("ICMP header length is too short: %d bytes", len(payload)))
		return psErr.InvalidPacket
	}

	buf := bytes.NewBuffer(payload)
	hdr := Hdr{}
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
	dump(&hdr, payload[HdrLen:])

	switch hdr.Type {
	case Echo:
		s := mw.IP(src[:])
		d := mw.IP(dst[:])
		iface := repo.IfaceRepo.Lookup(dev, mw.V4AddrFamily)
		if !iface.Unicast.EqualV4(dst) {
			d = iface.Unicast
		}
		msg := &mw.IcmpTxMessage{
			Type:    EchoReply,
			Code:    hdr.Code,
			Content: hdr.Content,
			Payload: payload[HdrLen:],
			Src:     d,
			Dst:     s,
		}
		mw.IcmpTxCh <- msg
		//if err := Send(EchoReply, hdr.Code, hdr.Content, payload[EthHdrLen:], d, s); err != psErr.OK {
		//	return psErr.Error
		//}
	}

	return psErr.OK
}

func Send(typ uint8, code uint8, content uint32, payload []byte, src mw.IP, dst mw.IP) psErr.E {
	hdr := Hdr{
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
	dump(&hdr, payload)

	mw.IpTxCh <- &mw.IpMessage{
		ProtoNum: ip.ICMP,
		Packet:   packet,
		Dst:      dst,
		Src:      src,
	}

	return psErr.OK
}

func StartService(wg *sync.WaitGroup) {
	wg.Add(2)
	go receiver(wg)
	go sender(wg)
}

func dump(hdr *Hdr, payload []byte) {
	psLog.I(fmt.Sprintf("\ttype:     %s (%d)", types[hdr.Type], hdr.Type))
	psLog.I(fmt.Sprintf("\tcode:     %d", hdr.Code))
	psLog.I(fmt.Sprintf("\tchecksum: 0x%04x", hdr.Checksum))

	switch hdr.Type {
	case Echo:
	case EchoReply:
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

func receiver(wg *sync.WaitGroup) {
	defer wg.Done()

	RcvTxCh <- &worker.Message{
		Current: worker.Running,
	}

	for {
		select {
		case msg := <-RcvRxCh:
			if msg.Desired == worker.Stopped {
				RcvTxCh <- &worker.Message{
					Current: worker.Stopped,
				}
				return
			}
		case msg := <-mw.IcmpRxCh:
			if err := Receive(msg.Payload, msg.Dst, msg.Src, msg.Dev); err != psErr.OK {
				return
			}
		}
	}
}

func sender(wg *sync.WaitGroup) {
	defer wg.Done()

	SndTxCh <- &worker.Message{
		Current: worker.Running,
	}

	for {
		select {
		case msg := <-SndRxCh:
			if msg.Desired == worker.Stopped {
				SndTxCh <- &worker.Message{
					Current: worker.Stopped,
				}
				return
			}
		case msg := <-mw.IcmpTxCh:
			if err := Send(msg.Type, msg.Code, msg.Content, msg.Payload, msg.Src, msg.Dst); err != psErr.OK {
				return
			}
		}
	}
}

func init() {
	RcvRxCh = make(chan *worker.Message, xChBufSize)
	RcvTxCh = make(chan *worker.Message, xChBufSize)
	SndRxCh = make(chan *worker.Message, xChBufSize)
	SndTxCh = make(chan *worker.Message, xChBufSize)
}
