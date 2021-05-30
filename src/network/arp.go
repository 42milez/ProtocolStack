package network

import (
	"bytes"
	"encoding/binary"
	"fmt"
	psErr "github.com/42milez/ProtocolStack/src/error"
	"github.com/42milez/ProtocolStack/src/ethernet"
	psLog "github.com/42milez/ProtocolStack/src/log"
)

var cache *ArpCache

func ArpInputHandler(payload []byte, dev ethernet.IDevice) psErr.E {
	if len(payload) < ArpPacketSize {
		psLog.E(fmt.Sprintf("ARP packet length is too short: %d bytes", len(payload)))
		return psErr.InvalidPacket
	}

	buf := bytes.NewBuffer(payload)
	packet := ArpPacket{}
	if err := binary.Read(buf, binary.BigEndian, &packet); err != nil {
		psLog.E(fmt.Sprintf("binary.Read() failed: %s", err))
		return psErr.Error
	}

	if packet.HT != ArpHwTypeEthernet || packet.HAL != ethernet.EthAddrLen {
		psLog.E("Value of ARP packet header is invalid (Hardware)")
		return psErr.InvalidPacket
	}

	if packet.PT != ethernet.EthTypeIpv4 || packet.PAL != V4AddrLen {
		psLog.E("Value of ARP packet header is invalid (Protocol)")
		return psErr.InvalidPacket
	}

	psLog.I("Incoming ARP packet")
	arpHdrDump(&packet)

	iface := IfaceRepo.Get(dev, V4AddrFamily)
	if iface == nil {
		psLog.E(fmt.Sprintf("Interface for %s is not registered", dev.DevName()))
		return psErr.InterfaceNotFound
	}

	if iface.Unicast.EqualV4(packet.TPA) {
		if err := cache.Update(&packet); err == psErr.NotFound {
			if err := cache.Add(&packet); err != psErr.OK {
				psLog.E(fmt.Sprintf("ArpCache.Add() failed: %s", err))
			}
		} else {
			psLog.I("ARP entry was updated")
			psLog.I(fmt.Sprintf("\tSPA: %v", packet.SPA.String()))
			psLog.I(fmt.Sprintf("\tSHA: %v", packet.SHA.String()))
		}
		if packet.Opcode == ArpOpRequest {
			if err := arpReply(packet.SHA, packet.SPA, iface); err != psErr.OK {
				psLog.E(fmt.Sprintf("arpReply() failed: %s", err))
				return psErr.Error
			}
		}
	} else {
		psLog.I("Ignored ARP packet (It was sent to different address)")
	}

	return psErr.OK
}

func arpHdrDump(packet *ArpPacket) {
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
	addr := iface.Dev.EthAddr()
	copy(packet.SHA[:], addr[:])
	copy(packet.SPA[:], iface.Unicast[:])

	psLog.I("Outgoing ARP packet (REPLY):")
	arpHdrDump(&packet)

	psLog.I(fmt.Sprintf("ARP packet (REPLY) will be sent from %s", iface.Unicast))

	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, &packet); err != nil {
		psLog.E(fmt.Sprintf("binary.Write() failed: %s", err))
		return psErr.Error
	}

	if err := iface.Dev.Transmit(tha, buf.Bytes(), ethernet.EthTypeArp); err != psErr.OK {
		psLog.E(fmt.Sprintf("IDevice.Transmit() failed: %s", err))
		return psErr.Error
	}

	return psErr.OK
}

// TODO:
//func arpRequest() {}

// TODO:
//func arpResolve() {}

// TODO:
//func arpTimer() {}

func init() {
	cache = &ArpCache{}
	cache.Init()
}
