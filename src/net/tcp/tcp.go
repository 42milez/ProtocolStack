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

const (
	//finFlag = 0x01
	synFlag = 0x02
	rstFlag = 0x04
	//pshFlag = 0x08
	ackFlag = 0x10
	//urgFlag = 0x20
)
const (
	HdrLenMax = 60 // byte
	HdrLenMin = 20
)
const xChBufSize = 5

var rcvMonCh chan *worker.Message
var rcvSigCh chan *worker.Message
var sndMonCh chan *worker.Message
var sndSigCh chan *worker.Message

var receiverID uint32
var senderID uint32

type Flag uint8

func (v Flag) IsSet(flag uint8) bool {
	return uint8(v)&flag == flag
}

type Hdr struct {
	Src      uint16 // source port
	Dst      uint16 // destination port
	Seq      uint32 // sequence number
	Ack      uint32 // acknowledgement number
	Offset   uint8  // offset (4bits), reserved (3bits), ns (1bit)
	Flag     Flag   // cwr, ece, urg, ack, psh, rst, syn, fin
	Wnd      uint16 // window size
	Checksum uint16
	Urg      uint16 // urgent pointer
}

// Pseudo-Headers: A Source of Controversy
// https://www.edn.com/pseudo-headers-a-source-of-controversy

type PseudoHdr struct {
	Src   mw.V4Addr // source address
	Dst   mw.V4Addr // destination address
	Zero  uint8     // zeros
	Proto uint8     // protocol
	Len   uint16    // segment length
}

type Segment struct {
	Seq  uint32 // sequence number
	Ack  uint32 // acknowledgement number
	Wnd  uint16 // window size
	Urg  uint16 // urgent pointer
	Data []byte // data
}

