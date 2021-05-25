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

type ArpHdr struct {
	hType ArpHwType // Hardware type
	pType ArpOp     // Protocol type
	hlen uint8      // Hardware address length
	plen uint8      // Protocol address length
}

type ArpMessage struct {
	ArpHdr
	SHA [ethernet.EthAddrLen]byte // Sender hardware address
	SPA [V4AddrLen]byte           // Sender protocol address
	THA [ethernet.EthAddrLen]byte // Target hardware address
	TPA [V4AddrLen]byte           // Target protocol address
}

type ArpCacheEntry struct {
	State ArpCacheState
	CreatedAt time.Time
	SHA [ethernet.EthAddrLen]byte
	SPA [V4AddrLen]byte
}

type ArpCache struct {
	entries [ArpCacheSize]ArpCacheEntry
	mtx sync.Mutex
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

func ArpInputHandler(payload []byte, dev ethernet.IDevice) psErr.Error {
	if len(payload) < ArpMessageSize {
		return psErr.Error{
			Code: psErr.InvalidPacket,
			Msg: "message size is too small",
		}
	}

	buf := bytes.NewBuffer(payload)
	msg := ArpMessage{}
	if err := binary.Read(buf, binary.LittleEndian, &msg); err != nil {
		return psErr.Error{Code: psErr.CantRead, Msg: err.Error()}
	}

	if msg.hType != ArpHwTypeEthernet || msg.hlen != ethernet.EthAddrLen {
		return psErr.Error{
			Code: psErr.InvalidPacket,
			Msg: "invalid arp packet",
		}
	}

	if msg.pType != ArpOpRequest || msg.plen != V4AddrLen {
		return psErr.Error{
			Code: psErr.InvalidPacket,
			Msg: "invalid arp packet",
		}
	}

	psLog.I("arp packet received")
	arpDump(&msg)

	//iface := middleware.GetIface(dev, FamilyV4)

	return psErr.Error{Code: psErr.OK}
}

func arpDump(msg *ArpMessage) {}

func arpReply() psErr.Error {
	return psErr.Error{Code: psErr.OK}
}

func init() {
	cache = &ArpCache{}
}
