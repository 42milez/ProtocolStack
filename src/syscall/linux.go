package syscall

import goSyscall "syscall"

type ISyscall interface {
	CloseExec() error
	EpollCreate1Exec() (int, error)
	EpollCtlExec() error
	EpollWait() (int, error)
	IoctlExec() (uintptr, uintptr, goSyscall.Errno)
	OpenExec() (int, error)
	ReadExec() (uintptr, uintptr, goSyscall.Errno)
	SocketExec() (int, error)
}

type Syscall struct {
	Close struct {
		FD int
	}
	EpollCreate1 struct {
		Flag int
	}
	EpollCtl struct {
		EPFD  int
		OP    int
		FD    int
		Event *goSyscall.EpollEvent
	}
	EpollWait struct {
		EPFD   int
		Events []goSyscall.EpollEvent
		MSEC   int
	}
	Ioctl struct {
		A1 uintptr
		A2 uintptr
		A3 uintptr
	}
	Open struct {
		Path string
		Mode int
		Perm uint32
	}
	Read struct {
		A1 uintptr
		A2 uintptr
		A3 uintptr
	}
	Socket struct {
		Domain int
		Typ    int
		Proto  int
	}
}

func (p *Syscall) CloseExec() error {
	return goSyscall.Close(p.Close.FD)
}

func (p *Syscall) EpollCreate1Exec() (int, error) {
	return goSyscall.EpollCreate1(0)
}

func (p *Syscall) EpollCtlExec() error {
	return goSyscall.EpollCtl(p.EpollCtl.EPFD, p.EpollCtl.OP, p.EpollCtl.FD, p.EpollCtl.Event)
}

func (p *Syscall) EpollWaitExec() (int, error) {
	return goSyscall.EpollWait(p.EpollWait.EPFD, p.EpollWait.Events, p.EpollWait.MSEC)
}

func (p *Syscall) IoctlExec() (uintptr, uintptr, goSyscall.Errno) {
	return goSyscall.Syscall(goSyscall.SYS_IOCTL, p.Ioctl.A1, p.Ioctl.A2, p.Ioctl.A3)
}

func (p *Syscall) OpenExec() (int, error) {
	return goSyscall.Open(p.Open.Path, p.Open.Mode, p.Open.Perm)
}

func (p *Syscall) ReadExec() (uintptr, uintptr, goSyscall.Errno) {
	return goSyscall.Syscall(goSyscall.SYS_READ, p.Read.A1, p.Read.A2, p.Read.A3)
}

func (p *Syscall) SocketExec() (int, error) {
	return goSyscall.Socket(p.Socket.Domain, p.Socket.Typ, p.Socket.Proto)
}
