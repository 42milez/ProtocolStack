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
const ArpPacketSize = 68

var cache *ArpCache

type ArpProtoAddr [V4AddrLen]byte

func (p *ArpProtoAddr) String() string {
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
	SHA       [ethernet.EthAddrLen]byte
	SPA       [V4AddrLen]byte
}

type ArpCache struct {
	entries [ArpCacheSize]ArpCacheEntry
	mtx     sync.Mutex
}

func (p *ArpCache) Delete() psErr.E {
	return psErr.OK
}

func (p *ArpCache) Insert(msg *ArpPacket) psErr.E {
	return psErr.OK
}

func (p *ArpCache) Select(ip [V4AddrLen]byte) *ArpCacheEntry {
	p.mtx.Lock()
	defer p.mtx.Unlock()
	for i, v := range p.entries {
		if v.SPA == ip {
			return &p.entries[i]
		}
	}
	return nil
}

func (p *ArpCache) Update(msg *ArpPacket) psErr.E {
	entry := p.Select(msg.SPA)
	if entry == nil {
		return psErr.NotFound
	}
	p.mtx.Lock()
	defer p.mtx.Unlock()
	entry.State = ArpCacheStateResolved
	entry.SHA = msg.SHA
	entry.CreatedAt = time.Now()
	return psErr.OK
}

func ArpInputHandler(payload []byte, dev ethernet.IDevice) psErr.E {
	if len(payload) < ArpPacketSize {
		psLog.E(fmt.Sprintf("ARP packet size is too small: %d bytes", len(payload)))
		return psErr.InvalidPacket
	}

	buf := bytes.NewBuffer(payload)
	msg := ArpPacket{}
	if err := binary.Read(buf, binary.BigEndian, &msg); err != nil {
		psLog.E(fmt.Sprintf("binary.Read() failed: %s", err))
		return psErr.Error
	}

	if msg.HT != ArpHwTypeEthernet || msg.HAL != ethernet.EthAddrLen {
		psLog.E("Value of ARP packet header is invalid (Hardware)")
		return psErr.InvalidPacket
	}

	if msg.PT != ethernet.EthTypeIpv4 || msg.PAL != V4AddrLen {
		psLog.E("Value of ARP packet header is invalid (Protocol)")
		return psErr.InvalidPacket
	}

	psLog.I("Incoming ARP packet:")
	arpDump(&msg)

	iface := IfaceRepo.Get(dev, FamilyV4)
	if iface == nil {
		devName, _ := dev.Names()
		psLog.E(fmt.Sprintf("Interface for %s is not registered", devName))
		return psErr.InterfaceNotFound
	}

	if isSameIP(msg.TPA, iface.Unicast) {
		if err := cache.Update(&msg); err == psErr.NotFound {
			if err := cache.Insert(&msg); err != psErr.OK {
				psLog.E(fmt.Sprintf("ArpCache.Insert() failed: %s", err))
			}
		} else {
			psLog.I("updated arp entry")
			psLog.I(fmt.Sprintf(
				"\tSPA: %v",
				fmt.Sprintf("%d.%d.%d.%d", msg.SPA[0], msg.SPA[1], msg.SPA[2], msg.SPA[3])))
			psLog.I(fmt.Sprintf(
				"\tSHA: %v",
				fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", msg.SHA[0], msg.SHA[1], msg.SHA[2], msg.SHA[3], msg.SHA[4],
					msg.SHA[5])))
		}
		if msg.Opcode == ArpOpRequest {
			if err := arpReply(msg.SHA, msg.SPA, iface); err != psErr.OK {
				psLog.E(fmt.Sprintf("arpReply() failed: %s", err))
				return psErr.Error
			}
		}
	} else {
		psLog.I("Ignored arp packet (It was sent to different address)")
	}

	return psErr.OK
}

func arpDump(msg *ArpPacket) {
	psLog.I(fmt.Sprintf("\thardware type:           %s", msg.HT))
	psLog.I(fmt.Sprintf("\tprotocol Type:           %s", msg.PT))
	psLog.I(fmt.Sprintf("\thardware address length: %d", msg.HAL))
	psLog.I(fmt.Sprintf("\tprotocol address length: %d", msg.PAL))
	psLog.I(fmt.Sprintf("\topcode:                  %s (%d)", msg.Opcode, uint16(msg.Opcode)))
	psLog.I(fmt.Sprintf("\tsender hardware address: %s", msg.SHA))
	psLog.I(fmt.Sprintf("\tsender protocol address: %s", msg.SPA))
	psLog.I(fmt.Sprintf("\ttarget hardware address: %s", msg.THA))
	psLog.I(fmt.Sprintf("\ttarget hardware address: %s", msg.TPA))
}

func arpReply(tha ethernet.EthAddr, tpa ArpProtoAddr, iface *Iface) psErr.E {
	addr, _, _ := iface.Dev.EthAddrs()
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

func isSameIP(a ArpProtoAddr, b IP) bool {
	return a[0] == b[0] && a[1] == b[1] && a[2] == b[2] && a[3] == b[3]
}

func init() {
	cache = &ArpCache{}
}
