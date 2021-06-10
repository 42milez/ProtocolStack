package tcp

import (
	"github.com/42milez/ProtocolStack/src/mw"
	"sync"
)

const (
	freeState int = iota
	//closedState
	listenState
	//synSentState
	//synReceivedState
	//establishedState
	//finWait1State
	//finWait2State
	//closingState
	//timeWaitState
	//closeWaitState
	//lastAckState
)
const tcpConnMax = 32

var PcbRepo *pcbRepo

type EndPoint struct {
	Addr mw.V4Addr
	Port uint16
}

type PCB struct {
	State   int
	Local   EndPoint
	Foreign EndPoint
	Backlog int
	MTU uint16
	MSS uint16
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
	IRS uint32
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

	isSameAddr := func(pcb *PCB, ep2 *EndPoint) bool {
		return mw.V4Any.EqualV4(pcb.Local.Addr) || pcb.Local.Addr == ep2.Addr
	}

	isSamePort := func(pcb *PCB, ep2 *EndPoint) bool {
		return pcb.Local.Port == ep2.Port
	}

	for _, pcb := range p.pcbs {
		if isSameAddr(pcb, local) && isSamePort(pcb, local) {
			return true
		}
	}

	return false
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

func (p *pcbRepo) init() {
	for i := range p.pcbs {
		p.pcbs[i] = &PCB{}
	}
}

func init() {
	PcbRepo = &pcbRepo{}
	PcbRepo.init()
}
