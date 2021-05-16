package syscall

import "syscall"

type EpollCreate1IF interface {
	Exec() (int, error)
}

type EpollCreate1Syscall struct {
	Flag int
}

func (e EpollCreate1Syscall) Exec() (int, error) {
	return syscall.EpollCreate1(0)
}
