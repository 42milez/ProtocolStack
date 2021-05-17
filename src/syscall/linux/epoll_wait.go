package linux

import goSyscall "syscall"

type EpollWaitSyscallInterface interface {
	Exec() (int, error)
}

type EpollWaitSyscall struct {
	EPFD   int
	Events []goSyscall.EpollEvent
	MSEC   int
}

func (p *EpollWaitSyscall) Exec() (int, error) {
	return goSyscall.EpollWait(p.EPFD, p.Events, p.MSEC)
}
