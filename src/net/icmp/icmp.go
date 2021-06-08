package icmp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	psErr "github.com/42milez/ProtocolStack/src/error"
	psLog "github.com/42milez/ProtocolStack/src/log"
	"github.com/42milez/ProtocolStack/src/monitor"
	"github.com/42milez/ProtocolStack/src/mw"
	"github.com/42milez/ProtocolStack/src/net/ip"
	"github.com/42milez/ProtocolStack/src/repo"
	"github.com/42milez/ProtocolStack/src/worker"
	"sync"
)

const Echo = 0x08
const EchoReply = 0x00
const HdrLen = 8 // byte
const replyQueueSize = 5
const xChBufSize = 5

var rcvMonCh chan *worker.Message
var rcvSigCh chan *worker.Message
var sndMonCh chan *worker.Message
var sndSigCh chan *worker.Message

var receiverID uint32
var senderID uint32

var ReplyQueue chan *Reply

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

type Reply struct {
	ID  uint16
	Seq uint16
}

type Hdr struct {
	Type     uint8
	Code     uint8
	Checksum uint16
	Content  uint32
}

func Receive(payload []byte, dst [mw.V4AddrLen]byte, src [mw.V4AddrLen]byte, dev mw.IDevice) psErr.E {
	if len(payload) < HdrLen {
		psLog.E(fmt.Sprintf("icmp header length is too short: %d bytes", len(payload)))
		return psErr.InvalidPacket
	}

	buf := bytes.NewBuffer(payload)
	hdr, err := ReadHeader(buf)
	if err != nil {
		return psErr.ReadFromBufError
	}

	checksum1 := uint16(payload[2])<<8 | uint16(payload[3])
	payload[2] = 0x00 // assign 0 to Checksum field (16bit)
	payload[3] = 0x00
	if checksum2 := mw.Checksum(payload); checksum2 != checksum1 {
		psLog.E(fmt.Sprintf("checksum mismatch: Expect = 0x%04x, Actual = 0x%04x", checksum1, checksum2))
		return psErr.ChecksumMismatch
	}

	psLog.D("incoming icmp packet", dump(hdr, payload[HdrLen:])...)

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
	case EchoReply:
		ReplyQueue <- &Reply{
			ID:  uint16((hdr.Content & 0xffff0000) >> 16),
			Seq: uint16(hdr.Content & 0x0000ffff),
		}
	default:
		psLog.E(fmt.Sprintf("unsupported icmp type: %d", hdr.Type))
		return psErr.Error
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

	psLog.D("outgoing icmp packet", dump(&hdr, payload)...)

	mw.IpTxCh <- &mw.IpMessage{
		ProtoNum: ip.ICMP,
		Packet:   packet,
		Dst:      dst,
		Src:      src,
	}

	return psErr.OK
}

func ReadHeader(buf *bytes.Buffer) (hdr *Hdr, err error) {
	hdr = &Hdr{}
	err = binary.Read(buf, binary.BigEndian, hdr)
	return
}

func SplitContent(content uint32) (id uint16, seq uint16) {
	id = uint16((content & 0xffff0000) >> 16)
	seq = uint16(content & 0x0000ffff)
	return
}

func Start(wg *sync.WaitGroup) psErr.E {
	wg.Add(2)
	go receiver(wg)
	go sender(wg)
	psLog.I("icmp service started")
	return psErr.OK
}

func Stop() {
	msg := &worker.Message{
		Desired: worker.Stopped,
	}
	rcvSigCh <- msg
	sndSigCh <- msg
}

func dump(hdr *Hdr, payload []byte) (ret []string) {
	ret = append(ret, fmt.Sprintf("type:     %s (%d)", types[hdr.Type], hdr.Type))
	ret = append(ret, fmt.Sprintf("code:     %d", hdr.Code))
	ret = append(ret, fmt.Sprintf("checksum: 0x%04x", hdr.Checksum))

	switch hdr.Type {
	case Echo:
	case EchoReply:
		ret = append(ret, fmt.Sprintf("id:       %d", (hdr.Content&0xffff0000)>>16))
		ret = append(ret, fmt.Sprintf("seq:      %d", hdr.Content&0x0000ffff))
	default:
		ret = append(ret, fmt.Sprintf("content:  %02x %02x %02x %02x",
			uint8((hdr.Content&0xf000)>>24),
			uint8((hdr.Content&0x0f00)>>16),
			uint8((hdr.Content&0x00f0)>>8),
			uint8(hdr.Content&0x000f)))
	}

	s := "payload:  "
	for i, v := range payload {
		s += fmt.Sprintf("%02x ", v)
		if (i+1)%20 == 0 {
			s += "\n                                  "
		}
	}
	ret = append(ret, s)

	return
}

func receiver(wg *sync.WaitGroup) {
	defer wg.Done()

	rcvMonCh <- &worker.Message{
		ID:      receiverID,
		Current: worker.Running,
	}

	for {
		select {
		case msg := <-rcvSigCh:
			if msg.Desired == worker.Stopped {
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

	sndMonCh <- &worker.Message{
		ID:      senderID,
		Current: worker.Running,
	}

	for {
		select {
		case msg := <-sndSigCh:
			if msg.Desired == worker.Stopped {
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
	rcvMonCh = make(chan *worker.Message, xChBufSize)
	rcvSigCh = make(chan *worker.Message, xChBufSize)
	receiverID = monitor.Register("ICMP Receiver", rcvMonCh, rcvSigCh)

	sndMonCh = make(chan *worker.Message, xChBufSize)
	sndSigCh = make(chan *worker.Message, xChBufSize)
	senderID = monitor.Register("ICMP Sender", sndMonCh, sndSigCh)

	ReplyQueue = make(chan *Reply, replyQueueSize)
}
