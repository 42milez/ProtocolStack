package syscall

import goSyscall "syscall"

type CloseSyscallIF interface {
	Exec() error
}

type CloseSyscall struct {
	FD int
}

func (cs CloseSyscall) Exec() error {
	return goSyscall.Close(cs.FD)
}
