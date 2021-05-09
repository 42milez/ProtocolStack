package ethernet

import (
	"fmt"
	"github.com/42milez/ProtocolStack/src/device"
	"syscall"
	"unsafe"
)

const MaxEpollEvents = 32

// error numbers @ errno-base.h
// https://github.com/torvalds/linux/blob/master/include/uapi/asm-generic/errno-base.h

// struct ifreq @ if.h
// https://github.com/torvalds/linux/blob/e48661230cc35b3d0f4367eddfc19f86463ab917/include/uapi/linux/if.h#L225

// struct sockaddr @ socket.h
// https://github.com/torvalds/linux/blob/master/include/uapi/linux/socket.h

type IfreqFlags struct {
	Name [syscall.IFNAMSIZ]byte
	Flags uint16
}

type IfreqSockAddr struct {
	Name [syscall.IFNAMSIZ]byte
	Addr struct {
		Data [14]byte
		Family uint16
	}
}

func tapOpen(dev *device.Device) error {
	fd, err := syscall.Open("/dev/net/tun", syscall.O_RDWR, 0666)
	if err != nil {
		return err
	}
	dev.Priv.FD = fd

	ifrFlags := IfreqFlags{}
	ifrFlags.Name = dev.Priv.Name
	ifrFlags.Flags = syscall.IFF_TAP | syscall.IFF_NO_PI
	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(fd),
		uintptr(syscall.TUNSETIFF),
		uintptr(unsafe.Pointer(&ifrFlags)))
	if errno != 0 {
		_ = syscall.Close(fd)
		return fmt.Errorf("ioctl error: %d", errno)
	}

	soc, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, 0)
	if err != nil {
		return err
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
		_ = syscall.Close(soc)
		return fmt.Errorf("ioctl error: %d", errno)
	}
	copy(dev.Addr[:], ifrSockAddr.Addr.Data[:])
	_ = syscall.Close(soc)

	return nil
}

func tapClose(dev *device.Device) error {
	return nil
}

func tapTransmit(dev *device.Device) error {
	return nil
}

func tapPoll(dev *device.Device) error {
	//var event syscall.EpollEvent
	//var events [MaxEpollEvents]syscall.EpollEvent
	//
	//epfd, errEpollCreate1 := syscall.EpollCreate1(0)
	//if errEpollCreate1 != nil {
	//	return errEpollCreate1
	//}
	//
	//event.Events = syscall.EPOLLIN
	//event.Fd = int32(dev.Priv.FD)
	//if errEpollCtl := syscall.EpollCtl(epfd, syscall.EPOLL_CTL_ADD, dev.Priv.FD, &event); errEpollCtl != nil {
	//	return errEpollCtl
	//}
	//
	//nEvents, errEpollWait := syscall.EpollWait(epfd, events[:], -1)
	//if errEpollWait != nil {
	//	return errEpollWait
	//}
	//
	//// TODO: send nevents to channel
	//// ...
	//
	//fmt.Printf("nEvents: %d\n", nEvents)

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
