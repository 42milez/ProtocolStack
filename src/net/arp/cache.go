package arp

import (
	"fmt"
	psErr "github.com/42milez/ProtocolStack/src/error"
	"github.com/42milez/ProtocolStack/src/mw"
	psTime "github.com/42milez/ProtocolStack/src/time"
	"sync"
	"time"
)

const cacheSize = 32
const lifetime = 24 * time.Hour
const (
	free cacheStatus = iota
	incomplete
	resolved
	//static
)

var cache *arpCache

type arpCache struct {
	entries [cacheSize]*arpCacheEntry
	mtx     sync.Mutex
}

func (p *arpCache) Init() {
	for i := range p.entries {
		p.entries[i] = &arpCacheEntry{
			Status:    free,
			CreatedAt: time.Unix(0, 0),
			HA:        mw.EthAddr{},
			PA:        mw.V4Addr{},
		}
	}
}

func (p *arpCache) Create(ha mw.EthAddr, pa mw.V4Addr, st cacheStatus) error {
	var entry *arpCacheEntry
	if entry = p.GetEntry(pa); entry != nil {
		return psErr.Exist
	}
	entry = p.GetReusableEntry()

	p.mtx.Lock()
	defer p.mtx.Unlock()

	entry.Status = st
	entry.CreatedAt = psTime.Time.Now()
	entry.HA = ha
	entry.PA = pa

	return psErr.OK
}

func (p *arpCache) Renew(pa mw.V4Addr, ha mw.EthAddr, st cacheStatus) error {
	entry := p.GetEntry(pa)
	if entry == nil {
		return psErr.NotFound
	}

	p.mtx.Lock()
	defer p.mtx.Unlock()

	entry.Status = st
	entry.CreatedAt = psTime.Time.Now()
	entry.HA = ha

	return psErr.OK
}

func (p *arpCache) GetEntry(ip mw.V4Addr) *arpCacheEntry {
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
		if entry.Status == free {
			return entry
		}
		if oldest.CreatedAt.After(entry.CreatedAt) {
			oldest = entry
		}
	}

	return oldest
}

func (p *arpCache) Clear(idx int) {
	p.entries[idx].Status = free
	p.entries[idx].CreatedAt = time.Unix(0, 0)
	p.entries[idx].HA = mw.EthAddr{}
	p.entries[idx].PA = mw.V4Addr{}
}

func (p *arpCache) Expire() (invalidations []string) {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	now := time.Now()
	for i, v := range p.entries {
		if v.CreatedAt != time.Unix(0, 0) && now.Sub(v.CreatedAt) > lifetime {
			invalidations = append(invalidations, fmt.Sprintf("%s (%s)", v.PA, v.HA))
			p.Clear(i)
		}
	}

	return
}

type arpCacheEntry struct {
	Status    cacheStatus
	CreatedAt time.Time
	HA        mw.EthAddr
	PA        mw.V4Addr
}

type cacheStatus uint8

func init() {
	cache = &arpCache{}
	cache.Init()
}
