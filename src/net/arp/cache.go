package arp

import (
	"fmt"
	psErr "github.com/42milez/ProtocolStack/src/error"
	"github.com/42milez/ProtocolStack/src/mw"
	psTime "github.com/42milez/ProtocolStack/src/time"
	"sync"
	"time"
)

var cache *arpCache

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

type arpCacheEntry struct {
	Status    CacheStatus
	CreatedAt time.Time
	HA        mw.Addr
	PA        ArpProtoAddr
}

func init() {
	cache = &arpCache{}
	cache.Init()
}