func Accept(id int, foreign *EndPoint) psErr.E {
	pcb := PcbRepo.Get(id)
	if pcb == nil {
		psLog.E("pcb not found")
		return psErr.PcbNotFound
	}

	if pcb.State != listenState {
		psLog.E("pcb is NOT in LISTEN state")
		return psErr.InvalidPcbState
	}

	// TODO: return id of pcb which accepted incoming connection
	// ...

	return psErr.OK
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

func Listen(id int, backlogSize int) psErr.E {
	pcb := PcbRepo.Get(id)
	if pcb == nil {
		psLog.E("pcb not found")
		return psErr.PcbNotFound
	}

	pcb.State = listenState
	pcb.Backlog = make([]*BacklogEntry, backlogSize)

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

	hdr := &Hdr{}
	if err := binary.Read(bytes.NewBuffer(msg.Segment), binary.BigEndian, &hdr); err != nil {
		return psErr.Error
	}

	offset := ((hdr.Flag & 0xf0) >> 4) << 2
	psLog.D("", dump(hdr, msg.Segment[offset:])...)

	local := &EndPoint{Addr: msg.Dst, Port: hdr.Dst}
	foreign := &EndPoint{Addr: msg.Src, Port: hdr.Src}

	if err := incomingSegment(hdr, msg.Segment[offset:], local, foreign); err != psErr.OK {
		psLog.E(fmt.Sprintf("can't process incoming segment: %s", err))
		return psErr.Error
	}

	return psErr.OK
}

func Send(pcb *PCB, flag uint8, data []byte) psErr.E {
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

	ret = make([]string, 10)
	ret = append(ret, fmt.Sprintf("src port: %d", hdr.Src))
	ret = append(ret, fmt.Sprintf("dst port: %d", hdr.Dst))
	ret = append(ret, fmt.Sprintf("seq:      %d", hdr.Seq))
	ret = append(ret, fmt.Sprintf("ack:      %d", hdr.Ack))
	ret = append(ret, fmt.Sprintf("offset:   %d", (hdr.Offset&0xf0)>>4))
	ret = append(ret, fmt.Sprintf("flag:       0b%09b", flag))
	ret = append(ret, fmt.Sprintf("window:   %d", hdr.Wnd))
	ret = append(ret, fmt.Sprintf("checksum: %d", hdr.Checksum))
	ret = append(ret, fmt.Sprintf("urg:      %d", hdr.Urg))

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

// SEGMENT ARRIVES
// https://datatracker.ietf.org/doc/html/rfc793#page-65

func incomingSegment(hdr *Hdr, data []byte, local *EndPoint, foreign *EndPoint) psErr.E {
	pcb := PcbRepo.LookUp(local, foreign)

	if pcb == nil {
		return psErr.Error
	}

	/* in CLOSED state */

	if pcb.State == closedState {
		// An incoming segment containing a RST is discarded.
		if hdr.Flag.IsSet(rstFlag) {
			return psErr.OK
		}
		// An incoming segment not containing a RST causes a RST to be sent in response.
		//if hdr.Flag.IsSet(ackFlag) {
		//	// TODO: send RST
		//	// ...
		//} else {
		//	// TODO: send RST,ACK
		//	// ...
		//}
	}

	/* in LISTEN state */

	isAcceptable := false

	switch pcb.State {
	case listenState:
		// 1st check for an RST
		if hdr.Flag.IsSet(rstFlag) {
			return psErr.OK
		}

		// 2nd check for an ACK
		if hdr.Flag.IsSet(ackFlag) {
			// TODO: send RST
			// ...
			return psErr.OK
		}

		// 3rd check for a SYN
		if hdr.Flag.IsSet(synFlag) {
			// TODO: check security and compartment
			// TODO: check precedence
			// ...

			newPcb, _ := PcbRepo.UnusedPcb()
			if newPcb == nil {
				return psErr.CantAllocatePcb
			}

			newPcb.Parent = pcb
			newPcb.Local = *local
			newPcb.Foreign = *foreign
			newPcb.RCV.WND = windowSize

			// Referenced from RFC793:
			// Set RCV.NXT to SEG.SEQ+1, IRS is set to SEG.SEQ and any other control or text should be queued for
			// processing later. ISS should be selected and a SYN segment sent of the form:
			// <SEQ=ISS><ACK=RCV.NXT><CTL=SYN,ACK>
			newPcb.RCV.NXT = hdr.Seq + 1
			newPcb.IRS = hdr.Seq
			newPcb.ISS = mw.RandU32()

			if err := Send(newPcb, synFlag|ackFlag, nil); err != psErr.OK {
				return psErr.Error
			}

			// Referenced from RFC793:
			// SND.NXT is set to ISS+1 and SND.UNA to ISS. The connection state should be changed to SYN-RECEIVED.
			// Note that any other incoming control or data (combined with SYN) will be processed in the SYN-RECEIVED
			// state, but processing of SYN and ACK should not be repeated. If the listen was not fully specified (i.e.,
			// the foreign socket was not fully specified), then the unspecified fields should be filled in now.
			newPcb.SND.NXT = newPcb.ISS + 1
			newPcb.SND.UNA = newPcb.ISS
			newPcb.State = synReceivedState

			return psErr.OK
		}

		// 4th, other text or control
		//
		// Referenced from RFC793:
		// Any other control or text-bearing segment (not containing SYN) must have an ACK and thus would be discarded
		// by the ACK processing. An incoming RST segment could not be valid, since it could not have been sent in
		// response to anything sent by this incarnation of the connection. So you are unlikely to get here, but if you
		// do, drop the segment, and return.
		return psErr.OK
	case synSentState:
		// 1st check the ACK bit
		if hdr.Flag.IsSet(ackFlag) {
			if hdr.Ack <= pcb.ISS || hdr.Ack > pcb.SND.NXT {
				// TODO: send RST
				// ...
				return psErr.OK
			}
			if pcb.SND.UNA <= hdr.Ack && hdr.Ack <= pcb.SND.NXT {
				isAcceptable = true
			}
		}

		// 2nd check the RST bit
		if hdr.Flag.IsSet(rstFlag) {
			if isAcceptable {
				psLog.E("connection reset")
				releasePCB(pcb)
			}
			return psErr.OK
		}

		// TODO: 3rd check the security and precedence
		// ...

		// 4th check the SYN bit
		if hdr.Flag.IsSet(synFlag) {
			pcb.RCV.NXT = hdr.Seq + 1
			pcb.IRS = hdr.Seq
			if isAcceptable {
				pcb.SND.UNA = hdr.Ack
				pcb.RefreshSndBuf()
			}
			if pcb.SND.UNA > pcb.ISS {
				pcb.State = establishedState
			}
		}

		// 5th, if neither of the SYN or RST bits is set then drop the segment and return
	}

	return psErr.OK
}

//func outgoingSegment(segment *Segment) psErr.E {
//	return psErr.OK
//}

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
			//if err := Send(); err != psErr.OK {
			//	return
			//}
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
