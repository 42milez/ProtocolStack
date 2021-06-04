// +build amd64,linux

package eth

import (
	"errors"
	"fmt"
	psErr "github.com/42milez/ProtocolStack/src/error"
	psLog "github.com/42milez/ProtocolStack/src/log"
	psSyscall "github.com/42milez/ProtocolStack/src/syscall"
	"syscall"
	"unsafe"
)

const EpollTimeout = 1000
const MaxEpollEvents = 32

// src/syscall/zerrors_linux_amd64.go
// https://golang.org/src/syscall/zerrors_linux_amd64.go

// error numbers
// https://github.com/torvalds/linux/blob/master/include/uapi/asm-generic/errno-base.h

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

type TapDevice struct {
	Device
}

func (p *TapDevice) Open() psErr.E {
	var fd int
	var err error

	fd, err = psSyscall.Syscall.Open(virtualNetworkDevice, syscall.O_RDWR, 0666)
	if err != nil {
		return psErr.CantOpenIOResource
	}

	ifrFlags := IfreqFlags{}
	ifrFlags.Flags = syscall.IFF_TAP | syscall.IFF_NO_PI
	copy(ifrFlags.Name[:], p.Priv().Name)
	if errno := psSyscall.Syscall.Ioctl(fd, syscall.TUNSETIFF, unsafe.Pointer(&ifrFlags)); errno != 0 {
		_ = psSyscall.Syscall.Close(fd)
		return psErr.CantModifyIOResourceParameter
	}

	// --------------------------------------------------

	var soc int
	soc, err = psSyscall.Syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, 0)
	if err != nil {
		_ = psSyscall.Syscall.Close(fd)
		return psErr.CantCreateEndpoint
	}

	ifrSockAddr := IfreqSockAddr{}
	ifrSockAddr.Addr.Family = syscall.AF_INET
	copy(ifrSockAddr.Name[:], p.Priv().Name)

	if errno := psSyscall.Syscall.Ioctl(soc, syscall.SIOCGIFHWADDR, unsafe.Pointer(&ifrSockAddr)); errno != 0 {
		_ = psSyscall.Syscall.Close(soc)
		return psErr.CantModifyIOResourceParameter
	}
	copy(p.Addr_[:], ifrSockAddr.Addr.Data[:])
	if err = psSyscall.Syscall.Close(soc); err != nil {
		return psErr.CantCloseIOResource
	}

	// --------------------------------------------------

	epfd, err = psSyscall.Syscall.EpollCreate1(0)
	if err != nil {
		return psErr.CantCreateEpollInstance
	}

	var event syscall.EpollEvent
	event.Events = syscall.EPOLLIN
	event.Fd = int32(fd)

	if err := psSyscall.Syscall.EpollCtl(epfd, syscall.EPOLL_CTL_ADD, fd, &event); err != nil {
		_ = psSyscall.Syscall.Close(epfd)
		return psErr.CantModifyIOResourceParameter
	}

	p.Priv_.FD = fd

	return psErr.OK
}

func (p *TapDevice) Close() psErr.E {
	if err := psSyscall.Syscall.Close(epfd); err != nil {
		return psErr.SyscallError
	}
	return psErr.OK
}

func (p *TapDevice) Poll(isTerminated bool) psErr.E {
	if isTerminated {
		_ = psSyscall.Syscall.Close(epfd)
		return psErr.Terminated
	}

	var events [MaxEpollEvents]syscall.EpollEvent
	nEvents, err := psSyscall.Syscall.EpollWait(epfd, events[:], EpollTimeout)
	if err != nil {
		// https://man7.org/linux/man-pages/man2/epoll_wait.2.html#RETURN_VALUE
		// ignore EINTR
		if !errors.Is(err, syscall.EINTR) {
			return psErr.SyscallError
		}
		return psErr.Interrupted
	}

	if nEvents > 0 {
		psLog.I("Event occurred")
		psLog.I(fmt.Sprintf("\tevents: %v", nEvents))
		psLog.I(fmt.Sprintf("\tdevice: %v (%v)", p.Name_, p.Priv_.Name))
		if packet, err := ReadEthFrame(p.Priv_.FD, p.Addr_); err != psErr.OK {
			if err != psErr.NoDataToRead {
				return psErr.Error
			}
		} else {
			packet.Dev = p
			RxCh <- packet
		}
	}

	return psErr.OK
}

func (p *TapDevice) Transmit(dst Addr, payload []byte, typ Type) psErr.E {
	return WriteEthFrame(p.Priv_.FD, dst, p.Addr_, typ, payload)
}

const virtualNetworkDevice = "/dev/net/tun"

var epfd int
