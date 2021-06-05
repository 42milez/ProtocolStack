package eth

import (
	"fmt"
	psErr "github.com/42milez/ProtocolStack/src/error"
	psLog "github.com/42milez/ProtocolStack/src/log"
	"github.com/42milez/ProtocolStack/src/mw"
	"github.com/42milez/ProtocolStack/src/worker"
	"sync"
)

const xChBufSize = 5

var RcvRxCh chan *worker.Message
var RcvTxCh chan *worker.Message
var SndRxCh chan *worker.Message
var SndTxCh chan *worker.Message

func StartService(wg *sync.WaitGroup) psErr.E {
	wg.Add(2)
	go receiver(wg)
	go sender(wg)
	return psErr.OK
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
				return
			}
		case msg := <-mw.EthRxCh:
			switch msg.Type {
			case mw.ARP:
				mw.ArpRxCh <- &mw.ArpRxMessage{
					Packet: msg.Content,
					Dev:    msg.Dev,
				}
			case mw.IPv4:
				mw.IpRxCh <- msg
			default:
				psLog.W(fmt.Sprintf("Unknown ether type: 0x%04x", uint16(msg.Type)))
			}
		case <-mw.EthTxCh:
			// TODO:
			psLog.E("not implemented")
		}
	}
}

func sender(wg *sync.WaitGroup) {
	defer wg.Done()

	SndTxCh <- &worker.Message{
		Current: worker.Running,
	}

	for {
		msg := <-SndRxCh
		if msg.Desired == worker.Stopped {
			return
		}
	}
}

func init() {
	RcvRxCh = make(chan *worker.Message, xChBufSize)
	RcvTxCh = make(chan *worker.Message, xChBufSize)
	SndRxCh = make(chan *worker.Message, xChBufSize)
	SndTxCh = make(chan *worker.Message, xChBufSize)
}
