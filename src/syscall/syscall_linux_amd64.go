//go:generate mockgen -source=syscall_linux_amd64.go -destination=syscall_mock_linux_amd64.go -package=$GOPACKAGE -self_package=github.com/42milez/ProtocolStack/src/$GOPACKAGE

package syscall

import (
	"syscall"
	"unsafe"
)

var Syscall ISyscall

type ISyscall interface {
	Open(path string, mode int, perm uint32) (fd int, err error)
	Close(fd int) (err error)
	EpollCreate1(flag int) (fd int, err error)
	EpollCtl(epfd int, op int, fd int, event *syscall.EpollEvent) (err error)
	EpollWait(epfd int, events []syscall.EpollEvent, msec int) (n int, err error)
	Ioctl(fd int, code int, data unsafe.Pointer) (err syscall.Errno)
	Socket(domain, typ, proto int) (fd int, err error)
	Read(fd int, p []byte) (n int, err error)
	Write(fd int, p []byte) (n int, err error)
}

type scImpl struct{}

func (scImpl) Open(path string, mode int, perm uint32) (fd int, err error) {
	return syscall.Open(path, mode, perm)
}

func (scImpl) Close(fd int) (err error) {
	return syscall.Close(fd)
}

func (scImpl) EpollCreate1(flag int) (fd int, err error) {
	return syscall.EpollCreate1(flag)
}

func (scImpl) EpollCtl(epfd int, op int, fd int, event *syscall.EpollEvent) (err error) {
	return syscall.EpollCtl(epfd, op, fd, event)
}

func (scImpl) EpollWait(epfd int, events []syscall.EpollEvent, msec int) (n int, err error) {
	return syscall.EpollWait(epfd, events, msec)
}

func (scImpl) Ioctl(fd int, code int, data unsafe.Pointer) (err syscall.Errno) {
	// doc: second return value of syscall.Syscall needs to be documented #29842
	// https://github.com/golang/go/issues/29842
	_, _, err = syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), uintptr(code), uintptr(data))
	return
}

func (scImpl) Socket(domain, typ, proto int) (fd int, err error) {
	return syscall.Socket(domain, typ, proto)
}

func (scImpl) Read(fd int, p []byte) (n int, err error) {
	return syscall.Read(fd, p)
}

func (scImpl) Write(fd int, p []byte) (n int, err error) {
	return syscall.Write(fd, p)
}

func init() {
	Syscall = &scImpl{}
}
