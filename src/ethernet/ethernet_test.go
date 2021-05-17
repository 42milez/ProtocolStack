package ethernet

import (
	"fmt"
	"github.com/42milez/ProtocolStack/src/log"
	"regexp"
	"strings"
	"testing"
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

func TestEthAddr_EqualByte(t *testing.T) {
	ethAddr1 := EthAddr([EthAddrLen]byte{11, 22, 33, 44, 55, 66})
	b1 := []byte{11, 22, 33, 44, 55, 66}
	var b2 []byte

	got := ethAddr1.EqualByte(b1)
	if got != true {
		t.Errorf("EthAddr.Equal() = %v; want %v", got, true)
	}

	got = ethAddr1.EqualByte(b2)
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
	got := log.CaptureLogOutput(func() {
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

func TestReadFrame(t *testing.T) {

}
