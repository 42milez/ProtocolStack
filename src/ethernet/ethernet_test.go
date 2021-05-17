package ethernet

import (
	"testing"
)

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
