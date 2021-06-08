package arp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	psErr "github.com/42milez/ProtocolStack/src/error"
	psLog "github.com/42milez/ProtocolStack/src/log"
	"github.com/42milez/ProtocolStack/src/monitor"
	"github.com/42milez/ProtocolStack/src/mw"
	"github.com/42milez/ProtocolStack/src/net"
	"github.com/42milez/ProtocolStack/src/repo"
	"github.com/42milez/ProtocolStack/src/worker"
	"sync"
	"time"
)

const PacketLen = 28 // byte
const Ethernet HwType = 0x0001
const Request Opcode = 0x0001
const Reply Opcode = 0x0002
const (
	Complete Status = iota
	Incomplete
	Error
)
const xChBufSize = 5

var rcvMonCh chan *worker.Message
var rcvSigCh chan *worker.Message
var sndMonCh chan *worker.Message
var sndSigCh chan *worker.Message
var tmrMonCh chan *worker.Message
var tmrSigCh chan *worker.Message

var receiverID uint32
var senderID uint32
var timerID uint32

// Hardware Types
// https://www.iana.org/assignments/arp-parameters/arp-parameters.xhtml#arp-parameters-2

var hwTypes = map[HwType]string{
	// 0: Reserved
	1:  "Ethernet (10Mb)",
	2:  "Experimental Ethernet (3Mb)",
	3:  "Amateur Radio AX.25",
	4:  "Proteon ProNET Token Ring",
	5:  "Chaos",
	6:  "IEEE 802 Networks",
	7:  "ARCNET",
	8:  "Hyperchannel",
	9:  "Lanstar",
	10: "Autonet Short Address",
	11: "LocalTalk",
	12: "LocalNet (IBM PCNet or SYTEK LocalNET)",
	13: "Ultra link",
	14: "SMDS",
	15: "Frame Relay",
	16: "Asynchronous Transmission Mode (ATM)",
	17: "HDLC",
	18: "Fibre Channel",
	19: "Asynchronous Transmission Mode (ATM)",
	20: "Serial Line",
	21: "Asynchronous Transmission Mode (ATM)",
	22: "MIL-STD-188-220",
	23: "Metricom",
	24: "IEEE 1394.1995",
	25: "MAPOS",
	26: "Twinaxial",
	27: "EUI-64",
	28: "HIPARP",
	29: "IP and ARP over ISO 7816-3",
	30: "ARPSec",
	31: "IPsec tunnel",
	32: "InfiniBand (TM)",
	33: "TIA-102 Project 25 Common Air Interface (CAI)",
	34: "Wiegand Interface",
	35: "Pure IP",
	36: "HW_EXP1",
	37: "HFI",
	// 38-255: Unassigned
	256: "HW_EXP2",
	257: "AEthernet",
	// 258-65534: Unassigned
	// 65535: Reserved
}

// Operation Codes
// https://www.iana.org/assignments/arp-parameters/arp-parameters.xhtml#arp-parameters-1

var opCodes = map[Opcode]string{
	// 0: Reserved
	1:  "REQUEST",
	2:  "REPLY",
	3:  "request Reverse",
	4:  "reply Reverse",
	5:  "DRARP-SendRequest",
	6:  "DRARP-SendReply",
	7:  "DRARP-Error",
	8:  "InARP-SendRequest",
	9:  "InARP-SendReply",
	10: "ARP-NAK",
	11: "MARS-SendRequest",
	12: "MARS-Multi",
	13: "MARS-MServ",
	14: "MARS-Join",
	15: "MARS-Leave",
	16: "MARS-NAK",
	17: "MARS-Unserv",
	18: "MARS-SJoin",
	19: "MARS-SLeave",
	20: "MARS-Grouplist-SendRequest",
	21: "MARS-Grouplist-SendReply",
	22: "MARS-Redirect-Map",
	23: "MAPOS-UNARP",
	24: "OP_EXP1",
	25: "OP_EXP2",
	// 26-65534: Unassigned
	// 65535: Reserved
}

// Address Resolution Protocol (ARP) Parameters
// https://www.iana.org/assignments/arp-parameters/arp-parameters.xhtml

// EtherType
// https://en.wikipedia.org/wiki/EtherType#Examples

// Notes: Protocol Type is same as EtherType.

