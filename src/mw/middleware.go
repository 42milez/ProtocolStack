package mw

import (
	"fmt"
	psErr "github.com/42milez/ProtocolStack/src/error"
	psLog "github.com/42milez/ProtocolStack/src/log"
)

const ethRxChBufSize = 10
const ethTxChBufSize = 10
const arpRxChBufSize = 10
const arpTxChBufSize = 10
const ipRxChBufSize = 10
const ipTxChBufSize = 10
const icmpRxChBufSize = 10
const icmpTxChBufSize = 10

var EthRxCh chan *EthMessage // channel for receiving packets
var EthTxCh chan *EthMessage // channel for sending packets

var ArpRxCh chan *EthMessage
var ArpTxCh chan *EthMessage
var IpRxCh chan *EthMessage
var IpTxCh chan *IpMessage
var IcmpRxCh chan *IcmpRxMessage
var IcmpTxCh chan *IcmpTxMessage

// Ethertypes
// https://www.iana.org/assignments/ieee-802-numbers/ieee-802-numbers.xhtml#ieee-802-numbers-1

type EthType uint16

func (v EthType) String() string {
	return ethTypes[v]
}

type EthMessage struct {
	Type    EthType
	Content []byte
	Dev     IDevice
}

// Computing the Internet Checksum
// https://datatracker.ietf.org/doc/html/rfc1071

func Checksum(b []byte) uint16 {
	var sum uint32
	// sum up all fields of IP header by each 16bits (except Header Checksum and Options)
	for i := 0; i < len(b); i += 2 {
		sum += uint32(uint16(b[i])<<8 | uint16(b[i+1]))
	}
	//
	sum = ((sum & 0xffff0000) >> 16) + (sum & 0x0000ffff)
	return ^(uint16(sum))
}

func Rx(packet *EthMessage) psErr.E {
	switch packet.Type {
	case ARP:
		ArpRxCh <- packet
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

func init() {
	EthRxCh = make(chan *EthMessage, ethRxChBufSize)
	EthTxCh = make(chan *EthMessage, ethTxChBufSize)
	ArpRxCh = make(chan *EthMessage, arpRxChBufSize)
	ArpTxCh = make(chan *EthMessage, arpTxChBufSize)
	IpRxCh = make(chan *EthMessage, ipRxChBufSize)
	IpTxCh = make(chan *IpMessage, ipTxChBufSize)
	IcmpRxCh = make(chan *IcmpRxMessage, icmpRxChBufSize)
	IcmpTxCh = make(chan *IcmpTxMessage, icmpTxChBufSize)
}

var ethTypes = map[EthType]string{
	0x0800: "IPv4",
	0x0806: "ARP",
	0x86dd: "IPv6",
}
