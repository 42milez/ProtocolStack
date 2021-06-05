package eth

import (
	"fmt"
	psErr "github.com/42milez/ProtocolStack/src/error"
	psLog "github.com/42milez/ProtocolStack/src/log"
	"github.com/42milez/ProtocolStack/src/mw"
	"github.com/42milez/ProtocolStack/src/worker"
	"sync"
)

const ReceiverID worker.ID = 1
const SenderID worker.ID = 2
const xChBufSize = 5

var MonitorCh chan *worker.Message
var ReceiverSigCh chan *worker.Message
var SenderSigCh chan *worker.Message

func StartService(wg *sync.WaitGroup) psErr.E {
	wg.Add(2)
	go receiver(wg)
	go sender(wg)
	return psErr.OK
}

func receiver(wg *sync.WaitGroup) {
	defer wg.Done()

	MonitorCh <- &worker.Message{
		ID:      ReceiverID,
		Current: worker.Running,
	}

	for {
		select {
		case msg := <-ReceiverSigCh:
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

	MonitorCh <- &worker.Message{
		ID:      SenderID,
		Current: worker.Running,
	}

	for {
		msg := <-SenderSigCh
		if msg.Desired == worker.Stopped {
			return
		}
	}
}

func init() {
	MonitorCh = make(chan *worker.Message, xChBufSize)
	ReceiverSigCh = make(chan *worker.Message, xChBufSize)
	SenderSigCh = make(chan *worker.Message, xChBufSize)
}
