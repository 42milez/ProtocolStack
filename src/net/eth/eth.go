package eth

import (
	"fmt"
	psErr "github.com/42milez/ProtocolStack/src/error"
	psLog "github.com/42milez/ProtocolStack/src/log"
	"github.com/42milez/ProtocolStack/src/monitor"
	"github.com/42milez/ProtocolStack/src/mw"
	"github.com/42milez/ProtocolStack/src/worker"
	"sync"
)

var rcvMonCh chan *worker.Message
var rcvSigCh chan *worker.Message
var sndMonCh chan *worker.Message
var sndSigCh chan *worker.Message

var receiverID uint32
var senderID uint32
var xChBufSize = 5

func Start(wg *sync.WaitGroup) psErr.E {
	wg.Add(2)
	go receiver(wg)
	go sender(wg)
	psLog.D("ethernet service started")
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
				rcvMonCh <- &worker.Message{
					ID:      receiverID,
					Current: worker.Stopped,
				}
				return
			}
		case msg := <-mw.EthRxCh:
			switch msg.Type {
			case mw.EtARP:
				mw.ArpRxCh <- &mw.ArpRxMessage{
					Packet: msg.Content,
					Dev:    msg.Dev,
				}
			case mw.EtIPV4:
				mw.IpRxCh <- msg
			default:
				psLog.W(fmt.Sprintf("unknown ether type: 0x%04x", uint16(msg.Type)))
			}
		case <-mw.EthTxCh:
			// TODO:
			psLog.E("not implemented")
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
		msg := <-sndSigCh
		if msg.Desired == worker.Stopped {
			sndMonCh <- &worker.Message{
				ID:      senderID,
				Current: worker.Stopped,
			}
			return
		}
	}
}

func init() {
	rcvMonCh = make(chan *worker.Message, xChBufSize)
	rcvSigCh = make(chan *worker.Message, xChBufSize)
	receiverID = monitor.Register("Eth Receiver", rcvMonCh, rcvSigCh)

	sndMonCh = make(chan *worker.Message, xChBufSize)
	sndSigCh = make(chan *worker.Message, xChBufSize)
	senderID = monitor.Register("Eth Sender", sndMonCh, sndSigCh)
}
