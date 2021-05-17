package ethernet

import (
	e "github.com/42milez/ProtocolStack/src/error"
	l "github.com/42milez/ProtocolStack/src/log"
	s "github.com/42milez/ProtocolStack/src/sys"
	"syscall"
	"unsafe"
)

const EpollTimeout = 1000
const MaxEpollEvents = 32

var epfd int

// error numbers @ errno-base.h
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

func tapOpen(dev *Device, sc *s.Syscall) e.Error {
	var err error
	var errno syscall.Errno
	var fd int
	var soc int

	fd, err = sc.Open(vnd, syscall.O_RDWR, 0666)
	if err != nil {
		l.E("can't open virtual networking device: %v ", vnd)
		return e.Error{Code: e.CantOpen}
	}

	// --------------------------------------------------
	ifrFlags := IfreqFlags{}
	ifrFlags.Flags = syscall.IFF_TAP | syscall.IFF_NO_PI
	copy(ifrFlags.Name[:], dev.Priv.Name)

	_, _, errno = sc.Ioctl(uintptr(fd), uintptr(syscall.TUNSETIFF), uintptr(unsafe.Pointer(&ifrFlags)))
	if errno != 0 {
		l.E("SYS_IOCTL (%v) failed: %v ", "TUNSETIFF", errno)
		_ = sc.Close(fd)
		return e.Error{Code: e.CantOpen}
	}

	// --------------------------------------------------
	soc, err = sc.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, 0)
	if err != nil {
		l.E("can't open socket: %v ", err)
		return e.Error{Code: e.CantOpen}
	}

	ifrSockAddr := IfreqSockAddr{}
	ifrSockAddr.Addr.Family = syscall.AF_INET
	copy(ifrSockAddr.Name[:], dev.Priv.Name)

	_, _, errno = sc.Ioctl(uintptr(soc), uintptr(syscall.SIOCGIFHWADDR), uintptr(unsafe.Pointer(&ifrSockAddr)))
	if errno != 0 {
		l.E("SYS_IOCTL (%v) failed: %v ", "SIOCGIFHWADDR", errno)
		_ = sc.Close(soc)
		return e.Error{Code: e.CantOpen}
	}

	dev.Addr = MAC(ifrSockAddr.Addr.Data[:])

	_ = sc.Close(soc)

	// --------------------------------------------------
	var event syscall.EpollEvent

	epfd, err = sc.EpollCreate1(0)
	if err != nil {
		l.E("can't open an epoll file descriptor: %v ", err)
		return e.Error{Code: e.CantOpen}
	}

	event.Events = syscall.EPOLLIN
	event.Fd = int32(fd)

	err = sc.EpollCtl(epfd, syscall.EPOLL_CTL_ADD, fd, &event)
	if err != nil {
		l.E("can't add an entry to the interest list of the epoll file descriptor: %v ", err)
		return e.Error{Code: e.CantOpen}
	}

	dev.Priv.FD = fd

	return e.Error{Code: e.OK}
}

func tapClose(dev *Device, sc *s.Syscall) e.Error {
	_ = sc.Close(epfd)
	return e.Error{Code: e.OK}
}

func tapTransmit(dev *Device, sc *s.Syscall) e.Error {
	return e.Error{Code: e.OK}
}

func tapPoll(dev *Device, sc *s.Syscall, isTerminated bool) e.Error {
	if isTerminated {
		_ = sc.Close(epfd)
		return e.Error{Code: e.OK}
	}

	var events [MaxEpollEvents]syscall.EpollEvent
	nEvents, err := sc.EpollWait(epfd, events[:], EpollTimeout)
	if err != nil {
		_ = sc.Close(epfd)
		return e.Error{Code: e.Interrupted}
	}

	// TODO: send events to channel
	// ...

	// TODO: for development (remove later)
	if nEvents > 0 {
		l.I("events occurred")
		l.I("\tevents: %v ", nEvents)
		l.I("\tdevice: %v (%v) ", dev.Name, dev.Priv.Name)
		_ = ReadFrame(dev, sc)
	} else {
		l.I("no event occurred")
	}

	return e.Error{Code: e.OK}
}

// GenTapDevice generates TAP device object.
func GenTapDevice(name string, mac MAC) (*Device, e.Error) {
	if len(name) > 16 {
		return nil, e.Error{Code: e.CantCreate, Msg: "device name must be less than or equal to 16 characters"}
	}

	dev := &Device{
		Type:      DevTypeEthernet,
		MTU:       EthPayloadSizeMax,
		FLAG:      DevFlagBroadcast | DevFlagNeedArp,
		HeaderLen: EthHeaderSize,
		Addr:      mac,
		AddrLen:   EthAddrLen,
		Broadcast: EthAddrBroadcast,
		Op: Operation{
			Open:     tapOpen,
			Close:    tapClose,
			Transmit: tapTransmit,
			Poll:     tapPoll,
		},
		Priv: Privilege{FD: -1, Name: name},
	}

	return dev, e.Error{Code: e.OK}
}
