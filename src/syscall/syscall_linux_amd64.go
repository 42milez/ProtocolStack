//go:generate mockgen -source=syscall_linux_amd64.go -destination=syscall_mock_linux_amd64.go -package=$GOPACKAGE -self_package=github.com/42milez/ProtocolStack/src/$GOPACKAGE

package syscall

import (
	"syscall"
)

type ISyscall interface {
	Close(fd int) (err error)
	EpollCreate1(flag int) (fd int, err error)
	EpollCtl(epfd int, op int, fd int, event *syscall.EpollEvent) (err error)
	EpollWait(epfd int, events []syscall.EpollEvent, msec int) (n int, err error)
	Ioctl(a1, a2, a3 uintptr) (r1, r2 uintptr, err syscall.Errno)
	Open(path string, mode int, perm uint32) (fd int, err error)
	Read(fd int, p []byte) (n int, err error)
	Socket(domain, typ, proto int) (fd int, err error)
}

type Syscall struct{}

func (Syscall) Close(fd int) (err error) {
	return syscall.Close(fd)
}

func (Syscall) EpollCreate1(flag int) (fd int, err error) {
	return syscall.EpollCreate1(flag)
}

func (Syscall) EpollCtl(epfd int, op int, fd int, event *syscall.EpollEvent) (err error) {
	return syscall.EpollCtl(epfd, op, fd, event)
}

func (Syscall) EpollWait(epfd int, events []syscall.EpollEvent, msec int) (n int, err error) {
	return syscall.EpollWait(epfd, events, msec)
}

func (Syscall) Ioctl(a1, a2, a3 uintptr) (r1, r2 uintptr, err syscall.Errno) {
	return syscall.Syscall(syscall.SYS_IOCTL, a1, a2, a3)
}

func (Syscall) Open(path string, mode int, perm uint32) (fd int, err error) {
	return syscall.Open(path, mode, perm)
}

func (Syscall) Read(fd int, p []byte) (n int, err error) {
	return syscall.Read(fd, p)
}

func (Syscall) Socket(domain, typ, proto int) (fd int, err error) {
	return syscall.Socket(domain, typ, proto)
}