type Hdr struct {
	HT     HwType     // hardware type
	PT     mw.EthType // protocol type
	HAL    uint8      // hardware address length
	PAL    uint8      // protocol address length
	Opcode Opcode
}

type HwType uint16

func (v HwType) String() string {
	return hwTypes[v]
}

type Opcode uint16

func (v Opcode) String() string {
	return opCodes[v]
}

type Packet struct {
	Hdr
	SHA mw.EthAddr // sender hardware address
	SPA mw.V4Addr  // sender protocol address
	THA mw.EthAddr // target hardware address
	TPA mw.V4Addr  // target protocol address
}

type Status int

func Receive(packet []byte, dev mw.IDevice) psErr.E {
	if len(packet) < PacketLen {
		psLog.E(fmt.Sprintf("arp packet length is too short: %d bytes", len(packet)))
		return psErr.InvalidPacket
	}

	buf := bytes.NewBuffer(packet)
	arpPacket := Packet{}
	if err := binary.Read(buf, binary.BigEndian, &arpPacket); err != nil {
		return psErr.ReadFromBufError
	}

	if arpPacket.HT != Ethernet || arpPacket.HAL != mw.EthAddrLen {
		psLog.E("invalid arp header field (Hardware)")
		return psErr.InvalidPacket
	}

	if arpPacket.PT != mw.IPv4 || arpPacket.PAL != mw.V4AddrLen {
		psLog.E("invalid arp header field (Protocol)")
		return psErr.InvalidPacket
	}

	psLog.D("incoming arp packet", dump(&arpPacket)...)

	iface := repo.IfaceRepo.Lookup(dev, mw.V4AddrFamily)
	if iface == nil {
		psLog.E(fmt.Sprintf("Interface for %s is not registered", dev.Name()))
		return psErr.InterfaceNotFound
	}

	if iface.Unicast.EqualV4(arpPacket.TPA) {
		if err := cache.Renew(arpPacket.SPA, arpPacket.SHA, resolved); err == psErr.NotFound {
			_ = cache.Create(arpPacket.SHA, arpPacket.SPA, resolved)
		} else {
			psLog.I("arp cache entry was renewed",
				fmt.Sprintf("spa: %s", arpPacket.SPA),
				fmt.Sprintf("sha: %s", arpPacket.SHA))
		}
		if arpPacket.Opcode == Request {
			if err := SendReply(arpPacket.SHA, arpPacket.SPA, iface); err != psErr.OK {
				return psErr.Error
			}
		}
	} else {
		psLog.I("arp packet was ignored (it was sent to different address)")
	}

	return psErr.OK
}

func SendReply(tha mw.EthAddr, tpa mw.V4Addr, iface *mw.Iface) psErr.E {
	packet := Packet{
		Hdr: Hdr{
			HT:     Ethernet,
			PT:     mw.IPv4,
			HAL:    mw.EthAddrLen,
			PAL:    mw.V4AddrLen,
			Opcode: Reply,
		},
		THA: tha,
		TPA: tpa,
	}
	addr := iface.Dev.Addr()
	copy(packet.SHA[:], addr[:])
	copy(packet.SPA[:], iface.Unicast[:])

	psLog.D("outgoing arp packet", dump(&packet)...)

	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, &packet); err != nil {
		return psErr.WriteToBufError
	}

	if err := iface.Dev.Transmit(tha, buf.Bytes(), mw.ARP); err != psErr.OK {
		return psErr.Error
	}

	return psErr.OK
}

func SendRequest(iface *mw.Iface, ip mw.IP) psErr.E {
	packet := Packet{
		Hdr: Hdr{
			HT:     Ethernet,
			PT:     mw.IPv4,
			HAL:    mw.EthAddrLen,
			PAL:    mw.V4AddrLen,
			Opcode: Request,
		},
		SHA: iface.Dev.Addr(),
		SPA: iface.Unicast.ToV4(),
		THA: mw.EthAddr{},
		TPA: ip.ToV4(),
	}

	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, &packet); err != nil {
		return psErr.WriteToBufError
	}
	payload := buf.Bytes()

	psLog.D("outgoing arp packet", dump(&packet)...)

	if err := net.Transmit(mw.EthBroadcast, payload, mw.ARP, iface); err != psErr.OK {
		return psErr.Error
	}

	return psErr.OK
}

