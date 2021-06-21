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
	"time"
)

const (
	finFlag = 0x01
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

type SegmentInfo struct {
	Seq  uint32 // sequence number
	Ack  uint32 // acknowledgement number
	Wnd  uint16 // window size
	Flag Flag
}

func Open() (int, psErr.E) {
	pcb := PcbRepo.UnusedPcb()
	if pcb == nil {
		psLog.E("all pcb is in used")
		return -1, psErr.CantAllocatePcb
	}
	return pcb.ID, psErr.OK
}

func Listen(id int, backlogSize int) psErr.E {
	pcb := PcbRepo.Get(id)
	if pcb == nil {
		psLog.E("pcb not found")
		return psErr.PcbNotFound
	}

	pcb.State = listenState
	pcb.backlog.size = backlogSize

	return psErr.OK
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

func Accept(id int) (int, EndPoint, psErr.E) {
	var foreign EndPoint

	pcb := PcbRepo.Get(id)
	if pcb == nil {
		psLog.E("pcb not found")
		return -1, foreign, psErr.PcbNotFound
	}

	if pcb.State != listenState {
		psLog.E("pcb is NOT in LISTEN state")
		return -1, foreign, psErr.InvalidPcbState
	}

	var newPcb *PCB
	var pcbId int
	for {
		if newPcb, pcbId = PcbRepo.PickNewConnection(); newPcb != nil {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	foreign = newPcb.Foreign

	return pcbId, foreign, psErr.OK
}

func Connect(id int, foreign EndPoint) psErr.E {
	return psErr.OK
}

func Receive(msg *mw.TcpRxMessage) psErr.E {
	if len(msg.RawSegment) < HdrLenMin {
		return psErr.InvalidPacket
	}

	pseudoHdr := PseudoHdr{
		Src:   msg.Src,
		Dst:   msg.Dst,
		Proto: msg.ProtoNum,
		Len:   uint16(len(msg.RawSegment)),
	}
	pseudoHdrBuf := new(bytes.Buffer)
	if err := binary.Write(pseudoHdrBuf, binary.BigEndian, &pseudoHdr); err != nil {
		return psErr.Error
	}
	if mw.Checksum(msg.RawSegment, uint32(^mw.Checksum(pseudoHdrBuf.Bytes(), 0))) != 0 {
		psLog.E("checksum mismatch (tcp)")
		return psErr.ChecksumMismatch
	}

	// TODO: handle broadcast
	if mw.V4Broadcast.EqualV4(msg.Dst) || msg.Iface.Broadcast.EqualV4(msg.Dst) {
		psLog.W("can't address broadcast (not supported)")
		return psErr.OK
	}

	hdr := &Hdr{}
	if err := binary.Read(bytes.NewBuffer(msg.RawSegment), binary.BigEndian, hdr); err != nil {
		return psErr.Error
	}

	hdrLen := int((hdr.Offset&0xf0)>>4) << 2
	psLog.D("incoming tcp segment", dump(msg.RawSegment, hdrLen)...)

	local := &EndPoint{Addr: msg.Dst, Port: hdr.Dst}
	foreign := &EndPoint{Addr: msg.Src, Port: hdr.Src}

	if err := receiveCore(hdr, msg.RawSegment[hdrLen:], local, foreign); err != psErr.OK {
		psLog.E(fmt.Sprintf("can't process incoming segment: %s", err))
		return psErr.Error
	}

	return psErr.OK
}

// SEGMENT ARRIVES
// https://datatracker.ietf.org/doc/html/rfc793#page-65

func receiveCore(hdr *Hdr, data []byte, local *EndPoint, foreign *EndPoint) psErr.E {
	pcb := PcbRepo.LookUp(local, foreign)

	if pcb == nil {
		return psErr.Error
	}

	// ==================================================
	//  CLOSED state
	// ==================================================

	// If the state is CLOSED (i.e., TCB does not exist) then all Data in the incoming segment is discarded. An incoming
	// segment containing a RST is discarded. An incoming segment not containing a RST causes a RST to be sent in
	// response. The acknowledgment and sequence field values are selected to make the reset sequence acceptable to the
	// TCP that sent the offending segment.
	if pcb.State == closedState {
		if hdr.Flag.IsSet(rstFlag) {
			return psErr.OK
		}

		// If the ACK bit is on, <SEQ=SEG.ACK><CTL=RST>
		if hdr.Flag.IsSet(ackFlag) {
			info := SegmentInfo{
				Seq:  hdr.Ack,
				Ack:  0,
				Wnd:  0,
				Flag: rstFlag,
			}
			if err := sendCore(info, nil, local, foreign); err != psErr.OK {
				return psErr.Error
			}
		} else {
			// If the ACK bit is off, sequence number zero is used, <SEQ=0><ACK=SEG.SEQ+SEG.LEN><CTL=RST,ACK>
			info := SegmentInfo{
				Seq:  0,
				Ack:  hdr.Seq + uint32(len(data)),
				Wnd:  0,
				Flag: ackFlag | rstFlag,
			}
			if err := sendCore(info, nil, local, foreign); err != psErr.OK {
				return psErr.Error
			}
		}

		return psErr.OK
	}

	// ==================================================
	//  LISTEN / SYN-SENT state
	// ==================================================

	isAcceptable := false

	switch pcb.State {
	case listenState:
		//  ▶ 1st check for an RST
		// ------------------------------
		if hdr.Flag.IsSet(rstFlag) {
			return psErr.OK
		}

		//  ▶ 2nd check for an ACK
		// ------------------------------
		if hdr.Flag.IsSet(ackFlag) {
			// Any acknowledgment is bad if it arrives on a connection still in the LISTEN state. An acceptable reset
			// segment should be formed for any arriving ACK-bearing segment. The RST should be formatted as follows:
			// <SEQ=SEG.ACK><CTL=RST> Return.
			info := SegmentInfo{
				Seq:  hdr.Ack,
				Ack:  0,
				Wnd:  0,
				Flag: rstFlag,
			}
			if err := sendCore(info, nil, local, foreign); err != psErr.OK {
				return psErr.Error
			}
			return psErr.OK
		}

		//  ▶ 3rd check for a SYN
		// ------------------------------
		if hdr.Flag.IsSet(synFlag) {
			// TODO: check security and compartment
			// TODO: check precedence
			// ...

			newPcb := PcbRepo.UnusedPcb()
			if newPcb == nil {
				return psErr.CantAllocatePcb
			}

			newPcb.parent = pcb
			newPcb.Local = *local
			newPcb.Foreign = *foreign
			newPcb.RCV.WND = windowSize

			// Set RCV.NXT to SEG.SEQ+1, IRS is set to SEG.SEQ and any other control or text should be queued for
			// processing later. ISS should be selected and a SYN segment sent of the form:
			// <SEQ=ISS><ACK=RCV.NXT><CTL=SYN,ACK>
			newPcb.RCV.NXT = hdr.Seq + 1
			newPcb.IRS = hdr.Seq
			newPcb.ISS = mw.RandU32()

			if err := Send(newPcb, synFlag|ackFlag, nil); err != psErr.OK {
				return psErr.Error
			}

			// SND.NXT is set to ISS+1 and SND.UNA to ISS. The connection state should be changed to SYN-RECEIVED.
			// Note that any other incoming control or Data (combined with SYN) will be processed in the SYN-RECEIVED
			// state, but processing of SYN and ACK should not be repeated. If the listen was not fully specified (i.e.,
			// the foreign socket was not fully specified), then the unspecified fields should be filled in now.
			newPcb.SND.NXT = newPcb.ISS + 1
			newPcb.SND.UNA = newPcb.ISS
			newPcb.State = synReceivedState

			return psErr.OK
		}

		//  ▶ 4th, other text or control
		// ------------------------------
		// Any other control or text-bearing segment (not containing SYN) must have an ACK and thus would be discarded
		// by the ACK processing. An incoming RST segment could not be valid, since it could not have been sent in
		// response to anything sent by this incarnation of the connection. So you are unlikely to get here, but if you
		// do, drop the segment, and return.
		return psErr.OK
	case synSentState:
		//  ▶ 1st check the ACK bit
		// ------------------------------
		if hdr.Flag.IsSet(ackFlag) {
			// If the ACK bit is set: If SEG.ACK =< ISS, or SEG.ACK > SND.NXT, send a reset (unless the RST bit is set,
			//                        if so drop the segment and return) <SEQ=SEG.ACK><CTL=RST> and discard the segment.
			//                        Return.
			if hdr.Ack <= pcb.ISS || hdr.Ack > pcb.SND.NXT {
				info := SegmentInfo{
					Seq:  hdr.Ack,
					Ack:  0,
					Wnd:  0,
					Flag: rstFlag,
				}
				return sendCore(info, data, local, foreign)
			}
			// If SND.UNA =< SEG.ACK =< SND.NXT then the ACK is acceptable.
			if pcb.SND.UNA <= hdr.Ack && hdr.Ack <= pcb.SND.NXT {
				isAcceptable = true
			}
		}

		//  ▶ 2nd check the RST bit
		// ------------------------------
		if hdr.Flag.IsSet(rstFlag) {
			if isAcceptable {
				psLog.E("connection reset")
				releasePCB(pcb)
			}
			return psErr.OK
		}

		// TODO: 3rd check the security and precedence
		// ...

		//  ▶ 4th check the SYN bit
		// ------------------------------
		// This step should be reached only if the ACK is ok, or there is no ACK, and it the segment did not contain a
		// RST.
		if hdr.Flag.IsSet(synFlag) {
			// If the SYN bit is on and the security/compartment and precedence are acceptable then, RCV.NXT is set to
			// SEG.SEQ+1, IRS is set to SEG.SEQ. SND.UNA should be advanced to equal SEG.ACK (if there is an ACK), and
			// any segments on the retransmission queue which
			// are thereby acknowledged should be removed.
			pcb.RCV.NXT = hdr.Seq + 1
			pcb.IRS = hdr.Seq
			if isAcceptable {
				pcb.SND.UNA = hdr.Ack
				pcb.refreshResendQueue()
			}
			// If SND.UNA > ISS (our SYN has been ACKed), change the connection state to ESTABLISHED, form an ACK
			// segment <SEQ=SND.NXT><ACK=RCV.NXT><CTL=ACK> and send it. Data or controls which were queued for
			// transmission may be included. If there are other controls or text in the segment then continue processing
			// at the sixth step below where the URG bit is checked, otherwise return.
			if pcb.SND.UNA > pcb.ISS {
				pcb.State = establishedState
				if err := Send(pcb, ackFlag, nil); err != psErr.OK {
					return psErr.Error
				}
				pcb.SND.WND = hdr.Wnd
				pcb.SND.WL1 = hdr.Seq
				pcb.SND.WL2 = hdr.Ack
				return psErr.OK
			} else {
				// Otherwise enter SYN-RECEIVED, form a SYN,ACK segment <SEQ=ISS><ACK=RCV.NXT><CTL=SYN,ACK> and send it.
				// If there are other controls or text in the segment, queue them for processing after the ESTABLISHED
				// state has been reached, return.
				pcb.State = synReceivedState
				if err := Send(pcb, synFlag|ackFlag, nil); err != psErr.OK {
					return psErr.Error
				}
				return psErr.OK
			}
		}

		//  ▶ 5th, if neither of the SYN or RST bits is set then drop the segment and return
		// ------------------------------
		return psErr.OK
	}

	// ==================================================
	//  Otherwise
	// ==================================================

	//  ▶ 1st check sequence number
	// ------------------------------
	// Segments are processed in sequence. Initial tests on arrival are used to discard old duplicates, but further
	// processing is done in SEG.SEQ order. If a segment's contents straddle the boundary between old and new, only the
	// new parts should be processed.
	switch pcb.State {
	case synReceivedState:
		fallthrough
	case establishedState:
		fallthrough
	case finWait1State:
		fallthrough
	case finWait2State:
		fallthrough
	case closeWaitState:
		fallthrough
	case closingState:
		fallthrough
	case lastAckState:
		fallthrough
	case timeWaitState:
		if len(data) == 0 {
			if pcb.RCV.WND == 0 {
				// case 1
				if hdr.Seq == pcb.RCV.NXT {
					isAcceptable = true
				}
			} else {
				// case 2
				a1 := pcb.RCV.NXT <= hdr.Seq
				a2 := hdr.Seq < (pcb.RCV.NXT + uint32(pcb.RCV.WND))
				if a1 && a2 {
					isAcceptable = true
				}
			}
		} else {
			if pcb.RCV.WND == 0 {
				// case 3: do nothing (not acceptable)
				// If the RCV.WND is zero, no segments will be acceptable, but special allowance should be made to accept
				// valid ACKs, URGs and RSTs.
			} else {
				// case 4
				a1 := pcb.RCV.NXT <= hdr.Seq
				a2 := hdr.Seq < (pcb.RCV.NXT + uint32(pcb.RCV.WND))
				b1 := pcb.RCV.NXT <= hdr.Seq+uint32(len(data))-1
				b2 := hdr.Seq+uint32(len(data))-1 < pcb.RCV.NXT+uint32(pcb.RCV.WND)
				if a1 && a2 || b1 && b2 {
					isAcceptable = true
				}
			}
		}

		// If an incoming segment is not acceptable, an acknowledgment should be sent in reply (unless the RST bit is
		// set, if so drop the segment and return): <SEQ=SND.NXT><ACK=RCV.NXT><CTL=ACK>
		// After sending the acknowledgment, drop the unacceptable segment and return.
		if !isAcceptable {
			if hdr.Flag.IsSet(rstFlag) {
				if err := Send(pcb, ackFlag, nil); err != psErr.OK {
					return psErr.Error
				}
			}
			return psErr.OK
		}

		// In the following it is assumed that the segment is the idealized segment that begins at RCV.NXT and does not
		// exceed the window. One could tailor actual segments to fit this assumption by trimming off any portions that
		// lie outside the window (including SYN and FIN), and only processing further if the segment then begins at
		// RCV.NXT. Segments with higher beginning sequence numbers may be held for later processing.
	}

	//  ▶ 2nd check the RST bit
	// ------------------------------
	switch pcb.State {
	case synReceivedState:
		// If this connection was initiated with a passive OPEN (i.e., came from the LISTEN state), then return this
		// connection to LISTEN state and return. The user need not be informed. If this connection was initiated with
		// an active OPEN (i.e., came from SYN-SENT state) then the connection was refused, signal the user "connection
		// refused". In either case, all segments on the retransmission queue should be removed. And in the active OPEN
		// case, enter the CLOSED state and delete the TCB, and return.
		if hdr.Flag.IsSet(rstFlag) {
			pcb.State = closedState
			releasePCB(pcb)
			return psErr.OK
		}
	case establishedState:
		fallthrough
	case finWait1State:
		fallthrough
	case finWait2State:
		fallthrough
	case closeWaitState:
		// If the RST bit is set then, any outstanding RECEIVEs and SEND should receive "reset" responses. All segment
		// queues should be flushed. Users should also receive an unsolicited general "connection reset" signal. Enter
		// the CLOSED state, delete the TCB, and return.
		if hdr.Flag.IsSet(rstFlag) {
			psLog.E("connection reset")
			pcb.State = closedState
			releasePCB(pcb)
			return psErr.OK
		}
	case closingState:
		fallthrough
	case lastAckState:
		fallthrough
	case timeWaitState:
		// If the RST bit is set then, enter the CLOSED state, delete the TCB, and return.
		if hdr.Flag.IsSet(rstFlag) {
			pcb.State = closedState
			releasePCB(pcb)
			return psErr.OK
		}
	}

	//  ▶ 3rd check security and precedence
	// ------------------------------
	// TODO:
	// ....

	//  ▶ 4th check the SYN bit
	// ------------------------------
	switch pcb.State {
	case synReceivedState:
		fallthrough
	case establishedState:
		fallthrough
	case finWait1State:
		fallthrough
	case finWait2State:
		fallthrough
	case closeWaitState:
		fallthrough
	case closingState:
		fallthrough
	case lastAckState:
		fallthrough
	case timeWaitState:
		// If the SYN is in the window it is an error, send a reset, any outstanding RECEIVEs and SEND should receive
		// "reset" responses, all segment queues should be flushed, the user should also receive an unsolicited general
		// "connection reset" signal, enter the CLOSED state, delete the TCB, and return.
		// If the SYN is not in the window this step would not be reached and an ack would have been sent in the First
		// step (sequence number check).
		if hdr.Flag.IsSet(synFlag) {
			if err := Send(pcb, rstFlag, nil); err != psErr.OK {
				return psErr.Error
			}
			psLog.E("connection reset")
			pcb.State = closedState
			releasePCB(pcb)
			return psErr.OK
		}
	}

	//  ▶ 5th check the ACK field
	// ------------------------------
	// if the ACK bit is off drop the segment and return
	if !hdr.Flag.IsSet(ackFlag) {
		return psErr.OK
	}

	switch pcb.State {
	case synReceivedState:
		// If SND.UNA =< SEG.ACK =< SND.NXT then enter ESTABLISHED state and continue processing.
		if pcb.SND.UNA <= hdr.Ack && hdr.Ack <= pcb.SND.NXT {
			pcb.State = establishedState
			if pcb.parent != nil {
				if err := pcb.parent.backlog.Push(pcb); err != psErr.OK {
					psLog.E("backlog full")
					return err
				}
			}
			// If the segment acknowledgment is not acceptable, form a reset segment, <SEQ=SEG.ACK><CTL=RST> and send
			// it.
		} else {
			info := SegmentInfo{
				Seq:  hdr.Ack,
				Ack:  0,
				Wnd:  0,
				Flag: rstFlag,
			}
			if err := sendCore(info, nil, local, foreign); err != psErr.OK {
				return psErr.Error
			}
			return psErr.OK
		}
		fallthrough
	case establishedState:
		fallthrough
	case finWait1State:
		fallthrough
	case finWait2State:
		fallthrough
	case closeWaitState:
		fallthrough
	case closingState:
		if pcb.SND.UNA < hdr.Ack && hdr.Ack <= pcb.SND.NXT {
			// If SND.UNA < SEG.ACK =< SND.NXT then, set SND.UNA <- SEG.ACK. Any segments on the retransmission queue
			// which are thereby entirely acknowledged are removed. Users should receive positive acknowledgments for
			// buffers which have been SENT and fully acknowledged (i.e., SEND buffer should be returned with "ok"
			// response).
			// If SND.UNA < SEG.ACK =< SND.NXT, the send window should be updated.
			pcb.SND.UNA = hdr.Ack

			// TODO: clean up resend queue
			// ...

			// If (SND.WL1 < SEG.SEQ or (SND.WL1 = SEG.SEQ and SND.WL2 =< SEG.ACK)), set SND.WND <- SEG.WND,
			// set SND.WL1 <- SEG.SEQ, and set SND.WL2 <- SEG.ACK.
			if pcb.SND.WL1 < hdr.Seq || (pcb.SND.WL1 == hdr.Seq && pcb.SND.WL2 <= hdr.Ack) {
				pcb.SND.WND = hdr.Wnd
				pcb.SND.WL1 = hdr.Seq
				pcb.SND.WL2 = hdr.Ack
			}
		} else if hdr.Ack < pcb.SND.UNA {
			// If the ACK is a duplicate (SEG.ACK < SND.UNA), it can be ignored.
		} else if hdr.Ack > pcb.SND.NXT {
			// If the ACK acks something not yet sent (SEG.ACK > SND.NXT) then send an ACK, drop the segment, and
			// return.
			if err := Send(pcb, ackFlag, nil); err != psErr.OK {
				return psErr.Error
			}
			return psErr.OK
		}

		// Note that SND.WND is an offset from SND.UNA, that SND.WL1 records the sequence number of the Last segment
		// used to update SND.WND, and that SND.WL2 records the acknowledgment number of the Last segment used to update
		// SND.WND. The check here prevents using old segments to update the window.

		switch pcb.State {
		case finWait1State:
			// In addition to the processing for the ESTABLISHED state, if our FIN is now acknowledged then enter
			// FIN-WAIT-2 and continue processing in that state.
			if hdr.Ack == pcb.SND.NXT {
				pcb.State = finWait2State
			}
		case finWait2State:
			// In addition to the processing for the ESTABLISHED state, if the retransmission queue is empty, the user's
			// CLOSE can be acknowledged ("ok") but do not delete the TCB.
		case closeWaitState:
			// Do the same processing as for the ESTABLISHED state.
		case closingState:
			// In addition to the processing for the ESTABLISHED state, if the ACK acknowledges our FIN then enter the
			// TIME-WAIT state, otherwise ignore the segment.
			if hdr.Ack == pcb.SND.NXT {
				pcb.State = timeWaitState
				// TODO: set time wait timer
				// ...
			}
		}
	case lastAckState:
		// The only thing that can arrive in this state is an acknowledgment of our FIN. If our FIN is now acknowledged,
		// delete the TCB, enter the CLOSED state, and return.
		if hdr.Ack == pcb.SND.NXT {
			pcb.State = closedState
			releasePCB(pcb)
		}
		return psErr.OK
	case timeWaitState:
		// The only thing that can arrive in this state is a retransmission of the remote FIN. Acknowledge it, and
		// restart the 2 MSL timeout.
		if hdr.Flag.IsSet(finFlag) {
			// TODO: set time wait timer
			pcb.setTimeWaitTimer()
		}
	}

	//  ▶ 6th check the URG bit
	// ------------------------------
	// TODO:
	// ...

	//  ▶ 7th process the segment text
	// ------------------------------
	switch pcb.State {
	case establishedState:
		fallthrough
	case finWait1State:
		fallthrough
	case finWait2State:
		// Once in the ESTABLISHED state, it is possible to deliver segment text to user RECEIVE buffers. Text from
		// segments can be moved into buffers until either the buffer is full or the segment is empty. If the segment
		// empties and carries an PUSH Flag, then the user is informed, when the buffer is returned, that a PUSH has
		// been received. When the TCP takes responsibility for delivering the Data to the user it must also acknowledge
		// the receipt of the Data. Once the TCP takes responsibility for the Data it advances RCV.NXT over the Data
		// accepted, and adjusts RCV.WND as appropriate to the current buffer availability. The total of RCV.NXT and
		// RCV.WND should not be reduced. Please note the window management suggestions in section 3.7. Send an
		// acknowledgment of the form: <SEQ=SND.NXT><ACK=RCV.NXT><CTL=ACK> This acknowledgment should be piggybacked on
		// a segment being transmitted if possible without incurring undue delay.
		if len(data) != 0 {
			copy(pcb.rcvBuf[:pcb.RCV.WND], data)
			pcb.RCV.NXT = hdr.Seq + uint32(len(data))
			pcb.RCV.WND -= uint16(len(data))
			if err := Send(pcb, ackFlag, nil); err != psErr.OK {
				return psErr.Error
			}
		}
	case closeWaitState:
		fallthrough
	case closingState:
		fallthrough
	case lastAckState:
		fallthrough
	case timeWaitState:
		// This should not occur, since a FIN has been received from the remote side. Ignore the segment text.
	}

	//  ▶ 8th check the FIN bit
	// ------------------------------
	if hdr.Flag.IsSet(finFlag) {
		switch pcb.State {
		case closedState:
			fallthrough
		case listenState:
			fallthrough
		case synSentState:
			// Do not process the FIN if the state is CLOSED, LISTEN or SYN-SENT since the SEG.SEQ cannot be validated;
			// drop the segment and return.
			return psErr.OK
		}

		// If the FIN bit is set, signal the user "connection closing" and return any pending RECEIVEs with same
		// message, advance RCV.NXT over the FIN, and send an acknowledgment for the FIN. Note that FIN implies PUSH for
		// any segment text not yet delivered to the user.
		psLog.I("connection closing")

		pcb.RCV.NXT = hdr.Seq + 1
		if err := Send(pcb, ackFlag, nil); err != psErr.OK {
			return psErr.Error
		}

		switch pcb.State {
		case synReceivedState:
			fallthrough
		case establishedState:
			pcb.State = closeWaitState
		case finWait1State:
			// If our FIN has been ACKed (perhaps in this segment), then enter TIME-WAIT, start the time-wait timer,
			// turn off the other timers; otherwise enter the CLOSING state.
			if hdr.Ack == pcb.SND.NXT {
				pcb.State = timeWaitState
				// TODO: set time wait timer
				// ...
			} else {
				pcb.State = closingState
			}
		case finWait2State:
			// Enter the TIME-WAIT state. Start the time-wait timer, turn off the other timers.
			pcb.State = timeWaitState
			// TODO: set time wait timer
			// ...
		case closeWaitState:
			// Remain in the CLOSE-WAIT state.
		case closingState:
			// Remain in the CLOSING state.
		case lastAckState:
			// Remain in the LAST-ACK state.
		case timeWaitState:
			// Remain in the TIME-WAIT state. Restart the 2 MSL time-wait timeout.
			// TODO: set time wait timer
			// ...
		}
	}

	return psErr.OK
}

func Send(pcb *PCB, flag Flag, data []byte) psErr.E {
	info := SegmentInfo{
		Seq:  pcb.SND.NXT,
		Ack:  pcb.RCV.NXT,
		Wnd:  pcb.RCV.WND,
		Flag: flag,
	}
	if flag.IsSet(synFlag) {
		info.Seq = pcb.ISS
	}
	if flag.IsSet(synFlag|finFlag) || len(data) != 0 {
		// TODO: add to retransmit queue
		pcb.resendQueue.Push()
	}

	return sendCore(info, data, &pcb.Local, &pcb.Foreign)
}

func sendCore(info SegmentInfo, data []byte, local *EndPoint, foreign *EndPoint) psErr.E {
	hdr := Hdr{
		Src:    local.Port,
		Dst:    foreign.Port,
		Seq:    info.Seq,
		Ack:    info.Ack,
		Offset: uint8(HdrLenMin>>2) << 4, // Data Offset occupies high 4bits
		Flag:   info.Flag,
		Wnd:    info.Wnd,
	}

	segBuf := new(bytes.Buffer)
	if err := binary.Write(segBuf, binary.BigEndian, &hdr); err != nil {
		return psErr.Error
	}
	if err := binary.Write(segBuf, binary.BigEndian, &data); err != nil {
		return psErr.Error
	}
	segment := segBuf.Bytes()

	pseudo := PseudoHdr{
		Src:   local.Addr,
		Dst:   foreign.Addr,
		Proto: uint8(mw.PnTCP),
		Len:   uint16(len(segment)),
	}
	pseudoBuf := new(bytes.Buffer)
	if err := binary.Write(pseudoBuf, binary.BigEndian, &pseudo); err != nil {
		return psErr.Error
	}

	csum := mw.Checksum(segment, uint32(^mw.Checksum(pseudoBuf.Bytes(), 0)))
	segment[16] = uint8((csum & 0xff00) >> 8)
	segment[17] = uint8(csum & 0x00ff)

	psLog.D("outgoing tcp segment", dump(segment, HdrLenMin)...)

	mw.IpTxCh <- &mw.IpMessage{
		ProtoNum: mw.PnTCP,
		Packet:   segment,
		Src:      local.Addr,
		Dst:      foreign.Addr,
	}

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

func dump(segment []byte, hdrLen int) (ret []string) {
	var flag uint16
	flag |= uint16(segment[12]&0x01) << 8
	flag |= uint16(segment[13] & (0x01 << 7))
	flag |= uint16(segment[13] & (0x01 << 6))
	flag |= uint16(segment[13] & (0x01 << 5))
	flag |= uint16(segment[13] & (0x01 << 4))
	flag |= uint16(segment[13] & (0x01 << 3))
	flag |= uint16(segment[13] & (0x01 << 2))
	flag |= uint16(segment[13] & (0x01 << 1))
	flag |= uint16(segment[13] & 0x01)

	ret = append(ret, fmt.Sprintf("src port: %d", uint16(segment[0])<<8|uint16(segment[1])))
	ret = append(ret, fmt.Sprintf("dst port: %d", uint16(segment[2])<<8|uint16(segment[3])))
	ret = append(ret, fmt.Sprintf("Seq:      %d",
		uint32(segment[4])<<24|
			uint32(segment[5])<<16|
			uint32(segment[6])<<8|
			uint32(segment[7])))
	ret = append(ret, fmt.Sprintf("ack:      %d",
		uint32(segment[8])<<24|
			uint32(segment[9])<<16|
			uint32(segment[10])<<8|
			uint32(segment[11])))
	ret = append(ret, fmt.Sprintf("offset:   %d", (segment[12]&0xf0)>>4))
	ret = append(ret, fmt.Sprintf("flag:     0b%09b", flag))
	ret = append(ret, fmt.Sprintf("window:   %d", uint16(segment[14])<<8|uint16(segment[15])))
	ret = append(ret, fmt.Sprintf("checksum: 0x%04x", uint16(segment[16])<<8|uint16(segment[17])))
	ret = append(ret, fmt.Sprintf("urg:      %d", uint16(segment[18])<<8|uint16(segment[19])))

	s := "data:     "
	data := segment[:hdrLen]
	for i, v := range data {
		s += fmt.Sprintf("%02x ", v)
		if (i+1)%20 == 0 && i+1 != len(data) {
			s += "\n                                  "
		}
	}
	ret = append(ret, s)

	return
}

func receiver(wg *sync.WaitGroup) {
	defer func() {
		psLog.D("tcp receiver stopped")
		wg.Done()
	}()

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
	defer func() {
		psLog.D("tcp sender stopped")
		wg.Done()
	}()

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
