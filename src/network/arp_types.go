package network

import (
	"fmt"
	psErr "github.com/42milez/ProtocolStack/src/error"
	"github.com/42milez/ProtocolStack/src/ethernet"
	"sync"
	"time"
)

// Address Resolution Protocol (ARP) Parameters
// https://www.iana.org/assignments/arp-parameters/arp-parameters.xhtml

// EtherType
// https://en.wikipedia.org/wiki/EtherType#Examples

// Notes:
//  - Protocol Type is same as EtherType.

const (
	ArpCacheStateFree ArpCacheState = iota
	ArpCacheStateIncomplete
	ArpCacheStateResolved
	ArpCacheStateStatic
)
const ArpHwTypeEthernet ArpHwType = 1
const ArpOpRequest ArpOpcode = 1
const ArpOpReply ArpOpcode = 2
const ArpPacketSize = 28 // byte

type ArpCache struct {
	entries [ArpCacheSize]*ArpCacheEntry
	mtx     sync.Mutex
}
type ArpCacheEntry struct {
	State     ArpCacheState
	CreatedAt time.Time
	HA        ethernet.EthAddr
	PA        ArpProtoAddr
}
type ArpCacheState uint8
type ArpHdr struct {
	HT     ArpHwType        // hardware type
	PT     ethernet.EthType // protocol type
	HAL    uint8            // hardware address length
	PAL    uint8            // protocol address length
	Opcode ArpOpcode
}
type ArpHwType uint16
type ArpOpcode uint16
type ArpPacket struct {
	ArpHdr
	SHA ethernet.EthAddr // sender hardware address
	SPA ArpProtoAddr     // sender protocol address
	THA ethernet.EthAddr // target hardware address
	TPA ArpProtoAddr     // target protocol address
}
type ArpProtoAddr [V4AddrLen]byte

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

func (p *ArpCache) Clear(idx int) psErr.E {
	p.entries[idx].State = ArpCacheStateFree
	p.entries[idx].CreatedAt = time.Unix(0, 0)
	p.entries[idx].HA = ethernet.EthAddr{}
	p.entries[idx].PA = ArpProtoAddr{}
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

func (p *ArpCache) Init() {
	for i := range p.entries {
		p.entries[i] = &ArpCacheEntry{
			State:     ArpCacheStateFree,
			CreatedAt: time.Unix(0, 0),
		}
	}
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

func (v ArpHwType) String() string {
	switch v {
	case ArpHwTypeEthernet:
		return "Ethernet"
	default:
		return "Unknown"
	}
}

func (v ArpOpcode) String() string {
	switch v {
	case ArpOpRequest:
		return "REQUEST"
	case ArpOpReply:
		return "REPLY"
	default:
		return "UNKNOWN"
	}
}

func (p ArpProtoAddr) String() string {
	return fmt.Sprintf("%d.%d.%d.%d", p[0], p[1], p[2], p[3])
}
