package ethernet

import (
	"errors"
	"fmt"
	"github.com/42milez/ProtocolStack/src/device"
	"syscall"
	"unsafe"
)

const iffTap = 0x0002
const iffNoPi = 0x1000
const tunSetIff = 0x400454ca

const afInet = 0x2
const siocgifhwaddr = 0x8927
const sockDgram = 0x2

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
	ifrFlags.Flags = iffTap | iffNoPi
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), uintptr(tunSetIff), uintptr(unsafe.Pointer(&ifrFlags)))
	if errno != 0 {
		_ = syscall.Close(fd)
		return errors.New(fmt.Sprintf("ioctl error: %d", errno))
	}

	soc, err := syscall.Socket(afInet, sockDgram, 0)
	if err != nil {
		return err
	}
	ifrSockAddr := IfreqSockAddr{}
	copy(ifrSockAddr.Name[:], dev.Name)
	ifrSockAddr.Addr.Family = syscall.AF_INET
	_, _, errno = syscall.Syscall(syscall.SYS_IOCTL, uintptr(soc), uintptr(siocgifhwaddr), uintptr(unsafe.Pointer(&ifrSockAddr)))
	if errno != 0 {
		_ = syscall.Close(soc)
		return errors.New(fmt.Sprintf("ioctl error: %d", errno))
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
	return nil
}

// GenTapDevice generates TAP device object.
func GenTapDevice(name string, mac MAC) (*device.Device, error) {
	if len(name) > 16 {
		return nil, errors.New("device name must be less than or equal to 16 characters")
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
