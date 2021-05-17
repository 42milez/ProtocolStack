package syscall

import goSyscall "syscall"

type EpollCtlSyscallInterface interface {
	Exec() error
}

type EpollCtlSyscall struct {
	EPFD  int
	OP    int
	FD    int
	Event *goSyscall.EpollEvent
}

func (p *EpollCtlSyscall) Exec() error {
	return goSyscall.EpollCtl(p.EPFD, p.OP, p.FD, p.Event)
}
