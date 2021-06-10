package tcp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	psErr "github.com/42milez/ProtocolStack/src/error"
	psLog "github.com/42milez/ProtocolStack/src/log"
	"github.com/42milez/ProtocolStack/src/monitor"
	"github.com/42milez/ProtocolStack/src/mw"
	"github.com/42milez/ProtocolStack/src/worker"
	"sync"
)

const HdrLenMax = 60 // byte
const HdrLenMin = 20
const xChBufSize = 5

var rcvMonCh chan *worker.Message
var rcvSigCh chan *worker.Message
var sndMonCh chan *worker.Message
var sndSigCh chan *worker.Message

var receiverID uint32
var senderID uint32

type Hdr struct {
	Src      uint16
	Dst      uint16
	Seq      uint32
	Ack      uint32
	Offset   uint8 // offset (4bits), reserved (3bits), ns (1bit)
	Flag     uint8 // cwr, ece, urg, ack, psh, rst, syn, fin
	Wnd      uint16
	Checksum uint16
	UrgPtr   uint16
}

// Pseudo-Headers: A Source of Controversy
// https://www.edn.com/pseudo-headers-a-source-of-controversy

type PseudoHdr struct {
	Src   mw.V4Addr
	Dst   mw.V4Addr
	Zero  uint8
	Proto uint8
	Len   uint16
}

func Open() (int, psErr.E) {
	pcb, idx := PcbRepo.UnusedPcb()
	if pcb == nil {
		psLog.E("all pcb is in used")
		return idx, psErr.CantAllocatePcb
	}
	return idx, psErr.OK
}

func Bind(id int, local EndPoint) psErr.E {
	if PcbRepo.Have(&local) {
		psLog.E(fmt.Sprintf("already bound: addr = %s, port = %d", local.Addr, local.Port))
		return psErr.AlreadyBound
	}

	pcb := PcbRepo.Get(id)
	if pcb == nil {
		psLog.E("pcb not found")
		return psErr.PcbNotFound
	}

	pcb.Local = local

	psLog.D(fmt.Sprintf("address and port assigned: addr = %s, port = %d", local.Addr, local.Port))

	return psErr.OK
}

func Listen(id int, backlog int) psErr.E {
	pcb := PcbRepo.Get(id)
	if pcb == nil {
		psLog.E("pcb not found")
		return psErr.PcbNotFound
	}

	pcb.State = listenState
	pcb.Backlog = backlog

	return psErr.OK
}

func Receive(msg *mw.TcpRxMessage) psErr.E {
	if len(msg.Segment) < HdrLenMin {
		return psErr.InvalidPacket
	}

	pseudoHdr := PseudoHdr{
		Src:   msg.Src,
		Dst:   msg.Dst,
		Proto: msg.ProtoNum,
		Len:   uint16(len(msg.Segment)),
	}
	pseudoHdrBuf := new(bytes.Buffer)
	if err := binary.Write(pseudoHdrBuf, binary.BigEndian, &pseudoHdr); err != nil {
		return psErr.Error
	}
	if mw.Checksum(msg.Segment, uint32(^mw.Checksum(pseudoHdrBuf.Bytes(), 0))) != 0 {
		psLog.E("checksum mismatch")
		return psErr.ChecksumMismatch
	}

	// TODO: handle broadcast
	if mw.V4Broadcast.EqualV4(msg.Dst) || msg.Iface.Broadcast.EqualV4(msg.Dst) {
		psLog.W("can't address broadcast (not supported)")
		return psErr.OK
	}

	hdr := Hdr{}
	if err := binary.Read(bytes.NewBuffer(msg.Segment), binary.BigEndian, &hdr); err != nil {
		return psErr.Error
	}

	offset := 4 * ((hdr.Flag & 0xf0) >> 4)
	psLog.D("", dump(&hdr, msg.Segment[offset:])...)

	// TODO:
	// ...

	return psErr.OK
}

func Send() psErr.E {
	return psErr.OK
}

func Start(wg *sync.WaitGroup) psErr.E {
	wg.Add(2)
	go receiver(wg)
	go sender(wg)
	psLog.D("tcp service started")
	return psErr.OK
}

func Stop() {
	msg := &worker.Message{
		Desired: worker.Stopped,
	}
	rcvSigCh <- msg
	sndSigCh <- msg
}

func dump(hdr *Hdr, data []byte) (ret []string) {
	var flag uint16
	flag |= uint16(hdr.Offset&0x01) << 8
	flag |= uint16(hdr.Flag & 0b10000000)
	flag |= uint16(hdr.Flag & 0b01000000)
	flag |= uint16(hdr.Flag & 0b00100000)
	flag |= uint16(hdr.Flag & 0b00010000)
	flag |= uint16(hdr.Flag & 0b00001000)
	flag |= uint16(hdr.Flag & 0b00000100)
	flag |= uint16(hdr.Flag & 0b00000010)
	flag |= uint16(hdr.Flag & 0b00000001)

	ret = append(ret, fmt.Sprintf("src port: %d", hdr.Src))
	ret = append(ret, fmt.Sprintf("dst port: %d", hdr.Dst))
	ret = append(ret, fmt.Sprintf("seq:      %d", hdr.Seq))
	ret = append(ret, fmt.Sprintf("ack:      %d", hdr.Ack))
	ret = append(ret, fmt.Sprintf("offset:   %d", (hdr.Offset&0xf0)>>4))
	ret = append(ret, fmt.Sprintf("flag:       0b%09b", flag))
	ret = append(ret, fmt.Sprintf("window:   %d", hdr.Wnd))
	ret = append(ret, fmt.Sprintf("checksum: %d", hdr.Checksum))
	ret = append(ret, fmt.Sprintf("urg:      %d", hdr.UrgPtr))

	s := "data:     "
	for i, v := range data {
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
		case msg := <-mw.TcpRxCh:
			if err := Receive(msg); err != psErr.OK {
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
		case <-mw.TcpTxCh:
			if err := Send(); err != psErr.OK {
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
}
