// +build amd64,linux

package ethernet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	psErr "github.com/42milez/ProtocolStack/src/error"
	psLog "github.com/42milez/ProtocolStack/src/log"
	"syscall"
	"unsafe"
)

const EpollTimeout = 1000
const MaxEpollEvents = 32

var epfd int

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

const vnd = "/dev/net/tun"

type TapDevice struct {
	Device
}

func (dev *TapDevice) Open() psErr.E {
	var err error
	var fd int

	fd, err = dev.Syscall.Open(vnd, syscall.O_RDWR, 0666)
	if err != nil {
		psLog.E("syscall.Open() failed: %s", err)
		return psErr.CantOpen
	}

	// --------------------------------------------------

	ifrFlags := IfreqFlags{}
	ifrFlags.Flags = syscall.IFF_TAP | syscall.IFF_NO_PI
	copy(ifrFlags.Name[:], dev.Priv.Name)

	if _, _, errno := dev.Syscall.Ioctl(uintptr(fd), uintptr(syscall.TUNSETIFF), uintptr(unsafe.Pointer(&ifrFlags))); errno != 0 {
		_ = dev.Syscall.Close(fd)
		psLog.E("syscall.Syscall(SYS_IOCTL, TUNSETIFF) failed: %s", errno)
		return psErr.CantInitialize
	}

	// --------------------------------------------------

	var soc int
	soc, err = dev.Syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, 0)
	if err != nil {
		psLog.E("syscall.Socket() failed: %s", err)
		return psErr.CantInitialize
	}

	ifrSockAddr := IfreqSockAddr{}
	ifrSockAddr.Addr.Family = syscall.AF_INET
	copy(ifrSockAddr.Name[:], dev.Priv.Name)

	if _, _, errno := dev.Syscall.Ioctl(uintptr(soc), uintptr(syscall.SIOCGIFHWADDR), uintptr(unsafe.Pointer(&ifrSockAddr))); errno != 0 {
		_ = dev.Syscall.Close(soc)
		psLog.E("syscall.Syscall(SYS_IOCTL, SIOCGIFHWADDR) failed: %s", errno)
		return psErr.CantInitialize
	}
	copy(dev.Addr[:], ifrSockAddr.Addr.Data[:])
	_ = dev.Syscall.Close(soc)

	// --------------------------------------------------

	epfd, err = dev.Syscall.EpollCreate1(0)
	if err != nil {
		psLog.E("syscall.EpollCreate1() failed: %s", err)
		return psErr.CantInitialize
	}

	var event syscall.EpollEvent
	event.Events = syscall.EPOLLIN
	event.Fd = int32(fd)

	if e := dev.Syscall.EpollCtl(epfd, syscall.EPOLL_CTL_ADD, fd, &event); e != nil {
		_ = dev.Syscall.Close(epfd)
		psLog.E("syscall.EpollCtl() failed: %s", err)
		return psErr.CantInitialize
	}

	dev.Priv.FD = fd

	return psErr.OK
}

func (dev *TapDevice) Close() psErr.E {
	_ = dev.Syscall.Close(epfd)
	return psErr.OK
}

func (dev *TapDevice) Poll(isTerminated bool) psErr.E {
	if isTerminated {
		_ = dev.Syscall.Close(epfd)
		return psErr.Terminated
	}

	var events [MaxEpollEvents]syscall.EpollEvent
	nEvents, err := dev.Syscall.EpollWait(epfd, events[:], EpollTimeout)
	if err != nil {
		// https://man7.org/linux/man-pages/man2/epoll_wait.2.html#RETURN_VALUE
		// ignore EINTR
		if !errors.Is(err, syscall.EINTR) {
			_ = dev.Syscall.Close(epfd)
			return psErr.Error
		}
		psLog.I("Syscall.EpollWait() was interrupted")
	}

	if nEvents > 0 {
		psLog.I("Events occurred")
		psLog.I("\tevents: %v", nEvents)
		psLog.I("\tdevice: %v (%v)", dev.Name, dev.Priv.Name)
		if packet, err := ReadFrame(dev.Priv.FD, dev.Addr, dev.Syscall); err != psErr.OK {
			if err != psErr.NoDataToRead {
				psLog.E("ReadFrame() failed: %s", err)
				return psErr.CantRead
			}
		} else {
			packet.Dev = dev
			RxCh <- packet
		}
	}

	return psErr.OK
}

func (dev *TapDevice) Transmit(dest EthAddr, payload []byte, typ EthType) psErr.E {
	hdr := EthHeader{
		Dst:  dest,
		Src:  dev.Addr,
		Type: typ,
	}

	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, &hdr); err != nil {
		psLog.E("binary.Write() failed: %s", err)
		return psErr.Error
	}
	if err := binary.Write(buf, binary.BigEndian, &payload); err != nil {
		psLog.E("binary.Write() failed: %s", err)
		return psErr.Error
	}

	if fsize := buf.Len(); fsize < EthFrameSizeMin {
		pad := make([]byte, EthFrameSizeMin-fsize)
		if err := binary.Write(buf, binary.BigEndian, &pad); err != nil {
			psLog.E("binary.Write() failed: %s", err)
			return psErr.Error
		}
	}

	psLog.I("Ethernet frame to be sent")
	psLog.I("\tdest:    %s", hdr.Dst)
	psLog.I("\tsrc:     %s", hdr.Src)
	psLog.I("\ttype:    %s", hdr.Type)
	s := "\tpayload: "
	for i, v := range payload {
		s += fmt.Sprintf("%02x", v)
		if (i+1)%10 == 0 {
			psLog.I("%s", s)
			s = "\t\t "
		}
	}

	if n, err := dev.Syscall.Write(dev.Priv.FD, buf.Bytes()); err != nil {
		psLog.E("syscall.Write() failed: %s", err)
		return psErr.Error
	} else {
		psLog.I("Ethernet frame has been written: %d bytes", n)
	}

	return psErr.OK
}
