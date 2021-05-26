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
const ArpMessageSize = 68

var cache *ArpCache

type ArpHwAddr [ethernet.EthAddrLen]byte

func (p *ArpHwAddr) String() string {
	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", p[0], p[1], p[2], p[3], p[4], p[5])
}

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

type ArpMessage struct {
	ArpHdr
	SHA ArpHwAddr    // sender hardware address
	SPA ArpProtoAddr // sender protocol address
	THA ArpHwAddr    // target hardware address
	TPA ArpProtoAddr // target protocol address
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

func (p *ArpCache) Delete() psErr.Error {
	return psErr.Error{Code: psErr.OK}
}

func (p *ArpCache) Insert() psErr.Error {
	return psErr.Error{Code: psErr.OK}
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

func (p *ArpCache) Update(msg *ArpMessage) psErr.Error {
	entry := p.Select(msg.SPA)
	if entry == nil {
		return psErr.Error{Code: psErr.NotFound}
	}
	p.mtx.Lock()
	defer p.mtx.Unlock()
	entry.State = ArpCacheStateResolved
	entry.SHA = msg.SHA
	entry.CreatedAt = time.Now()
	psLog.I("updated arp entry")
	psLog.I("\tSPA: %v", fmt.Sprintf("%d.%d.%d.%d", msg.SPA[0], msg.SPA[1], msg.SPA[2], msg.SPA[3]))
	psLog.I("\tSHA: %v", fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", msg.SHA[0], msg.SHA[1], msg.SHA[2], msg.SHA[3], msg.SHA[4], msg.SHA[5]))
	return psErr.Error{Code: psErr.OK}
}

func ArpInputHandler(payload []byte) psErr.Error {
	if len(payload) < ArpMessageSize {
		return psErr.Error{
			Code: psErr.InvalidPacket,
			Msg:  "message size is too small",
		}
	}

	buf := bytes.NewBuffer(payload)
	msg := ArpMessage{}
	if err := binary.Read(buf, binary.BigEndian, &msg); err != nil {
		return psErr.Error{Code: psErr.CantRead, Msg: err.Error()}
	}

	if msg.HT != ArpHwTypeEthernet || msg.HAL != ethernet.EthAddrLen {
		return psErr.Error{
			Code: psErr.InvalidPacket,
			Msg:  "invalid arp packet",
		}
	}

	if msg.PT != ethernet.EthTypeIpv4 || msg.PAL != V4AddrLen {
		return psErr.Error{
			Code: psErr.InvalidPacket,
			Msg:  "invalid arp packet",
		}
	}

	psLog.I("arp packet received")
	arpDump(&msg)

	_ = cache.Update(&msg)

	return psErr.Error{Code: psErr.OK}
}

func arpDump(msg *ArpMessage) {
	psLog.I("\thardware type:           %s", msg.HT)
	psLog.I("\tprotocol Type:           %s", msg.PT)
	psLog.I("\thardware address length: %d", msg.HAL)
	psLog.I("\tprotocol address length: %d", msg.PAL)
	psLog.I("\topcode:                  %s (%d)", msg.Opcode, uint16(msg.Opcode))
	psLog.I("\tsender hardware address: %s", msg.SHA.String())
	psLog.I("\tsender protocol address: %s", msg.SPA.String())
	psLog.I("\ttarget hardware address: %s", msg.THA.String())
	psLog.I("\ttarget hardware address: %s", msg.TPA.String())
}

//func arpReply() psErr.Error {
//	return psErr.Error{Code: psErr.OK}
//}

func init() {
	cache = &ArpCache{}
}
