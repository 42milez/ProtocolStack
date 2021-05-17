package syscall

import (
	"syscall"
)

type ReadSyscallInterface interface {
	Exec() (uintptr, uintptr, syscall.Errno)
}

type ReadSyscall struct {
	A1 uintptr
	A2 uintptr
	A3 uintptr
}

func (p *ReadSyscall) Exec() (uintptr, uintptr, syscall.Errno) {
	return syscall.Syscall(syscall.SYS_READ, p.A1, p.A2, p.A3)
}
