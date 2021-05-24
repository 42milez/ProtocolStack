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

func format(s string) string {
	ret := strings.Replace(s, "\t", "", -1)
	ret = strings.Replace(ret, "\n", "", -1)
	return ret
}

func TestEthAddr_Equal_SUCCESS_Equal(t *testing.T) {
	ethAddr1 := EthAddr([EthAddrLen]byte{11, 22, 33, 44, 55, 66})
	ethAddr2 := EthAddr([EthAddrLen]byte{11, 22, 33, 44, 55, 66})

	got := ethAddr1.Equal(ethAddr2)
	if got != true {
		t.Errorf("EthAddr.Equal() = %v; want %v", got, true)
	}
}

func TestEthAddr_Equal_SUCCESS_NotEqual(t *testing.T) {
	ethAddr1 := EthAddr([EthAddrLen]byte{11, 22, 33, 44, 55, 66})
	ethAddr2 := EthAddr([EthAddrLen]byte{})

	got := ethAddr1.Equal(ethAddr2)
	if got != false {
		t.Errorf("EthAddr.Equal() = %v; want %v", got, false)
	}
}

func TestEthType_String_SUCCESS_A(t *testing.T) {
	ethType := EthType(0x0608)
	want := "ARP"
	got := ethType.String()
	if got != want {
		t.Errorf("EthType.String() = %v; want %v", got, want)
	}
}

func TestEthType_String_SUCCESS_B(t *testing.T) {
	want := "IPv4"
	got := EthType(0x0008).String()
	if got != want {
		t.Errorf("EthType.String() = %v; want %v", got, want)
	}
}

func TestEthType_String_SUCCESS_C(t *testing.T) {
	want := "IPv6"
	got := EthType(0xdd86).String()
	if got != want {
		t.Errorf("EthType.String() = %v; want %v", got, want)
	}
}

func TestEthType_String_SUCCESS_D(t *testing.T) {
	want := "UNKNOWN"
	got := EthType(0x0000).String()
	if got != want {
		t.Errorf("EthType.String() = %v; want %v", got, want)
	}
}

func TestEthDump_SUCCESS(t *testing.T) {
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
				Type: EthType(0x0008), // big endian
			}
			b := new(bytes.Buffer)
			_ = binary.Write(b, binary.BigEndian, hdr)
			copy(buf, b.Bytes())
		}).
		Return(150, nil)

	dev := &Device{Addr: EthAddr{11, 12, 13, 14, 15, 16}}

	got := ReadFrame(dev.Priv.FD, dev.Addr, m)
	if got.Code != psErr.OK {
		t.Errorf("ReadFrame() = %v; want %v", got.Code, psErr.OK)
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

	got := ReadFrame(dev.Priv.FD, dev.Addr, m)
	if got.Code != psErr.CantRead {
		t.Errorf("ReadFrame() = %v; want %v", got.Code, psErr.CantRead)
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

	got := ReadFrame(dev.Priv.FD, dev.Addr, m)
	if got.Code != psErr.InvalidHeader {
		t.Errorf("ReadFrame() = %v; want %v", got.Code, psErr.InvalidHeader)
	}
}

func TestReadFrame_SUCCESS_WhenNoDataToReadExists(t *testing.T) {
	psLog.DisableOutput()
	defer psLog.EnableOutput()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dev := &Device{Addr: EthAddr{33, 44, 55, 66, 77, 88}}
	m := psSyscall.NewMockISyscall(ctrl)
	m.EXPECT().Read(gomock.Any(), gomock.Any()).Return(150, nil)

	got := ReadFrame(dev.Priv.FD, dev.Addr, m)
	if got.Code != psErr.NoDataToRead {
		t.Errorf("ReadFrame() = %v; want %v", got.Code, psErr.NoDataToRead)
	}
}
