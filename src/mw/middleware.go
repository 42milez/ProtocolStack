package mw

import (
	"math/rand"
	"time"
)

const (
	EtARP  EthType = 0x0806
	EtIPV4 EthType = 0x0800
	EtIPV6 EthType = 0x86dd
)
const (
	PnICMP ProtocolNumber = 1
	PnTCP  ProtocolNumber = 6
	PnUDP  ProtocolNumber = 17
)
const maxUint8 = ^uint8(0)
const maxUint16 = ^uint16(0)
const xChBufSize = 10

var EthRxCh chan *EthMessage // channel for receiving packets
var EthTxCh chan *EthMessage // channel for sending packets
var ArpRxCh chan *ArpRxMessage
var ArpTxCh chan *ArpTxMessage
var IpRxCh chan *EthMessage
var IpTxCh chan *IpMessage
var IcmpDeadLetterQueue chan *IcmpQueueEntry
var IcmpRxCh chan *IcmpRxMessage
var IcmpTxCh chan *IcmpTxMessage
var TcpRxCh chan *TcpRxMessage
var TcpTxCh chan *TcpTxMessage

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
	Dst      [V4AddrLen]byte
	Src      [V4AddrLen]byte
}

type IcmpQueueEntry struct {
	Packet []byte
}

type IcmpRxMessage struct {
	Packet []byte
	Dst    [V4AddrLen]byte
	Src    [V4AddrLen]byte
	Dev    IDevice
}

type IcmpTxMessage struct {
	Type    uint8
	Code    uint8
	Content uint32
	Data    []byte
	Src     IP
	Dst     IP
}

type TcpRxMessage struct {
	ProtoNum   uint8
	RawSegment []byte
	Dst        [V4AddrLen]byte
	Src        [V4AddrLen]byte
	Iface      *Iface
}

type TcpTxMessage struct{}

func RandU8() uint8 {
	return uint8(rand.Intn(int(maxUint8) + 1))
}

func RandU16() uint16 {
	return uint16(rand.Intn(int(maxUint16) + 1))
}

func RandU32() uint32 {
	return rand.Uint32()
}

func init() {
	EthRxCh = make(chan *EthMessage, xChBufSize)
	EthTxCh = make(chan *EthMessage, xChBufSize)
	ArpRxCh = make(chan *ArpRxMessage, xChBufSize)
	ArpTxCh = make(chan *ArpTxMessage, xChBufSize)
	IpRxCh = make(chan *EthMessage, xChBufSize)
	IpTxCh = make(chan *IpMessage, xChBufSize)
	IcmpDeadLetterQueue = make(chan *IcmpQueueEntry, xChBufSize)
	IcmpRxCh = make(chan *IcmpRxMessage, xChBufSize)
	IcmpTxCh = make(chan *IcmpTxMessage, xChBufSize)
	TcpRxCh = make(chan *TcpRxMessage, xChBufSize)
	TcpTxCh = make(chan *TcpTxMessage, xChBufSize)

	rand.Seed(time.Now().UnixNano())
}
