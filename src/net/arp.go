package net

import (
	"bytes"
	"encoding/binary"
	"fmt"
	psErr "github.com/42milez/ProtocolStack/src/error"
	"github.com/42milez/ProtocolStack/src/eth"
	psLog "github.com/42milez/ProtocolStack/src/log"
	psTime "github.com/42milez/ProtocolStack/src/time"
	"github.com/42milez/ProtocolStack/src/timer"
	"sync"
	"time"
)

var ARP *arp
var ArpCondCh chan timer.Condition
var ArpSigCh chan timer.Signal

type arp struct {
	cache arpCache
}

func (p *arp) Receive(packet []byte, dev eth.IDevice) psErr.E {
	if len(packet) < ArpPacketLen {
		psLog.E(fmt.Sprintf("ARP packet length is too short: %d bytes", len(packet)))
		return psErr.InvalidPacket
	}

	buf := bytes.NewBuffer(packet)
	arpPacket := ArpPacket{}
	if err := binary.Read(buf, binary.BigEndian, &arpPacket); err != nil {
		return psErr.ReadFromBufError
	}

	if arpPacket.HT != ArpHwTypeEthernet || arpPacket.HAL != eth.EthAddrLen {
		psLog.E("Value of ARP packet header is invalid (Hardware)")
		return psErr.InvalidPacket
	}

	if arpPacket.PT != eth.EthTypeIpv4 || arpPacket.PAL != V4AddrLen {
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
		if err := p.cache.Renew(arpPacket.SPA, arpPacket.SHA, ArpCacheStateResolved); err == psErr.NotFound {
			_ = p.cache.Create(arpPacket.SHA, arpPacket.SPA, ArpCacheStateResolved)
		} else {
			psLog.I("ARP entry was renewed")
			psLog.I(fmt.Sprintf("\tspa: %s", arpPacket.SPA))
			psLog.I(fmt.Sprintf("\tsha: %s", arpPacket.SHA))
		}
		if arpPacket.Opcode == ArpOpRequest {
			if err := p.SendReply(arpPacket.SHA, arpPacket.SPA, iface); err != psErr.OK {
				return psErr.Error
			}
		}
	} else {
		psLog.I("ARP packet was ignored (It was sent to different address)")
	}

	return psErr.OK
}

func (p *arp) SendReply(tha eth.EthAddr, tpa ArpProtoAddr, iface *Iface) psErr.E {
	packet := ArpPacket{
		ArpHdr: ArpHdr{
			HT:     ArpHwTypeEthernet,
			PT:     eth.EthTypeIpv4,
			HAL:    eth.EthAddrLen,
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
		return psErr.WriteToBufError
	}

	if err := iface.Dev.Transmit(tha, buf.Bytes(), eth.EthTypeArp); err != psErr.OK {
		return psErr.Error
	}

	return psErr.OK
}

func (p *arp) SendRequest(iface *Iface, ip IP) psErr.E {
	packet := ArpPacket{
		ArpHdr: ArpHdr{
			HT:     ArpHwTypeEthernet,
			PT:     eth.EthTypeIpv4,
			HAL:    eth.EthAddrLen,
			PAL:    V4AddrLen,
			Opcode: ArpOpRequest,
		},
		SHA: iface.Dev.Addr(),
		SPA: iface.Unicast.ToV4(),
		THA: eth.EthAddr{},
		TPA: ip.ToV4(),
	}

	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, &packet); err != nil {
		return psErr.WriteToBufError
	}
	payload := buf.Bytes()

	psLog.I("Outgoing ARP packet")
	dumpArpPacket(&packet)

	if err := Transmit(eth.EthAddrBroadcast, payload, eth.EthTypeArp, iface); err != psErr.OK {
		return psErr.Error
	}

	return psErr.OK
}

func (p *arp) Resolve(iface *Iface, ip IP) (eth.EthAddr, ArpStatus) {
	if iface.Dev.Type() != eth.DevTypeEthernet {
		psLog.E(fmt.Sprintf("Unsupported device type: %s", iface.Dev.Type()))
		return eth.EthAddr{}, ArpStatusError
	}

	if iface.Family != V4AddrFamily {
		psLog.E(fmt.Sprintf("Unsupported address family: %s", iface.Family))
		return eth.EthAddr{}, ArpStatusError
	}

	ethAddr, found := p.cache.EthAddr(ip.ToV4())
	if !found {
		if err := p.cache.Create(eth.EthAddr{}, ip.ToV4(), ArpCacheStateIncomplete); err != psErr.OK {
			return eth.EthAddr{}, ArpStatusError
		}
		if err := p.SendRequest(iface, ip); err != psErr.OK {
			return eth.EthAddr{}, ArpStatusError
		}
		return eth.EthAddr{}, ArpStatusIncomplete
	}

	return ethAddr, ArpStatusComplete
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
				ret := p.cache.expire()
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
	State     ArpCacheState
	CreatedAt time.Time
	HA        eth.EthAddr
	PA        ArpProtoAddr
}

type arpCache struct {
	entries [ArpCacheSize]*arpCacheEntry
	mtx     sync.Mutex
}

func (p *arpCache) Init() {
	for i := range p.entries {
		p.entries[i] = &arpCacheEntry{
			State:     ArpCacheStateFree,
			CreatedAt: time.Unix(0, 0),
			HA:        eth.EthAddr{},
			PA:        ArpProtoAddr{},
		}
	}
}

func (p *arpCache) Create(ha eth.EthAddr, pa ArpProtoAddr, state ArpCacheState) psErr.E {
	var entry *arpCacheEntry
	if entry = p.get(pa); entry != nil {
		return psErr.Exist
	}
	entry = p.reusableEntry()

	p.mtx.Lock()
	defer p.mtx.Unlock()

	entry.State = state
	entry.CreatedAt = psTime.Time.Now()
	entry.HA = ha
	entry.PA = pa

	return psErr.OK
}

func (p *arpCache) Renew(pa ArpProtoAddr, ha eth.EthAddr, state ArpCacheState) psErr.E {
	entry := p.get(pa)
	if entry == nil {
		return psErr.NotFound
	}

	p.mtx.Lock()
	defer p.mtx.Unlock()

	entry.State = state
	entry.CreatedAt = psTime.Time.Now()
	entry.HA = ha

	return psErr.OK
}

func (p *arpCache) EthAddr(pa ArpProtoAddr) (eth.EthAddr, bool) {
	if entry := p.get(pa); entry != nil {
		return entry.HA, true
	} else {
		return eth.EthAddr{}, false
	}
}

func (p *arpCache) clear(idx int) {
	p.entries[idx].State = ArpCacheStateFree
	p.entries[idx].CreatedAt = time.Unix(0, 0)
	p.entries[idx].HA = eth.EthAddr{}
	p.entries[idx].PA = ArpProtoAddr{}
}

func (p *arpCache) get(ip ArpProtoAddr) *arpCacheEntry {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	for i, v := range p.entries {
		if v.PA == ip {
			return p.entries[i]
		}
	}

	return nil
}

func (p *arpCache) reusableEntry() *arpCacheEntry {
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

func (p *arpCache) expire() (invalidations []string) {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	now := time.Now()
	for i, v := range p.entries {
		if v.CreatedAt != time.Unix(0, 0) && now.Sub(v.CreatedAt) > arpCacheLifetime {
			invalidations = append(invalidations, fmt.Sprintf("%s (%s)", v.PA, v.HA))
			p.clear(i)
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

func init() {
	ARP = &arp{}
	ARP.cache.Init()
	ArpCondCh = make(chan timer.Condition)
	ArpSigCh = make(chan timer.Signal)
}
