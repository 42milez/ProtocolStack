package network

import (
	"github.com/42milez/ProtocolStack/src/ethernet"
	psTime "github.com/42milez/ProtocolStack/src/time"
	"sync"
	"testing"
	"time"
)

func TestRunArpTimer_1(t *testing.T) {
	ctrl, teardown := setup(t)
	defer teardown()
	defer cache.Init()

	createdAt, _ := time.Parse(time.RFC3339, "2021-01-01T00:00:00Z")
	m := psTime.NewMockITime(ctrl)
	m.EXPECT().Now().Return(createdAt)
	psTime.Time = m

	ha := ethernet.EthAddr{0x11, 0x22, 0x33, 0x44, 0x55, 0x66}
	pa := ArpProtoAddr{192, 168, 1, 1}
	state := ArpCacheStateResolved
	_ = cache.Add(ha, pa, state)

	var wg sync.WaitGroup
	RunArpTimer(&wg)
	<-ArpCondCh
	StopArpTimer()
	wg.Wait()

	_, got := cache.EthAddr(pa)
	if got {
		t.Errorf("ARP cache is not expired")
	}
}

func TestRunArpTimer_2(t *testing.T) {
	defer cache.Init()

	var wg sync.WaitGroup
	RunArpTimer(&wg)
	<-ArpCondCh
	StopArpTimer()
	wg.Wait()

	_, got := cache.EthAddr(ArpProtoAddr{192, 168, 0, 1})
	if got {
		t.Errorf("ARP cache exists")
	}
}
