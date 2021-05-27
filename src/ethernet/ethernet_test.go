package ethernet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	psErr "github.com/42milez/ProtocolStack/src/error"
	psLog "github.com/42milez/ProtocolStack/src/log"
	psSyscall "github.com/42milez/ProtocolStack/src/syscall"
	"github.com/golang/mock/gomock"
	"regexp"
	"strings"
	"testing"
)

func trim(s string) string {
	ret := strings.Replace(s, "\t", "", -1)
	ret = strings.Replace(ret, "\n", "", -1)
	return ret
}

func TestEthAddr_Equal(t *testing.T) {
	ethAddr1 := EthAddr([EthAddrLen]byte{11, 22, 33, 44, 55, 66})
	ethAddr2 := EthAddr([EthAddrLen]byte{11, 22, 33, 44, 55, 66})
	got := ethAddr1.Equal(ethAddr2)
	if got != true {
		t.Errorf("EthAddr.Equal() = %t; want %t", got, true)
	}

	ethAddr1 = EthAddr([EthAddrLen]byte{11, 22, 33, 44, 55, 66})
	ethAddr2 = EthAddr([EthAddrLen]byte{})
	got = ethAddr1.Equal(ethAddr2)
	if got != false {
		t.Errorf("EthAddr.Equal() = %t; want %t", got, false)
	}
}

func TestEthType_String(t *testing.T) {
	want := "ARP"
	got := EthType(0x0806).String()
	if got != want {
		t.Errorf("EthType.String() = %s; want %s", got, want)
	}

	want = "IPv4"
	got = EthType(0x0800).String()
	if got != want {
		t.Errorf("EthType.String() = %s; want %s", got, want)
	}

	want = "IPv6"
	got = EthType(0x86dd).String()
	if got != want {
		t.Errorf("EthType.String() = %s; want %s", got, want)
	}

	want = "UNKNOWN"
	got = EthType(0x0000).String()
	if got != want {
		t.Errorf("EthType.String() = %s; want %s", got, want)
	}
}

func TestEthDump(t *testing.T) {
	regexpDatetime := "[0-9]{4}/[0-9]{2}/[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2}"
	macDst := EthAddr{0x01, 0x02, 0x03, 0x04, 0x05, 0x06}
	macSrc := EthAddr{0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}
	ethType := EthType(0x0800)
	want, _ := regexp.Compile(fmt.Sprintf(
		"^.+ %v dst:  %v.+ %v src:  %v.+ %v type: 0x%04x \\(%v\\)$",
		regexpDatetime,
		macDst.String(),
		regexpDatetime,
		macSrc.String(),
		regexpDatetime,
		uint16(ethType),
		ethType.String()))
	got := psLog.CaptureLogOutput(func() {
		hdr := EthHeader{Dst: macDst, Src: macSrc, Type: ethType}
		EthDump(&hdr)
	})
	got = trim(got)
	if !want.MatchString(got) {
		t.Errorf("EthDump() = %v; want %v", got, want)
	}
}

func TestReadFrame_SUCCESS(t *testing.T) {
	psLog.DisableOutput()
	defer psLog.EnableOutput()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := psSyscall.NewMockISyscall(ctrl)
	m.EXPECT().
		Read(gomock.Any(), gomock.Any()).
		Do(func(_ int, buf []byte) {
			hdr := EthHeader{
				Dst:  EthAddr{11, 12, 13, 14, 15, 16},
				Src:  EthAddr{21, 22, 23, 24, 25, 26},
				Type: EthType(0x0800),
			}
			b := new(bytes.Buffer)
			_ = binary.Write(b, binary.BigEndian, hdr)
			copy(buf, b.Bytes())
		}).
		Return(150, nil)

	dev := &Device{Addr: EthAddr{11, 12, 13, 14, 15, 16}}

	_, got := ReadFrame(dev.Priv.FD, dev.Addr, m)
	if got != psErr.OK {
		t.Errorf("ReadFrame() = %v; want %v", got, psErr.OK)
	}
}

func TestReadFrame_FAIL_WhenReadSyscallFailed(t *testing.T) {
	psLog.DisableOutput()
	defer psLog.EnableOutput()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dev := &Device{Addr: EthAddr{11, 12, 13, 14, 15, 16}}
	m := psSyscall.NewMockISyscall(ctrl)
	m.EXPECT().Read(gomock.Any(), gomock.Any()).Return(-1, errors.New(""))

	_, got := ReadFrame(dev.Priv.FD, dev.Addr, m)
	if got != psErr.Error {
		t.Errorf("ReadFrame() = %v; want %v", got, psErr.Error)
	}
}

func TestReadFrame_FAIL_WhenHeaderLengthIsInvalid(t *testing.T) {
	psLog.DisableOutput()
	defer psLog.EnableOutput()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dev := &Device{Addr: EthAddr{11, 12, 13, 14, 15, 16}}
	m := psSyscall.NewMockISyscall(ctrl)
	m.EXPECT().Read(gomock.Any(), gomock.Any()).Return(10, nil)

	_, got := ReadFrame(dev.Priv.FD, dev.Addr, m)
	if got != psErr.Error {
		t.Errorf("ReadFrame() = %v; want %v", got, psErr.Error)
	}
}

func TestReadFrame_SUCCESS_WhenNoDataExists(t *testing.T) {
	psLog.DisableOutput()
	defer psLog.EnableOutput()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dev := &Device{Addr: EthAddr{33, 44, 55, 66, 77, 88}}
	m := psSyscall.NewMockISyscall(ctrl)
	m.EXPECT().Read(gomock.Any(), gomock.Any()).Return(150, nil)

	_, got := ReadFrame(dev.Priv.FD, dev.Addr, m)
	if got != psErr.NoDataToRead {
		t.Errorf("ReadFrame() = %v; want %v", got, psErr.NoDataToRead)
	}
}
