package mw

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

func TestEthAddr_Equal(t *testing.T) {
	ethAddr1 := EthAddr{0x01, 0x02, 0x03, 0x04, 0x05, 0x06}
	ethAddr2 := EthAddr{0x01, 0x02, 0x03, 0x04, 0x05, 0x06}
	got := ethAddr1.Equal(ethAddr2)
	if got != true {
		t.Errorf("EthAddr.Equal() = %t; want %t", got, true)
	}

	ethAddr3 := EthAddr{}
	got = ethAddr1.Equal(ethAddr3)
	if got != false {
		t.Errorf("EthAddr.Equal() = %t; want %t", got, false)
	}
}

func TestEthType_String(t *testing.T) {
	want := "IPv4"
	got := EthType(0x0800).String()
	if got != want {
		t.Errorf("Type.String() = %s; want %s", got, want)
	}

	want = "ARP"
	got = EthType(0x0806).String()
	if got != want {
		t.Errorf("Type.String() = %s; want %s", got, want)
	}

	want = "IPv6"
	got = EthType(0x86dd).String()
	if got != want {
		t.Errorf("Type.String() = %s; want %s", got, want)
	}
}

func TestDumpFrame(t *testing.T) {
	macDst := EthAddr{0x01, 0x02, 0x03, 0x04, 0x05, 0x06}
	macSrc := EthAddr{0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}
	ethType := EthType(0x0800)
	want, _ := regexp.Compile(fmt.Sprintf(
		"^type:%s\\(0x%04x\\)dst:%ssrc:%spayload:0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728$",
		ethType.String(),
		uint16(ethType),
		macDst,
		macSrc))
	got := psLog.CaptureLogOutput(func() {
		hdr := EthHdr{Dst: macDst, Src: macSrc, Type: ethType}
		payload := []byte{
			0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a,
			0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14,
			0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e,
			0x1f, 0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28,
		}
		psLog.I("", dumpFrame(&hdr, payload)...)
	})
	got = Trim(got)
	if !want.MatchString(got) {
		t.Errorf("dumpFrame() = %v; want %v", got, want)
	}
}

func TestReadFrame_1(t *testing.T) {
	ctrl, teardown := SetupReadFrameTest(t)
	defer teardown()

	m := psSyscall.NewMockISyscall(ctrl)
	m.EXPECT().
		Read(gomock.Any(), gomock.Any()).
		Do(func(_ int, buf []byte) {
			hdr := EthHdr{
				Dst:  EthAddr{11, 12, 13, 14, 15, 16},
				Src:  EthAddr{21, 22, 23, 24, 25, 26},
				Type: EthType(0x0800),
			}
			b := new(bytes.Buffer)
			_ = binary.Write(b, binary.BigEndian, &hdr)
			copy(buf, b.Bytes())
		}).
		Return(150, nil)
	psSyscall.Syscall = m

	dev := &Device{Addr_: EthAddr{11, 12, 13, 14, 15, 16}}

	_, got := ReadFrame(dev.Priv().FD, dev.Addr())
	if got != psErr.OK {
		t.Errorf("ReadFrame() = %v; want %v", got, psErr.OK)
	}
}

// Fail when Read() returns error.
func TestReadFrame_2(t *testing.T) {
	ctrl, teardown := SetupReadFrameTest(t)
	defer teardown()

	dev := &Device{Addr_: EthAddr{11, 12, 13, 14, 15, 16}}
	m := psSyscall.NewMockISyscall(ctrl)
	m.EXPECT().Read(gomock.Any(), gomock.Any()).Return(-1, errors.New(""))
	psSyscall.Syscall = m

	_, got := ReadFrame(dev.Priv().FD, dev.Addr())
	if got != psErr.Error {
		t.Errorf("ReadFrame() = %v; want %v", got, psErr.Error)
	}
}

// Fail when header length is invalid.
func TestReadFrame_3(t *testing.T) {
	ctrl, teardown := SetupReadFrameTest(t)
	defer teardown()

	m := psSyscall.NewMockISyscall(ctrl)
	m.EXPECT().Read(gomock.Any(), gomock.Any()).Return(10, nil)
	psSyscall.Syscall = m

	dev := &Device{Addr_: EthAddr{11, 12, 13, 14, 15, 16}}

	_, got := ReadFrame(dev.Priv().FD, dev.Addr())
	if got != psErr.Error {
		t.Errorf("ReadFrame() = %v; want %v", got, psErr.Error)
	}
}

// Success when no data exits.
func TestReadFrame_4(t *testing.T) {
	ctrl, teardown := SetupReadFrameTest(t)
	defer teardown()

	m := psSyscall.NewMockISyscall(ctrl)
	m.EXPECT().Read(gomock.Any(), gomock.Any()).Return(150, nil)
	psSyscall.Syscall = m

	dev := &Device{Addr_: EthAddr{33, 44, 55, 66, 77, 88}}

	_, got := ReadFrame(dev.Priv().FD, dev.Addr())
	if got != psErr.NoDataToRead {
		t.Errorf("ReadFrame() = %v; want %v", got, psErr.NoDataToRead)
	}
}

// Fail when Read() system call returns error.
func TestReadFrame_5(t *testing.T) {
	ctrl, teardown := SetupReadFrameTest(t)
	defer teardown()

	m := psSyscall.NewMockISyscall(ctrl)
	m.EXPECT().Read(gomock.Any(), gomock.Any()).Return(150, nil)
	psSyscall.Syscall = m

	dev := &Device{Addr_: EthAddr{33, 44, 55, 66, 77, 88}}

	_, got := ReadFrame(dev.Priv().FD, dev.Addr())
	if got != psErr.NoDataToRead {
		t.Errorf("ReadFrame() = %v; want %v", got, psErr.NoDataToRead)
	}
}

func SetupReadFrameTest(t *testing.T) (ctrl *gomock.Controller, teardown func()) {
	psLog.DisableOutput()
	ctrl = gomock.NewController(t)
	teardown = func() {
		ctrl.Finish()
		psLog.EnableOutput()
	}
	return
}

func Trim(s string) string {
	ret := strings.Replace(s, "", "", -1)
	ret = strings.Replace(ret, "\n", "", -1)
	ret = strings.Replace(ret, " ", "", -1)
	return ret
}
