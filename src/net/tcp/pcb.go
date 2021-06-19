package tcp

import (
	"container/list"
	psErr "github.com/42milez/ProtocolStack/src/error"
	"github.com/42milez/ProtocolStack/src/mw"
	"reflect"
	"sync"
	"time"
)

const (
	freeState int = iota
	closedState
	listenState
	synSentState
	synReceivedState
	establishedState
	finWait1State
	finWait2State
	closingState
	timeWaitState
	closeWaitState
	lastAckState
)
const tcpConnMax = 32
const windowSize = 65535

var PcbRepo *pcbRepo

type Backlog struct {
	entries list.List
	size    int
}

func (p *Backlog) Pop() *PCB {
	entry := p.entries.Front()
	if entry == nil {
		return nil
	}
	p.entries.Remove(entry)
	return entry.Value.(*PCB)
}

func (p *Backlog) Push(pcb *PCB) psErr.E {
	if p.entries.Len() > p.size {
		return psErr.BacklogFull
	}
	p.entries.PushBack(pcb)
	return psErr.OK
}

type EndPoint struct {
	Addr mw.V4Addr // ipv4 address
	Port uint16    // port number
}

type PCB struct {
	ID      int
	State   int      // tcp state
	Local   EndPoint // local address/port
	Foreign EndPoint // foreign address/port
	MTU     uint16   // maximum transmission unit
	MSS     uint16   // maximum segment size
	SND     struct {
		UNA uint32 // oldest unacknowledged sequence number
		NXT uint32 // next sequence number to be sent
		WND uint16 // window size
		UP  uint16 // urgent pointer to be sent
		WL1 uint32 // segment sequence number at last window update
		WL2 uint32 // segment acknowledgment number at last window update
	}
	RCV struct {
		NXT uint32 // next sequence number to receive
		WND uint16 // window size
		UP  uint16 // urgent pointer to receive
	}
	ISS         uint32
	IRS         uint32 // initial receive sequence number
	backlog     Backlog
	parent      *PCB             // parent pcb
	rcvBuf      [windowSize]byte // receive buffer
	resendQueue resendQueue      // resend buffer
}

func (p *PCB) refreshResendQueue() {
	for entry := p.resendQueue.entries.Front(); entry != nil; entry.Next() {
		if entry.Value.(*resendQueueEntry).seq >= p.SND.UNA {
			break
		}
		p.resendQueue.entries.Remove(entry)
	}
}

func (p *PCB) setTimeWaitTimer() {}

type resendQueue struct {
	entries list.List
}

func (p *resendQueue) Push() {}

type resendQueueEntry struct {
	first time.Time
	last  time.Time
	rto   uint32
	seq   uint32
	flag  uint8
	data  []byte
}

type pcbRepo struct {
	mtx  sync.Mutex
	pcbs [tcpConnMax]*PCB
}

func (p *pcbRepo) Get(id int) *PCB {
	defer p.mtx.Unlock()
	p.mtx.Lock()
	for _, v := range p.pcbs {
		if v.ID == id {
			return v
		}
	}
	return nil
}

func (p *pcbRepo) Have(local *EndPoint) bool {
	defer p.mtx.Unlock()
	p.mtx.Lock()

	for _, pcb := range p.pcbs {
		if isSameLocalEndpoint(pcb, local) {
			return true
		}
	}

	return false
}

func (p *pcbRepo) LookUp(local *EndPoint, foreign *EndPoint) *PCB {
	if local == nil || foreign == nil {
		return nil
	}

	defer p.mtx.Unlock()
	p.mtx.Lock()

	var ret *PCB
	for _, pcb := range p.pcbs {
		if isSameLocalEndpoint(pcb, local) {
			if isSameForeignEndpoint(pcb, foreign) {
				return pcb
			}
			if isListen(pcb) {
				ret = pcb
			}
		}
	}
	return ret
}

func (p *pcbRepo) PickNewPcb() (*PCB, int) {
	for i, pcb := range p.pcbs {
		if newPcb := pcb.backlog.Pop(); newPcb != nil {
			return newPcb, i
		}
	}
	return nil, -1
}

func (p *pcbRepo) UnusedPcb() *PCB {
	defer p.mtx.Unlock()
	p.mtx.Lock()
	for _, v := range p.pcbs {
		if v.State == freeState {
			return v
		}
	}
	return nil
}

func (p *pcbRepo) init() {
	defer p.mtx.Unlock()
	p.mtx.Lock()
	for i := range p.pcbs {
		p.pcbs[i] = &PCB{
			ID: i,
		}
	}
}

func isSameLocalEndpoint(pcb *PCB, ep *EndPoint) bool {
	return (mw.V4Any.EqualV4(pcb.Local.Addr) || pcb.Local.Addr == ep.Addr) && pcb.Local.Port == ep.Port
}

func isSameForeignEndpoint(pcb *PCB, ep *EndPoint) bool {
	return pcb.Foreign.Addr == ep.Addr && pcb.Foreign.Port == ep.Port
}

func isListen(pcb *PCB) bool {
	if pcb.State == listenState {
		if mw.V4Any.EqualV4(pcb.Foreign.Addr) || pcb.Foreign.Port == 0 {
			return true
		}
	}
	return false
}

func releasePCB(v *PCB) {
	p := reflect.ValueOf(v).Elem()
	p.Set(reflect.Zero(p.Type()))
}

func init() {
	PcbRepo = &pcbRepo{}
	PcbRepo.init()
}
