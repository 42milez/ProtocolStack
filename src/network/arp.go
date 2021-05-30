package network

import (
	"bytes"
	"encoding/binary"
	"fmt"
	psErr "github.com/42milez/ProtocolStack/src/error"
	"github.com/42milez/ProtocolStack/src/ethernet"
	psLog "github.com/42milez/ProtocolStack/src/log"
	"sync"
	"time"
)

const ArpCacheSize = 32
const ArpPacketSize = 28 // byte

var cache *ArpCache

type ArpProtoAddr [V4AddrLen]byte

func (p ArpProtoAddr) String() string {
	return fmt.Sprintf("%d.%d.%d.%d", p[0], p[1], p[2], p[3])
}

type ArpHdr struct {
	HT     ArpHwType        // hardware type
	PT     ethernet.EthType // protocol type
	HAL    uint8            // hardware address length
	PAL    uint8            // protocol address length
	Opcode ArpOpcode
}

type ArpPacket struct {
	ArpHdr
	SHA ethernet.EthAddr // sender hardware address
	SPA ArpProtoAddr     // sender protocol address
	THA ethernet.EthAddr // target hardware address
	TPA ArpProtoAddr     // target protocol address
}

type ArpCacheEntry struct {
	State     ArpCacheState
	CreatedAt time.Time
	HA        ethernet.EthAddr
	PA        ArpProtoAddr
}

type ArpCache struct {
	entries [ArpCacheSize]*ArpCacheEntry
	mtx     sync.Mutex
}

func (p *ArpCache) Init() {
	for i := range p.entries {
		p.entries[i] = &ArpCacheEntry{
			State:     ArpCacheStateFree,
			CreatedAt: time.Unix(0, 0),
		}
	}
}

func (p *ArpCache) Clear(idx int) psErr.E {
	p.entries[idx].State = ArpCacheStateFree
	p.entries[idx].CreatedAt = time.Unix(0, 0)
	p.entries[idx].HA = ethernet.EthAddr{}
	p.entries[idx].PA = ArpProtoAddr{}
	return psErr.OK
}

func (p *ArpCache) Add(packet *ArpPacket) psErr.E {
	if ret := p.Get(packet.SPA); ret != nil {
		return psErr.Exist
	}
	entry := p.danglingEntry()
	p.mtx.Lock()
	defer p.mtx.Unlock()
	entry.State = ArpCacheStateResolved
	entry.CreatedAt = time.Now()
	entry.HA = packet.SHA
	entry.PA = packet.SPA
	return psErr.OK
}

func (p *ArpCache) Get(ip [V4AddrLen]byte) *ArpCacheEntry {
	p.mtx.Lock()
	defer p.mtx.Unlock()
	for i, v := range p.entries {
		if v.PA == ip {
			return p.entries[i]
		}
	}
	return nil
}

func (p *ArpCache) Update(packet *ArpPacket) psErr.E {
	entry := p.Get(packet.SPA)
	if entry == nil {
		return psErr.NotFound
	}
	p.mtx.Lock()
	defer p.mtx.Unlock()
	entry.State = ArpCacheStateResolved
	entry.HA = packet.SHA
	entry.CreatedAt = time.Now()
	return psErr.OK
}

func (p *ArpCache) danglingEntry() *ArpCacheEntry {
	p.mtx.Lock()
	defer p.mtx.Unlock()
	oldest := p.entries[0]
	for _, entry := range p.entries {
		if entry.State == ArpCacheStateFree {
			return entry
		}
		if oldest.CreatedAt.After(entry.CreatedAt) {
			oldest = entry
		}
	}
	return oldest
}

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
	arpDump(&packet)

	iface := IfaceRepo.Get(dev, V4AddrFamily)
	if iface == nil {
		psLog.E(fmt.Sprintf("Interface for %s is not registered", dev.DevName()))
		return psErr.InterfaceNotFound
	}

	if isSameIP(packet.TPA, iface.Unicast) {
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

func arpDump(packet *ArpPacket) {
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
	arpDump(&packet)

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

func isSameIP(a ArpProtoAddr, b IP) bool {
	return a[0] == b[0] && a[1] == b[1] && a[2] == b[2] && a[3] == b[3]
}

func init() {
	cache = &ArpCache{}
	cache.Init()
}
