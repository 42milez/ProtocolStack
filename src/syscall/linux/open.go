package linux

import goSyscall "syscall"

type OpenSyscallInterface interface {
	Exec() (int, error)
}

type OpenSyscall struct {
	Path string
	Mode int
	Perm uint32
}

func (p *OpenSyscall) Exec() (int, error) {
	return goSyscall.Open(p.Path, p.Mode, p.Perm)
}
