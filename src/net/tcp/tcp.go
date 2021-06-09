package tcp

import (
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
	Offset   uint8
	Flag     uint8
	Wnd      uint16
	Checksum uint16
	UrgPtr   uint16
}

// Pseudo-Headers: A Source of Controversy
// https://www.edn.com/pseudo-headers-a-source-of-controversy

type PseudoHdr struct {
	Src   [mw.V4AddrLen]byte
	Dst   [mw.V4AddrLen]byte
	Zero  uint8
	Proto uint8
	Len   uint16
}

func Receive(msg *mw.TcpRxMessage) psErr.E {
	//if len(msg.Payload) < HdrLenMin {
	//	return psErr.InvalidPacket
	//}
	//
	//hdr := Hdr{}
	//buf := bytes.NewBuffer(msg.Payload)
	//if err := binary.Read(buf, binary.BigEndian, &hdr); err != nil {
	//	return psErr.Error
	//}
	//
	//pHdr := PseudoHdr{
	//	Src: msg.Src,
	//	Dst: msg.Dst,
	//	Proto: msg.ProtoNum,
	//	Len: uint16(len(msg.Payload)),
	//}
	//
	//
	//mw.Checksum(pHdr)

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
