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
		t.Errorf("EthType.String() = %s; want %s", got, want)
	}

	want = "ARP"
	got = EthType(0x0806).String()
	if got != want {
		t.Errorf("EthType.String() = %s; want %s", got, want)
	}

	want = "IPv6"
	got = EthType(0x86dd).String()
	if got != want {
		t.Errorf("EthType.String() = %s; want %s", got, want)
	}
}

func TestEthFrameDump(t *testing.T) {
	regexpDatetime := "[0-9]{4}/[0-9]{2}/[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2}"
	macDst := EthAddr{0x01, 0x02, 0x03, 0x04, 0x05, 0x06}
	macSrc := EthAddr{0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}
	ethType := EthType(0x0800)
	want, _ := regexp.Compile(fmt.Sprintf(
		"^.+ %s dst:           %s.+ %s src:           %s.+ %s type:          0x%04x \\(%s\\).+ %s payload \\(nbo\\): 01 02 03 04 05 06 07 08 09 0a 0b 0c 0d 0e 0f 10 11 12 13 14 .+ %s        15 16 17 18 19 1a 1b 1c 1d 1e 1f 20 21 22 23 24 25 26 27 28 $",
		regexpDatetime,
		macDst,
		regexpDatetime,
		macSrc,
		regexpDatetime,
		uint16(ethType),
		ethType.String(),
		regexpDatetime,
		regexpDatetime))
	got := psLog.CaptureLogOutput(func() {
		buf := new(bytes.Buffer)
		hdr := EthHdr{Dst: macDst, Src: macSrc, Type: ethType}
		payload := [...]byte{
			0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a,
			0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14,
			0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e,
			0x1f, 0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28,
		}
		_ = binary.Write(buf, binary.BigEndian, &hdr)
		_ = binary.Write(buf, binary.BigEndian, &payload)
		frame := buf.Bytes()
		EthFrameDump(frame[:EthHdrLen], frame[EthHdrLen:])
	})
	got = trim(got)
	if !want.MatchString(got) {
		t.Errorf("EthFrameDump() = %v; want %v", got, want)
	}
}

func TestReadFrame_1(t *testing.T) {
	psLog.DisableOutput()
	defer psLog.EnableOutput()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

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

	_, got := ReadEthFrame(dev.Priv().FD, dev.Addr())
	if got != psErr.OK {
		t.Errorf("ReadEthFrame() = %v; want %v", got, psErr.OK)
	}
}

// Fail when Read() returns error.
func TestReadFrame_2(t *testing.T) {
	psLog.DisableOutput()
	defer psLog.EnableOutput()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dev := &Device{Addr_: EthAddr{11, 12, 13, 14, 15, 16}}
	m := psSyscall.NewMockISyscall(ctrl)
	m.EXPECT().Read(gomock.Any(), gomock.Any()).Return(-1, errors.New(""))
	psSyscall.Syscall = m

	_, got := ReadEthFrame(dev.Priv().FD, dev.Addr())
	if got != psErr.Error {
		t.Errorf("ReadEthFrame() = %v; want %v", got, psErr.Error)
	}
}

// Fail when header length is invalid.
func TestReadFrame_3(t *testing.T) {
	psLog.DisableOutput()
	defer psLog.EnableOutput()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := psSyscall.NewMockISyscall(ctrl)
	m.EXPECT().Read(gomock.Any(), gomock.Any()).Return(10, nil)
	psSyscall.Syscall = m

	dev := &Device{Addr_: EthAddr{11, 12, 13, 14, 15, 16}}

	_, got := ReadEthFrame(dev.Priv().FD, dev.Addr())
	if got != psErr.Error {
		t.Errorf("ReadEthFrame() = %v; want %v", got, psErr.Error)
	}
}

// Success when no data exits.
func TestReadFrame_4(t *testing.T) {
	psLog.DisableOutput()
	defer psLog.EnableOutput()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := psSyscall.NewMockISyscall(ctrl)
	m.EXPECT().Read(gomock.Any(), gomock.Any()).Return(150, nil)
	psSyscall.Syscall = m

	dev := &Device{Addr_: EthAddr{33, 44, 55, 66, 77, 88}}

	_, got := ReadEthFrame(dev.Priv().FD, dev.Addr())
	if got != psErr.NoDataToRead {
		t.Errorf("ReadEthFrame() = %v; want %v", got, psErr.NoDataToRead)
	}
}
