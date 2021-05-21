package ethernet

import (
	"bytes"
	"encoding/binary"
	"fmt"
	psErr "github.com/42milez/ProtocolStack/src/error"
	psLog "github.com/42milez/ProtocolStack/src/log"
	psSyscall "github.com/42milez/ProtocolStack/src/syscall"
	"github.com/golang/mock/gomock"
	"regexp"
	"strings"
	"syscall"
	"testing"
	"unsafe"
)

func format(s string) string {
	ret := strings.Replace(s, "\t", "", -1)
	ret = strings.Replace(ret, "\n", "", -1)
	return ret
}

func TestEthAddr_Equal(t *testing.T) {
	ethAddr1 := EthAddr([EthAddrLen]byte{11, 22, 33, 44, 55, 66})
	ethAddr2 := EthAddr([EthAddrLen]byte{11, 22, 33, 44, 55, 66})
	ethAddr3 := EthAddr([EthAddrLen]byte{})

	got := ethAddr1.Equal(ethAddr2)
	if got != true {
		t.Errorf("EthAddr.Equal() = %v; want %v", got, true)
	}

	got = ethAddr1.Equal(ethAddr3)
	if got != false {
		t.Errorf("EthAddr.Equal() = %v; want %v", got, false)
	}
}

func TestEthDump(t *testing.T) {
	regexpDatetime := "[0-9]{4}/[0-9]{2}/[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2}"
	macDst := EthAddr{11, 12, 13, 14, 15, 16}
	macSrc := EthAddr{21, 22, 23, 24, 25, 26}
	ethType := EthType(0x0008)
	want, _ := regexp.Compile(fmt.Sprintf(
		"^.+ %v mac \\(dst\\): %v.+ %v mac \\(src\\): %v.+ %v eth_type:  0x%04x \\(%v\\)$",
		regexpDatetime,
		macDst.String(),
		regexpDatetime,
		macSrc.String(),
		regexpDatetime,
		ethType.String(),
		ethType.String()))
	got := psLog.CaptureLogOutput(func() {
		hdr := EthHeader{Dst: macDst, Src: macSrc, Type: ethType}
		EthDump(&hdr)
	})
	got = format(got)
	if !want.MatchString(got) {
		t.Errorf("EthDump() = %v; want %v", got, want)
	}
}

func TestEthType_String(t *testing.T) {
	ethType := EthType(0x0608)
	want := "ARP"
	got := ethType.String()
	if got != want {
		t.Errorf("EthType.String() = %v; want %v", got, want)
	}

	ethType = EthType(0x0008)
	want = "IPv4"
	got = ethType.String()
	if got != want {
		t.Errorf("EthType.String() = %v; want %v", got, want)
	}

	ethType = EthType(0xdd86)
	want = "IPv6"
	got = ethType.String()
	if got != want {
		t.Errorf("EthType.String() = %v; want %v", got, want)
	}

	ethType = EthType(0x0000)
	want = "UNKNOWN"
	got = ethType.String()
	if got != want {
		t.Errorf("EthType.String() = %v; want %v", got, want)
	}
}

func TestReadFrame1(t *testing.T) {
	psLog.DisableOutput()
	defer psLog.EnableOutput()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r1 := EthHeaderSize * 8
	r2 := 0
	m := psSyscall.NewMockISyscall(ctrl)
	m.EXPECT().
		Read(gomock.Any(), gomock.Any(), gomock.Any()).
		Do(func(_ int, buf unsafe.Pointer, _ int) {
			hdr := EthHeader{
				Dst:  EthAddr{11, 12, 13, 14, 15, 16},
				Src:  EthAddr{21, 22, 23, 24, 25, 26},
				Type: EthType(0x0008), // big endian
			}
			b := new(bytes.Buffer)
			_ = binary.Write(b, binary.BigEndian, hdr)
			copy((*(*[]byte)(buf))[:], b.Bytes())
		}).
		Return(uintptr(r1), uintptr(r2), syscall.Errno(0))

	dev := &Device{Addr: EthAddr{11, 12, 13, 14, 15, 16}}

	got := ReadFrame(dev.Priv.FD, dev.Addr, m)
	if got.Code != psErr.OK {
		t.Errorf("ReadFrame() = %v; want %v", got.Code, psErr.OK)
	}
}

func TestReadFrame2(t *testing.T) {
	psLog.DisableOutput()
	defer psLog.EnableOutput()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dev := &Device{Addr: EthAddr{11, 12, 13, 14, 15, 16}}
	hdrLen := EthHeaderSize * 8
	m := psSyscall.NewMockISyscall(ctrl)
	m.EXPECT().Read(gomock.Any(), gomock.Any(), gomock.Any()).Return(uintptr(hdrLen), uintptr(0), syscall.EINTR)

	got := ReadFrame(dev.Priv.FD, dev.Addr, m)
	if got.Code != psErr.CantRead {
		t.Errorf("ReadFrame() = %v; want %v", got.Code, psErr.CantRead)
	}
}

func TestReadFrame3(t *testing.T) {
	psLog.DisableOutput()
	defer psLog.EnableOutput()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dev := &Device{Addr: EthAddr{11, 12, 13, 14, 15, 16}}
	hdrLen := 0
	m := psSyscall.NewMockISyscall(ctrl)
	m.EXPECT().
		Read(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(uintptr(hdrLen), uintptr(0), syscall.Errno(0))

	got := ReadFrame(dev.Priv.FD, dev.Addr, m)
	if got.Code != psErr.InvalidHeader {
		t.Errorf("ReadFrame() = %v; want %v", got.Code, psErr.InvalidHeader)
	}
}

func TestReadFrame4(t *testing.T) {
	psLog.DisableOutput()
	defer psLog.EnableOutput()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dev := &Device{Addr: EthAddr{33, 44, 55, 66, 77, 88}}
	hdrLen := EthHeaderSize * 8
	m := psSyscall.NewMockISyscall(ctrl)
	m.EXPECT().Read(gomock.Any(), gomock.Any(), gomock.Any()).Return(uintptr(hdrLen), uintptr(0), syscall.Errno(0))

	got := ReadFrame(dev.Priv.FD, dev.Addr, m)
	if got.Code != psErr.NoDataToRead {
		t.Errorf("ReadFrame() = %v; want %v", got.Code, psErr.NoDataToRead)
	}
}
