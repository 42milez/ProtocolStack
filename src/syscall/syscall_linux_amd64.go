//go:generate mockgen -source=syscall_linux_amd64.go -destination=syscall_mock_linux_amd64.go -package=$GOPACKAGE -self_package=github.com/42milez/ProtocolStack/src/$GOPACKAGE

package syscall

import (
	"syscall"
	"unsafe"
)

type ISyscall interface {
	Close(fd int) error
	EpollCreate1(flag int) (int, error)
	EpollCtl(epfd int, op int, fd int, event *syscall.EpollEvent) error
	EpollWait(epfd int, events []syscall.EpollEvent, msec int) (int, error)
	Ioctl(a1, a2, a3 uintptr) (uintptr, uintptr, syscall.Errno)
	Open(path string, mode int, perm uint32) (int, error)
	Read(fd int, buf unsafe.Pointer, size int) (uintptr, uintptr, syscall.Errno)
	Socket(domain, typ, proto int) (int, error)
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

func (Syscall) Read(fd int, buf unsafe.Pointer, size int) (uintptr, uintptr, syscall.Errno) {
	return syscall.Syscall(syscall.SYS_READ, uintptr(fd), uintptr(buf), uintptr(size))
}

func (Syscall) Socket(domain, typ, proto int) (int, error) {
	return syscall.Socket(domain, typ, proto)
}
