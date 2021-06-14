package tcp

import (
	"github.com/42milez/ProtocolStack/src/mw"
	"sync"
)

const (
	freeState int = iota
	closedState
	listenState
	synSentState
	synReceivedState
	//establishedState
	//finWait1State
	//finWait2State
	//closingState
	//timeWaitState
	//closeWaitState
	//lastAckState
)
const bufSize = 65535
const tcpConnMax = 32

var PcbRepo *pcbRepo

type BacklogEntry struct{}

type EndPoint struct {
	Addr mw.V4Addr
	Port uint16
}

type PCB struct {
	State   int
	Local   EndPoint
	Foreign EndPoint
	MTU     uint16
	MSS     uint16
	SND     struct {
		UNA uint32
		NXT uint32
		WND uint16
		UP  uint16
		WL1 uint32
		WL2 uint32
	}
	ISS uint32
	RCV struct {
		NXT uint32
		WND uint16
		UP  uint16
	}
	IRS     uint32
	Buf     [bufSize]byte
	Backlog []*BacklogEntry
	Child   []*PCB
	Parent  *PCB
}

func (p *PCB) BelongsTo(pcb *PCB) {
	p.Parent = pcb
	pcb.Child = append(pcb.Child, p)
}

type pcbRepo struct {
	mtx  sync.Mutex
	pcbs [tcpConnMax]*PCB
}

func (p *pcbRepo) Get(idx int) *PCB {
	defer p.mtx.Unlock()
	p.mtx.Lock()
	return p.pcbs[idx]
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

func (p *pcbRepo) UnusedPcb() (*PCB, int) {
	defer p.mtx.Unlock()
	p.mtx.Lock()
	for i := 0; i < tcpConnMax; i++ {
		if p.pcbs[i].State == freeState {
			return p.pcbs[i], i
		}
	}
	return nil, -1
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

func (p *pcbRepo) init() {
	defer p.mtx.Unlock()
	p.mtx.Lock()
	for i := range p.pcbs {
		p.pcbs[i] = &PCB{}
	}
}

func init() {
	PcbRepo = &pcbRepo{}
	PcbRepo.init()
}
