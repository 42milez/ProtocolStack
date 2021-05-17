package linux

import (
	goSyscall "syscall"
)

type IoctlSyscallInterface interface {
	Exec() (uintptr, uintptr, goSyscall.Errno)
}

type IoctlSyscall struct {
	A1 uintptr
	A2 uintptr
	A3 uintptr
}

func (i *IoctlSyscall) Exec() (uintptr, uintptr, goSyscall.Errno) {
	return goSyscall.Syscall(goSyscall.SYS_IOCTL, i.A1, i.A2, i.A3)
}
