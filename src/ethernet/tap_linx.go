package ethernet

import (
	"fmt"
	"github.com/42milez/ProtocolStack/src/device"
	e "github.com/42milez/ProtocolStack/src/error"
	"log"
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

func tapOpen(dev *device.Device) error {
	var err error
	var errno syscall.Errno
	var fd int
	var soc int

	fd, err = syscall.Open(vnd, syscall.O_RDWR, 0666)
	if err != nil {
		log.Printf("can't open virtual networking device: %v\n", vnd)
		return e.CantOpen
	}

	// --------------------------------------------------
	ifrFlags := IfreqFlags{}
	ifrFlags.Name = dev.Priv.Name
	ifrFlags.Flags = syscall.IFF_TAP | syscall.IFF_NO_PI

	_, _, errno = syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(fd),
		uintptr(syscall.TUNSETIFF),
		uintptr(unsafe.Pointer(&ifrFlags)))
	if errno != 0 {
		log.Printf("SYS_IOCTL (%v) failed: %v\n", "TUNSETIFF", errno)
		_ = syscall.Close(fd)
		return e.CantOpen
	}

	// --------------------------------------------------
	soc, err = syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, 0)
	if err != nil {
		log.Printf("can't open socket: %v\n", err)
		return e.CantOpen
	}

	ifrSockAddr := IfreqSockAddr{}
	ifrSockAddr.Name = dev.Priv.Name
	ifrSockAddr.Addr.Family = syscall.AF_INET

	_, _, errno = syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(soc),
		uintptr(syscall.SIOCGIFHWADDR),
		uintptr(unsafe.Pointer(&ifrSockAddr)))
	if errno != 0 {
		log.Printf("SYS_IOCTL (%v) failed: %v\n", "SIOCGIFHWADDR", errno)
		_ = syscall.Close(soc)
		return e.CantOpen
	}

	copy(dev.Addr[:], ifrSockAddr.Addr.Data[:])
	_ = syscall.Close(soc)

	// --------------------------------------------------
	var event syscall.EpollEvent

	epfd, err = syscall.EpollCreate1(0)
	if err != nil {
		log.Printf("can't open an epoll file descriptor: %v\n", err)
		return e.CantOpen
	}

	event.Events = syscall.EPOLLIN
	event.Fd = int32(fd)

	err = syscall.EpollCtl(epfd, syscall.EPOLL_CTL_ADD, fd, &event)
	if err != nil {
		log.Printf("can't add an entry to the interest list of the epoll file descriptor: %v\n", err)
		return e.CantOpen
	}

	dev.Priv.FD = fd

	return e.OK
}

func tapClose(dev *device.Device) error {
	return syscall.Close(epfd)
}

func tapTransmit(dev *device.Device) error {
	return nil
}

func tapPoll(dev *device.Device, isTerminated bool) error {
	if isTerminated {
		_ = syscall.Close(epfd)
		return nil
	}

	var events [MaxEpollEvents]syscall.EpollEvent
	nEvents, err := syscall.EpollWait(epfd, events[:], EpollTimeout)
	if err != nil {
		_ = syscall.Close(epfd)
		return err
	}

	// TODO: send events to channel
	// ...

	// TODO: for development (remove later)
	if nEvents > 0 {
		log.Println("events occurred")
		log.Printf("\tEvents: %v\n", nEvents)
		log.Printf("\tDevice: %v (%v)\n", dev.Name, dev.Priv.Name)
		_ = ReadFrame(dev)
	} else {
		log.Printf("no event occurred")
	}

	return nil
}

// GenTapDevice generates TAP device object.
func GenTapDevice(name string, mac MAC) (*device.Device, error) {
	if len(name) > 16 {
		return nil, fmt.Errorf("device name must be less than or equal to 16 characters")
	}

	dev := &device.Device{
		Type:      device.DevTypeEthernet,
		MTU:       EthPayloadSizeMax,
		FLAG:      device.DevFlagBroadcast | device.DevFlagNeedArp,
		HeaderLen: EthHeaderSize,
		AddrLen:   EthAddrLen,
		Broadcast: EthAddrBroadcast,
		Op: device.Operation{
			Open:     tapOpen,
			Close:    tapClose,
			Transmit: tapTransmit,
			Poll:     tapPoll,
		},
		Priv: device.Privilege{FD: -1},
	}
	copy(dev.Priv.Name[:], name)

	if addr, err := mac.Byte(); err != nil {
		return nil, err
	} else {
		dev.Addr = addr
	}

	return dev, nil
}
