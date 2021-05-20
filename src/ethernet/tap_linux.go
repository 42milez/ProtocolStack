// +build amd64,linux

package ethernet

import (
	"fmt"
	psErr "github.com/42milez/ProtocolStack/src/error"
	psLog "github.com/42milez/ProtocolStack/src/log"
	"syscall"
	"unsafe"
)

const EpollTimeout = 1000
const MaxEpollEvents = 32

var epfd int

// error numbers @ errno-baspsErr.h
// https://github.com/torvalds/linux/blob/master/include/uapi/asm-generic/errno-baspsErr.h

// struct ifreq @ if.h
// https://github.com/torvalds/linux/blob/e48661230cc35b3d0f4367eddfc19f86463ab917/include/uapi/linux/if.h#L225

// struct sockaddr @ socket.h
// https://github.com/torvalds/linux/blob/master/include/uapi/linux/socket.h

type IfreqFlags struct {
	Name  [syscall.IFNAMSIZ]byte
	Flags uint16
}

type IfreqSockAddr struct {
	Name [syscall.IFNAMSIZ]byte
	Addr struct {
		Data   [14]byte
		Family uint16
	}
}

const vnd = "/dev/net/tun"

type TapDevice struct {
	Device
}

func (dev *TapDevice) Open() psErr.Error {
	var err error
	var errno syscall.Errno
	var fd int
	var soc int

	fd, err = dev.Syscall.Open(vnd, syscall.O_RDWR, 0666)
	if err != nil {
		psLog.E("can't open virtual networking device: %v ", vnd)
		return psErr.Error{Code: psErr.CantOpen, Msg: err.Error()}
	}

	// --------------------------------------------------
	ifrFlags := IfreqFlags{}
	ifrFlags.Flags = syscall.IFF_TAP | syscall.IFF_NO_PI
	copy(ifrFlags.Name[:], dev.Priv.Name)

	_, _, errno = dev.Syscall.Ioctl(uintptr(fd), uintptr(syscall.TUNSETIFF), uintptr(unsafe.Pointer(&ifrFlags)))
	if errno != 0 {
		psLog.E("SYS_IOCTL (%v) failed: %v ", "TUNSETIFF", errno)
		_ = dev.Syscall.Close(fd)
		return psErr.Error{Code: psErr.CantOpen, Msg: fmt.Sprintf("errno: %v", errno)}
	}

	// --------------------------------------------------
	soc, err = dev.Syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, 0)
	if err != nil {
		psLog.E("can't open socket: %v ", err)
		return psErr.Error{Code: psErr.CantOpen, Msg: err.Error()}
	}

	ifrSockAddr := IfreqSockAddr{}
	ifrSockAddr.Addr.Family = syscall.AF_INET
	copy(ifrSockAddr.Name[:], dev.Priv.Name)

	_, _, errno = dev.Syscall.Ioctl(uintptr(soc), uintptr(syscall.SIOCGIFHWADDR), uintptr(unsafe.Pointer(&ifrSockAddr)))
	if errno != 0 {
		psLog.E("SYS_IOCTL (%v) failed: %v ", "SIOCGIFHWADDR", errno)
		_ = dev.Syscall.Close(soc)
		return psErr.Error{Code: psErr.CantOpen, Msg: fmt.Sprintf("errno: %v", errno)}
	}

	copy(dev.Addr[:], ifrSockAddr.Addr.Data[:])

	_ = dev.Syscall.Close(soc)

	// --------------------------------------------------
	var event syscall.EpollEvent

	epfd, err = dev.Syscall.EpollCreate1(0)
	if err != nil {
		psLog.E("can't open an epoll file descriptor: %v ", err)
		return psErr.Error{Code: psErr.CantOpen, Msg: err.Error()}
	}

	event.Events = syscall.EPOLLIN
	event.Fd = int32(fd)

	err = dev.Syscall.EpollCtl(epfd, syscall.EPOLL_CTL_ADD, fd, &event)
	if err != nil {
		psLog.E("can't add an entry to the interest list of the epoll file descriptor: %v ", err)
		return psErr.Error{Code: psErr.CantOpen, Msg: err.Error()}
	}

	dev.Priv.FD = fd

	return psErr.Error{Code: psErr.OK}
}

func (dev *TapDevice) Close() psErr.Error {
	_ = dev.Syscall.Close(epfd)
	return psErr.Error{Code: psErr.OK}
}

func (dev *TapDevice) Poll(isTerminated bool) psErr.Error {
	if isTerminated {
		_ = dev.Syscall.Close(epfd)
		return psErr.Error{Code: psErr.OK}
	}

	var events [MaxEpollEvents]syscall.EpollEvent
	nEvents, err := dev.Syscall.EpollWait(epfd, events[:], EpollTimeout)
	if err != nil {
		_ = dev.Syscall.Close(epfd)
		return psErr.Error{Code: psErr.Interrupted}
	}

	// TODO: send events to channel
	// ...

	// TODO: for development (remove later)
	if nEvents > 0 {
		psLog.I("events occurred")
		psLog.I("\tevents: %v ", nEvents)
		psLog.I("\tdevice: %v (%v) ", dev.Name, dev.Priv.Name)
		_ = ReadFrame(dev.Priv.FD, dev.Addr, dev.Syscall)
	} else {
		psLog.I("no event occurred")
	}

	return psErr.Error{Code: psErr.OK}
}

func (dev *TapDevice) Transmit() psErr.Error {
	return psErr.Error{Code: psErr.OK}
}
