package network

import (
	"bytes"
	"encoding/binary"
	"fmt"
	psErr "github.com/42milez/ProtocolStack/src/error"
	"github.com/42milez/ProtocolStack/src/ethernet"
	psLog "github.com/42milez/ProtocolStack/src/log"
	"github.com/42milez/ProtocolStack/src/timer"
	"sync"
	"time"
)

var ArpCondCh chan timer.Condition
var ArpSigCh chan timer.Signal

func ArpInputHandler(packet []byte, dev ethernet.IDevice) psErr.E {
	if len(packet) < ArpPacketLen {
		psLog.E(fmt.Sprintf("ARP packet length is too short: %d bytes", len(packet)))
		return psErr.InvalidPacket
	}

	buf := bytes.NewBuffer(packet)
	arpPacket := ArpPacket{}
	if err := binary.Read(buf, binary.BigEndian, &arpPacket); err != nil {
		return psErr.Error
	}

	if arpPacket.HT != ArpHwTypeEthernet || arpPacket.HAL != ethernet.EthAddrLen {
		psLog.E("Value of ARP packet header is invalid (Hardware)")
		return psErr.InvalidPacket
	}

	if arpPacket.PT != ethernet.EthTypeIpv4 || arpPacket.PAL != V4AddrLen {
		psLog.E("Value of ARP packet header is invalid (Protocol)")
		return psErr.InvalidPacket
	}

	psLog.I("Incoming ARP packet")
	dumpArpPacket(&arpPacket)

	iface := IfaceRepo.Lookup(dev, V4AddrFamily)
	if iface == nil {
		psLog.E(fmt.Sprintf("Interface for %s is not registered", dev.Name()))
		return psErr.InterfaceNotFound
	}

	if iface.Unicast.EqualV4(arpPacket.TPA) {
		if err := cache.Renew(arpPacket.SPA, arpPacket.SHA, ArpCacheStateResolved); err == psErr.NotFound {
			if err := cache.Add(arpPacket.SHA, arpPacket.SPA, ArpCacheStateResolved); err == psErr.Exist {
				psLog.W(fmt.Sprintf("Cache entry already exists: %s", arpPacket.SPA))
			}
		} else {
			psLog.I("ARP entry was renewed")
			psLog.I(fmt.Sprintf("\tspa: %s", arpPacket.SPA))
			psLog.I(fmt.Sprintf("\tsha: %s", arpPacket.SHA))
		}
		if arpPacket.Opcode == ArpOpRequest {
			if err := arpReply(arpPacket.SHA, arpPacket.SPA, iface); err != psErr.OK {
				return psErr.Error
			}
		}
	} else {
		psLog.I("ARP packet was ignored (It was sent to different address)")
	}

	return psErr.OK
}

func RunArpTimer(wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		ArpCondCh <-timer.Condition{
			CurrentState: timer.Running,
		}
		for {
			select {
			case signal := <-ArpSigCh:
				if signal == timer.Stop {
					return
				}
			default:
				arpTimer()
				time.Sleep(time.Second)
			}
		}
	}()
}

func StopArpTimer() {
	ArpSigCh <-timer.Stop
}

func arpReply(tha ethernet.EthAddr, tpa ArpProtoAddr, iface *Iface) psErr.E {
	packet := ArpPacket{
		ArpHdr: ArpHdr{
			HT:     ArpHwTypeEthernet,
			PT:     ethernet.EthTypeIpv4,
			HAL:    ethernet.EthAddrLen,
			PAL:    V4AddrLen,
			Opcode: ArpOpReply,
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
		return psErr.Error
	}

	if err := iface.Dev.Transmit(tha, buf.Bytes(), ethernet.EthTypeArp); err != psErr.OK {
		return psErr.Error
	}

	return psErr.OK
}

func arpRequest(iface *Iface, ip IP) psErr.E {
	packet := ArpPacket{
		ArpHdr: ArpHdr{
			HT:     ArpHwTypeEthernet,
			PT:     ethernet.EthTypeIpv4,
			HAL:    ethernet.EthAddrLen,
			PAL:    V4AddrLen,
			Opcode: ArpOpRequest,
		},
		SHA: iface.Dev.Addr(),
		SPA: iface.Unicast.ToV4(),
		THA: ethernet.EthAddr{},
		TPA: ip.ToV4(),
	}

	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, &packet); err != nil {
		return psErr.Error
	}
	payload := buf.Bytes()

	psLog.I("Outgoing ARP packet")
	dumpArpPacket(&packet)

	if err := Transmit(ethernet.EthAddrBroadcast, payload, ethernet.EthTypeArp, iface); err != psErr.OK {
		return psErr.Error
	}

	return psErr.OK
}

func arpResolve(iface *Iface, ip IP) (ethernet.EthAddr, ArpStatus) {
	if iface.Dev.Type() != ethernet.DevTypeEthernet {
		psLog.E(fmt.Sprintf("Unsupported device type: %s", iface.Dev.Type()))
		return ethernet.EthAddr{}, ArpStatusError
	}

	if iface.Family != V4AddrFamily {
		psLog.E(fmt.Sprintf("Unsupported address family: %s", iface.Family))
		return ethernet.EthAddr{}, ArpStatusError
	}

	ethAddr, found := cache.EthAddr(ip.ToV4())
	if !found {
		if err := cache.Add(ethernet.EthAddr{}, ip.ToV4(), ArpCacheStateIncomplete); err != psErr.OK {
			return ethernet.EthAddr{}, ArpStatusError
		}
		if err := arpRequest(iface, ip); err != psErr.OK {
			return ethernet.EthAddr{}, ArpStatusError
		}
		return ethernet.EthAddr{}, ArpStatusIncomplete
	}

	return ethAddr, ArpStatusComplete
}

func arpTimer() {
	ret := cache.expire()
	if len(ret) == 0 {
		return
	}
	psLog.I("ARP cache entries were expired:")
	for i, v := range ret {
		psLog.I(fmt.Sprintf("\t%d: %s", i+1, v))
	}
}

func dumpArpPacket(packet *ArpPacket) {
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

func init() {
	ArpCondCh = make(chan timer.Condition)
	ArpSigCh = make(chan timer.Signal)
}
