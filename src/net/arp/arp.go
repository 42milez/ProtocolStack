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
	psTime "github.com/42milez/ProtocolStack/src/time"
	"github.com/42milez/ProtocolStack/src/timer"
	"sync"
	"time"
)

const CacheSize = 32

var ARP *arp
var ArpCondCh chan timer.Condition
var ArpSigCh chan timer.Signal

type arp struct {
	cache arpCache
}

func (p *arp) Receive(packet []byte, dev mw.IDevice) psErr.E {
	if len(packet) < ArpPacketLen {
		psLog.E(fmt.Sprintf("ARP packet length is too short: %d bytes", len(packet)))
		return psErr.InvalidPacket
	}

	buf := bytes.NewBuffer(packet)
	arpPacket := ArpPacket{}
	if err := binary.Read(buf, binary.BigEndian, &arpPacket); err != nil {
		return psErr.ReadFromBufError
	}

	if arpPacket.HT != ArpHwTypeEthernet || arpPacket.HAL != mw.AddrLen {
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
		if err := p.cache.Renew(arpPacket.SPA, arpPacket.SHA, cacheStatusResolved); err == psErr.NotFound {
			_ = p.cache.Create(arpPacket.SHA, arpPacket.SPA, cacheStatusResolved)
		} else {
			psLog.I("ARP entry was renewed")
			psLog.I(fmt.Sprintf("\tspa: %s", arpPacket.SPA))
			psLog.I(fmt.Sprintf("\tsha: %s", arpPacket.SHA))
		}
		if arpPacket.Opcode == ArpOpRequest {
			if err := p.Reply(arpPacket.SHA, arpPacket.SPA, iface); err != psErr.OK {
				return psErr.Error
			}
		}
	} else {
		psLog.I("ARP packet was ignored (It was sent to different address)")
	}

	return psErr.OK
}

func (p *arp) Reply(tha mw.Addr, tpa ArpProtoAddr, iface *mw.Iface) psErr.E {
	packet := ArpPacket{
		ArpHdr: ArpHdr{
			HT:     ArpHwTypeEthernet,
			PT:     mw.IPv4,
			HAL:    mw.AddrLen,
			PAL:    mw.V4AddrLen,
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
		return psErr.WriteToBufError
	}

	if err := iface.Dev.Transmit(tha, buf.Bytes(), mw.ARP); err != psErr.OK {
		return psErr.Error
	}

	return psErr.OK
}

func (p *arp) Request(iface *mw.Iface, ip mw.IP) psErr.E {
	packet := ArpPacket{
		ArpHdr: ArpHdr{
			HT:     ArpHwTypeEthernet,
			PT:     mw.IPv4,
			HAL:    mw.AddrLen,
			PAL:    mw.V4AddrLen,
			Opcode: ArpOpRequest,
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

func (p *arp) Resolve(iface *mw.Iface, ip mw.IP) (mw.Addr, ArpStatus) {
	if iface.Dev.Type() != mw.DevTypeEthernet {
		psLog.E(fmt.Sprintf("Unsupported device type: %s", iface.Dev.Type()))
		return mw.Addr{}, ArpStatusError
	}

	if iface.Family != mw.V4AddrFamily {
		psLog.E(fmt.Sprintf("Unsupported address family: %s", iface.Family))
		return mw.Addr{}, ArpStatusError
	}

	entry := p.cache.GetEntry(ip.ToV4())
	if entry == nil {
		if err := p.cache.Create(mw.Addr{}, ip.ToV4(), cacheStatusIncomplete); err != psErr.OK {
			return mw.Addr{}, ArpStatusError
		}
		if err := p.Request(iface, ip); err != psErr.OK {
			return mw.Addr{}, ArpStatusError
		}
		return mw.Addr{}, ArpStatusIncomplete
	}

	return entry.HA, ArpStatusComplete
}

func (p *arp) RunTimer(wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		ArpCondCh <- timer.Condition{
			CurrentState: timer.Running,
		}
		for {
			select {
			case signal := <-ArpSigCh:
				if signal == timer.Stop {
					return
				}
			default:
				ret := p.cache.Expire()
				if len(ret) != 0 {
					psLog.I("ARP cache entries were expired:")
					for i, v := range ret {
						psLog.I(fmt.Sprintf("\t%d: %s", i+1, v))
					}
				}
				time.Sleep(time.Second)
			}
		}
	}()
}

func (p *arp) StopTimer() {
	ArpSigCh <- timer.Stop
}

type arpCacheEntry struct {
	Status    CacheStatus
	CreatedAt time.Time
	HA        mw.Addr
	PA        ArpProtoAddr
}

type arpCache struct {
	entries [CacheSize]*arpCacheEntry
	mtx     sync.Mutex
}

func (p *arpCache) Init() {
	for i := range p.entries {
		p.entries[i] = &arpCacheEntry{
			Status:    cacheStatusFree,
			CreatedAt: time.Unix(0, 0),
			HA:        mw.Addr{},
			PA:        ArpProtoAddr{},
		}
	}
}

func (p *arpCache) Create(ha mw.Addr, pa ArpProtoAddr, state CacheStatus) psErr.E {
	var entry *arpCacheEntry
	if entry = p.GetEntry(pa); entry != nil {
		return psErr.Exist
	}
	entry = p.GetReusableEntry()

	p.mtx.Lock()
	defer p.mtx.Unlock()

	entry.Status = state
	entry.CreatedAt = psTime.Time.Now()
	entry.HA = ha
	entry.PA = pa

	return psErr.OK
}

func (p *arpCache) Renew(pa ArpProtoAddr, ha mw.Addr, state CacheStatus) psErr.E {
	entry := p.GetEntry(pa)
	if entry == nil {
		return psErr.NotFound
	}

	p.mtx.Lock()
	defer p.mtx.Unlock()

	entry.Status = state
	entry.CreatedAt = psTime.Time.Now()
	entry.HA = ha

	return psErr.OK
}

func (p *arpCache) GetEntry(ip ArpProtoAddr) *arpCacheEntry {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	for i, v := range p.entries {
		if v.PA == ip {
			return p.entries[i]
		}
	}

	return nil
}

func (p *arpCache) GetReusableEntry() *arpCacheEntry {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	oldest := p.entries[0]
	for _, entry := range p.entries {
		if entry.Status == cacheStatusFree {
			return entry
		}
		if oldest.CreatedAt.After(entry.CreatedAt) {
			oldest = entry
		}
	}

	return oldest
}

func (p *arpCache) Clear(idx int) {
	p.entries[idx].Status = cacheStatusFree
	p.entries[idx].CreatedAt = time.Unix(0, 0)
	p.entries[idx].HA = mw.Addr{}
	p.entries[idx].PA = ArpProtoAddr{}
}

func (p *arpCache) Expire() (invalidations []string) {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	now := time.Now()
	for i, v := range p.entries {
		if v.CreatedAt != time.Unix(0, 0) && now.Sub(v.CreatedAt) > arpCacheLifetime {
			invalidations = append(invalidations, fmt.Sprintf("%s (%s)", v.PA, v.HA))
			p.Clear(i)
		}
	}

	return
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

func StartService() {
	go func() {
		for {
			packet := <-mw.ArpRxCh
			if err := ARP.Receive(packet.Content, packet.Dev); err != psErr.OK {
				return
			}
		}
	}()
}

func init() {
	ARP = &arp{}
	ARP.cache.Init()
	ArpCondCh = make(chan timer.Condition)
	ArpSigCh = make(chan timer.Signal)
}