func Resolve(iface *mw.Iface, ip mw.IP) (mw.EthAddr, Status) {
	if iface.Dev.Type() != mw.EthernetDevice {
		psLog.E(fmt.Sprintf("unsupported device type: %s", iface.Dev.Type()))
		return mw.EthAddr{}, Error
	}

	if iface.Family != mw.V4AddrFamily {
		psLog.E(fmt.Sprintf("unsupported address family: %s", iface.Family))
		return mw.EthAddr{}, Error
	}

	entry := cache.GetEntry(ip.ToV4())
	if entry == nil {
		if err := cache.Create(mw.EthAddr{}, ip.ToV4(), incomplete); err != psErr.OK {
			return mw.EthAddr{}, Error
		}
		if err := SendRequest(iface, ip); err != psErr.OK {
			return mw.EthAddr{}, Error
		}
		return mw.EthAddr{}, Incomplete
	}

	return entry.HA, Complete
}

func Start(wg *sync.WaitGroup) psErr.E {
	wg.Add(3)
	go receiver(wg)
	go sender(wg)
	go timer(wg)
	psLog.D("arp service started")
	return psErr.OK
}

func Stop() {
	msg := &worker.Message{
		Desired: worker.Stopped,
	}
	rcvSigCh <- msg
	sndSigCh <- msg
	tmrSigCh <- msg
}

func dump(packet *Packet) (ret []string) {
	ret = append(ret, fmt.Sprintf("hardware type:           %s", packet.HT))
	ret = append(ret, fmt.Sprintf("protocol Type:           %s", packet.PT))
	ret = append(ret, fmt.Sprintf("hardware address length: %d", packet.HAL))
	ret = append(ret, fmt.Sprintf("protocol address length: %d", packet.PAL))
	ret = append(ret, fmt.Sprintf("opcode:                  %s (%d)", packet.Opcode, uint16(packet.Opcode)))
	ret = append(ret, fmt.Sprintf("sender hardware address: %s", packet.SHA))
	ret = append(ret, fmt.Sprintf("sender protocol address: %v", packet.SPA))
	ret = append(ret, fmt.Sprintf("target hardware address: %s", packet.THA))
	ret = append(ret, fmt.Sprintf("target protocol address: %v", packet.TPA))
	return
}

func receiver(wg *sync.WaitGroup) {
	defer wg.Done()

	rcvMonCh <- &worker.Message{
		ID:      receiverID,
		Current: worker.Running,
	}

	for {
		select {
		case msg := <-rcvSigCh:
			if msg.Desired == worker.Stopped {
				return
			}
		case msg := <-mw.ArpRxCh:
			if err := Receive(msg.Packet, msg.Dev); err != psErr.OK {
				return
			}
		}
	}
}

func sender(wg *sync.WaitGroup) {
	defer wg.Done()

	sndMonCh <- &worker.Message{
		ID:      senderID,
		Current: worker.Running,
	}

	for {
		msg := <-sndSigCh
		if msg.Desired == worker.Stopped {
			return
		}
	}
}

func timer(wg *sync.WaitGroup) {
	defer wg.Done()

	tmrMonCh <- &worker.Message{
		ID:      timerID,
		Current: worker.Running,
	}

	for {
		select {
		case msg := <-tmrSigCh:
			if msg.Desired == worker.Stopped {
				return
			}
		default:
			ret := cache.Expire()
			if len(ret) != 0 {
				psLog.I("arp cache entries were expired")
				for i, v := range ret {
					psLog.I(fmt.Sprintf("%d: %s", i+1, v))
				}
			}
		}
		time.Sleep(time.Second)
	}
}

func init() {
	rcvMonCh = make(chan *worker.Message, xChBufSize)
	rcvSigCh = make(chan *worker.Message, xChBufSize)
	receiverID = monitor.Register("ARP Receiver", rcvMonCh, rcvSigCh)

	sndMonCh = make(chan *worker.Message, xChBufSize)
	sndSigCh = make(chan *worker.Message, xChBufSize)
	senderID = monitor.Register("ARP Sender", sndMonCh, sndSigCh)

	tmrMonCh = make(chan *worker.Message, xChBufSize)
	tmrSigCh = make(chan *worker.Message, xChBufSize)
	timerID = monitor.Register("ARP Timer", tmrMonCh, tmrSigCh)
}
