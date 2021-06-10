package tcp

import "github.com/42milez/ProtocolStack/src/mw"

var PcbRepo []*PCB

type EndPoint struct {
	Addr mw.V4Addr
	Port uint16
}

type PCB struct {
	Local EndPoint
	Foreign EndPoint
	SND struct {
		UNA uint32
		NXT uint32
		WND uint16
		UP uint16
		WL1 uint32
		WL2 uint32
	}
	ISS uint32
	RCV struct {
		NXT uint32
		WND uint16
		UP uint16
	}
	IRS uint32
	MTU uint16
	MSS uint16
}
