package monitor

import (
	psErr "github.com/42milez/ProtocolStack/src/error"
	"github.com/42milez/ProtocolStack/src/worker"
	"sync"
	"time"
)

const (
	Green ServiceStatus = iota
	Yellow
	Red
)
const watchInterval = 100 * time.Millisecond
const xChBufSize = 5

var sigCh chan *worker.Message
var serviceRepo *serviceRepo_

type ServiceStatus int

type serviceRepo_ struct {
	id       uint32
	services []*service
	mtx      sync.Mutex
}

func (p *serviceRepo_) Add(entry *service) (id uint32) {
	defer p.mtx.Unlock()
	p.mtx.Lock()

	id = p.id
	p.id += 1

	entry.ID = id
	p.services = append(p.services, entry)

	return
}

func (p *serviceRepo_) IsAllRunning() bool {
	defer p.mtx.Unlock()
	p.mtx.Lock()

	for _, s := range p.services {
		if s.State != worker.Running {
			return false
		}
	}

	return true
}

func (p *serviceRepo_) Watch() {
	defer p.mtx.Unlock()
	p.mtx.Lock()

	for _, s := range p.services {
		select {
		case msg := <-s.MonitorCh:
			if msg.Current == worker.Running && s.State != worker.Running {
				s.State = worker.Running
			}
		default:
			continue
		}
	}
}

type service struct {
	ID        uint32
	Name      string
	State     worker.State
	MonitorCh <-chan *worker.Message
	SignalCh  chan<- *worker.Message
}

func Register(name string, monCh <-chan *worker.Message, sigCh chan<- *worker.Message) uint32 {
	return serviceRepo.Add(&service{
		Name:      name,
		MonitorCh: monCh,
		SignalCh:  sigCh,
	})
}

func Status() ServiceStatus {
	if !serviceRepo.IsAllRunning() {
		return Red
	}
	return Green
}

func StartService(wg *sync.WaitGroup) psErr.E {
	wg.Add(1)
	go watcher(wg)
	return psErr.OK
}

func StopService() {
	sigCh <- &worker.Message{
		Desired: worker.Stopped,
	}
}

func watcher(wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case msg := <- sigCh:
			if msg.Desired == worker.Stopped {
				return
			}
		default:
			serviceRepo.Watch()
		}
		time.Sleep(watchInterval)
	}
}

func init() {
	sigCh = make(chan *worker.Message, xChBufSize)
	serviceRepo = &serviceRepo_{}
}
