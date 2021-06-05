package mw

import (
	"fmt"
	psErr "github.com/42milez/ProtocolStack/src/error"
	psLog "github.com/42milez/ProtocolStack/src/log"
)

const xChBufSize = 10

var EthRxCh chan *EthMessage // channel for receiving packets
var EthTxCh chan *EthMessage // channel for sending packets
var ArpRxCh chan *ArpRxMessage
var ArpTxCh chan *ArpTxMessage
var IpRxCh chan *EthMessage
var IpTxCh chan *IpMessage
var IcmpRxCh chan *IcmpRxMessage
var IcmpTxCh chan *IcmpTxMessage

// Ethertypes
// https://www.iana.org/assignments/ieee-802-numbers/ieee-802-numbers.xhtml#ieee-802-numbers-1

var ethTypes = map[EthType]string{
	0x0800: "IPv4",
	0x0806: "ARP",
	0x86dd: "IPv6",
}

type EthType uint16

func (v EthType) String() string {
	return ethTypes[v]
}

type EthMessage struct {
	Type    EthType
	Content []byte
	Dev     IDevice
}

type ArpRxMessage struct {
	Packet []byte
	Dev    IDevice
}

type ArpTxMessage struct {
	Iface *Iface
	IP    IP
}

type IpMessage struct {
	ProtoNum ProtocolNumber
	Packet   []byte
	Dst      IP
	Src      IP
}

type IcmpRxMessage struct {
	Payload []byte
	Dst     [V4AddrLen]byte
	Src     [V4AddrLen]byte
	Dev     IDevice
}

type IcmpTxMessage struct {
	Type    uint8
	Code    uint8
	Content uint32
	Payload []byte
	Src     IP
	Dst     IP
}

func Rx(packet *EthMessage) psErr.E {
	switch packet.Type {
	case ARP:
		ArpRxCh <- &ArpRxMessage{
			Packet: packet.Content,
			Dev:    packet.Dev,
		}
	case IPv4:
		IpRxCh <- packet
	default:
		psLog.E(fmt.Sprintf("Unknown ether type: 0x%04x", uint16(packet.Type)))
		return psErr.Error
	}
	return psErr.OK
}

func Tx(msg *EthMessage) psErr.E {
	return psErr.OK
}

func init() {
	EthRxCh = make(chan *EthMessage, xChBufSize)
	EthTxCh = make(chan *EthMessage, xChBufSize)
	ArpRxCh = make(chan *ArpRxMessage, xChBufSize)
	ArpTxCh = make(chan *ArpTxMessage, xChBufSize)
	IpRxCh = make(chan *EthMessage, xChBufSize)
	IpTxCh = make(chan *IpMessage, xChBufSize)
	IcmpRxCh = make(chan *IcmpRxMessage, xChBufSize)
	IcmpTxCh = make(chan *IcmpTxMessage, xChBufSize)
}
