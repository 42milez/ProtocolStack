package sys

import "syscall"

type SyscallInterface interface {
	Close(fd int) error
	EpollCreate1() (int, error)
	EpollCtl() error
	EpollWait() (int, error)
	Ioctl() (uintptr, uintptr, syscall.Errno)
	Open() (int, error)
	Read() (uintptr, uintptr, syscall.Errno)
	Socket() (int, error)
}

type Syscall struct{}

func (Syscall) Close(fd int) error {
	return syscall.Close(fd)
}

func (Syscall) EpollCreate1(flag int) (int, error) {
	return syscall.EpollCreate1(flag)
}

func (Syscall) EpollCtl(epfd int, op int, fd int, event *syscall.EpollEvent) error {
	return syscall.EpollCtl(epfd, op, fd, event)
}

func (Syscall) EpollWait(epfd int, events []syscall.EpollEvent, msec int) (int, error) {
	return syscall.EpollWait(epfd, events, msec)
}

func (Syscall) Ioctl(a1, a2, a3 uintptr) (uintptr, uintptr, syscall.Errno) {
	return syscall.Syscall(syscall.SYS_IOCTL, a1, a2, a3)
}

func (Syscall) Open(path string, mode int, perm uint32) (int, error) {
	return syscall.Open(path, mode, perm)
}

func (Syscall) Read(a1, a2, a3 uintptr) (uintptr, uintptr, syscall.Errno) {
	return syscall.Syscall(syscall.SYS_READ, a1, a2, a3)
}

func (Syscall) Socket(domain, typ, proto int) (int, error) {
	return syscall.Socket(domain, typ, proto)
}
