package linux

import goSyscall "syscall"

type SocketSyscallInterface interface {
	Exec() (int, error)
}

type SocketSyscall struct {
	Domain int
	Typ    int
	Proto  int
}

func (p *SocketSyscall) Exec() (int, error) {
	return goSyscall.Socket(p.Domain, p.Typ, p.Proto)
}
