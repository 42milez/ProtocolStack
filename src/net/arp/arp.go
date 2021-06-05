package arp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	psErr "github.com/42milez/ProtocolStack/src/error"
	psLog "github.com/42milez/ProtocolStack/src/log"
	"github.com/42milez/ProtocolStack/src/mw"
	"github.com/42milez/ProtocolStack/src/net"
	"github.com/42milez/ProtocolStack/src/repo"
	"github.com/42milez/ProtocolStack/src/worker"
	"sync"
	"time"
)

const Ethernet HwType = 0x0001
const (
	Complete Status = iota
	Incomplete
	Error
)
const xChBufSize = 5

var RcvRxCh chan *worker.Message
var RcvTxCh chan *worker.Message
var SndRxCh chan *worker.Message
var SndTxCh chan *worker.Message
var TimerRxCh chan *worker.Message
var TimerTxCh chan *worker.Message

// Hardware Types
// https://www.iana.org/assignments/arp-parameters/arp-parameters.xhtml#arp-parameters-2

var arpHwTypes = map[HwType]string{
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

var arpOpCodes = map[Opcode]string{
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
	return arpHwTypes[v]
}

type Opcode uint16

func (v Opcode) String() string {
	return arpOpCodes[v]
}

type Packet struct {
	Hdr
	SHA mw.Addr      // sender hardware address
	SPA ArpProtoAddr // sender protocol address
	THA mw.Addr      // target hardware address
	TPA ArpProtoAddr // target protocol address
}

type Status int

func Receive(packet []byte, dev mw.IDevice) psErr.E {
	if len(packet) < PacketLen {
		psLog.E(fmt.Sprintf("ARP packet length is too short: %d bytes", len(packet)))
		return psErr.InvalidPacket
	}

	buf := bytes.NewBuffer(packet)
	arpPacket := Packet{}
	if err := binary.Read(buf, binary.BigEndian, &arpPacket); err != nil {
		return psErr.ReadFromBufError
	}

	if arpPacket.HT != Ethernet || arpPacket.HAL != mw.AddrLen {
		psLog.E("Value of ARP packet header is invalid (Hardware)")
		return psErr.InvalidPacket
	}

	if arpPacket.PT != mw.IPv4 || arpPacket.PAL != mw.V4AddrLen {
		psLog.E("Value of ARP packet header is invalid (Protocol)")
		return psErr.InvalidPacket
	}

	psLog.I("Incoming ARP packet")
	dumpArpPacket(&arpPacket)

	iface := repo.IfaceRepo.Lookup(dev, mw.V4AddrFamily)
	if iface == nil {
		psLog.E(fmt.Sprintf("Interface for %s is not registered", dev.Name()))
		return psErr.InterfaceNotFound
	}

	if iface.Unicast.EqualV4(arpPacket.TPA) {
		if err := cache.Renew(arpPacket.SPA, arpPacket.SHA, resolved); err == psErr.NotFound {
			_ = cache.Create(arpPacket.SHA, arpPacket.SPA, resolved)
		} else {
			psLog.I("ARP entry was renewed")
			psLog.I(fmt.Sprintf("\tspa: %s", arpPacket.SPA))
			psLog.I(fmt.Sprintf("\tsha: %s", arpPacket.SHA))
		}
		if arpPacket.Opcode == Request {
			if err := SendReply(arpPacket.SHA, arpPacket.SPA, iface); err != psErr.OK {
				return psErr.Error
			}
		}
	} else {
		psLog.I("ARP packet was ignored (It was sent to different address)")
	}

	return psErr.OK
}

func SendReply(tha mw.Addr, tpa ArpProtoAddr, iface *mw.Iface) psErr.E {
	packet := Packet{
		Hdr: Hdr{
			HT:     Ethernet,
			PT:     mw.IPv4,
			HAL:    mw.AddrLen,
			PAL:    mw.V4AddrLen,
			Opcode: Reply,
		},
		THA: tha,
		TPA: tpa,
	}
	addr := iface.Dev.Addr()
	copy(packet.SHA[:], addr[:])
	copy(packet.SPA[:], iface.Unicast[:])

	psLog.I("Outgoing ARP packet (REPLY):")
	dumpArpPacket(&packet)

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
			HAL:    mw.AddrLen,
			PAL:    mw.V4AddrLen,
			Opcode: Request,
		},
		SHA: iface.Dev.Addr(),
		SPA: iface.Unicast.ToV4(),
		THA: mw.Addr{},
		TPA: ip.ToV4(),
	}

	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, &packet); err != nil {
		return psErr.WriteToBufError
	}
	payload := buf.Bytes()

	psLog.I("Outgoing ARP packet")
	dumpArpPacket(&packet)

	if err := net.Transmit(mw.Broadcast, payload, mw.ARP, iface); err != psErr.OK {
		return psErr.Error
	}

	return psErr.OK
}

func Resolve(iface *mw.Iface, ip mw.IP) (mw.Addr, Status) {
	if iface.Dev.Type() != mw.DevTypeEthernet {
		psLog.E(fmt.Sprintf("Unsupported device type: %s", iface.Dev.Type()))
		return mw.Addr{}, Error
	}

	if iface.Family != mw.V4AddrFamily {
		psLog.E(fmt.Sprintf("Unsupported address family: %s", iface.Family))
		return mw.Addr{}, Error
	}

	entry := cache.GetEntry(ip.ToV4())
	if entry == nil {
		if err := cache.Create(mw.Addr{}, ip.ToV4(), incomplete); err != psErr.OK {
			return mw.Addr{}, Error
		}
		if err := SendRequest(iface, ip); err != psErr.OK {
			return mw.Addr{}, Error
		}
		return mw.Addr{}, Incomplete
	}

	return entry.HA, Complete
}

func StartService(wg *sync.WaitGroup) {
	wg.Add(3)
	go receiver(wg)
	go sender(wg)
	go timer(wg)
}

func StopService() {
	msg := &worker.Message{
		Desired: worker.Stopped,
	}
	RcvRxCh <- msg
	SndRxCh <- msg
	TimerRxCh <- msg
}

func dumpArpPacket(packet *Packet) {
	psLog.I(fmt.Sprintf("\thardware type:           %s", packet.HT))
	psLog.I(fmt.Sprintf("\tprotocol Type:           %s", packet.PT))
	psLog.I(fmt.Sprintf("\thardware address length: %d", packet.HAL))
	psLog.I(fmt.Sprintf("\tprotocol address length: %d", packet.PAL))
	psLog.I(fmt.Sprintf("\topcode:                  %s (%d)", packet.Opcode, uint16(packet.Opcode)))
	psLog.I(fmt.Sprintf("\tsender hardware address: %s", packet.SHA))
	psLog.I(fmt.Sprintf("\tsender protocol address: %v", packet.SPA))
	psLog.I(fmt.Sprintf("\ttarget hardware address: %s", packet.THA))
	psLog.I(fmt.Sprintf("\ttarget protocol address: %v", packet.TPA))
}

func receiver(wg *sync.WaitGroup) {
	defer wg.Done()

	RcvTxCh <- &worker.Message{
		Current: worker.Running,
	}

	for {
		select {
		case msg := <-RcvRxCh:
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

	SndTxCh <- &worker.Message{
		Current: worker.Running,
	}

	for {
		msg := <-SndRxCh
		if msg.Desired == worker.Stopped {
			return
		}
	}
}

func timer(wg *sync.WaitGroup) {
	defer wg.Done()

	TimerTxCh <- &worker.Message{
		Current: worker.Running,
	}

	for {
		select {
		case msg := <-TimerRxCh:
			if msg.Desired == worker.Stopped {
				return
			}
		default:
			ret := cache.Expire()
			if len(ret) != 0 {
				psLog.I("ARP cache entries were expired:")
				for i, v := range ret {
					psLog.I(fmt.Sprintf("\t%d: %s", i+1, v))
				}
			}
		}
		time.Sleep(time.Second)
	}
}

func init() {
	RcvRxCh = make(chan *worker.Message, xChBufSize)
	RcvTxCh = make(chan *worker.Message, xChBufSize)
	SndRxCh = make(chan *worker.Message, xChBufSize)
	SndTxCh = make(chan *worker.Message, xChBufSize)
	TimerRxCh = make(chan *worker.Message, xChBufSize)
	TimerTxCh = make(chan *worker.Message, xChBufSize)
}
