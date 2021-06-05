package arp

import (
	"fmt"
	psErr "github.com/42milez/ProtocolStack/src/error"
	psLog "github.com/42milez/ProtocolStack/src/log"
	"github.com/42milez/ProtocolStack/src/mw"
	psTime "github.com/42milez/ProtocolStack/src/time"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"sync"
	"testing"
	"time"
)

func TestCache_Create_1(t *testing.T) {
	defer cache.Init()

	ha := mw.EthAddr{0x11, 0x22, 0x33, 0x44, 0x55, 0x66}
	pa := mw.V4Addr{192, 168, 1, 1}
	state := resolved

	want := psErr.OK
	got := cache.Create(ha, pa, state)
	if got != want {
		t.Errorf("ArpCache.Add() = %s; want %s", got, want)
	}
}

func TestCache_Create_2(t *testing.T) {
	defer cache.Init()

	ha := mw.EthAddr{0x11, 0x22, 0x33, 0x44, 0x55, 0x66}
	pa := mw.V4Addr{192, 168, 1, 1}
	state := resolved

	_ = cache.Create(ha, pa, state)

	want := psErr.Exist
	got := cache.Create(ha, pa, state)
	if got != want {
		t.Errorf("ArpCache.Add() = %s; want %s", got, want)
	}
}

func TestCache_Create_3(t *testing.T) {
	ctrl, teardown := SetupCacheTest(t)
	defer teardown()
	defer cache.Init()

	want := psErr.OK
	var got psErr.E

	m := psTime.NewMockITime(ctrl)
	psTime.Time = m
	for i := cacheSize; i >= 0; i-- {
		ha := mw.EthAddr{0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
		pa := mw.V4Addr{192, 168, byte(i), 1}
		state := resolved
		createdAt, _ := time.Parse(time.RFC3339, fmt.Sprintf("2021-01-01T00:%02d:00Z", i))
		m.EXPECT().Now().Return(createdAt)
		got = cache.Create(ha, pa, state)
	}

	if got != want {
		t.Errorf("ArpCache.Add() = %s; want %s", got, want)
	}
}

func TestCache_Renew_1(t *testing.T) {
	ctrl, teardown := SetupCacheTest(t)
	defer teardown()
	defer cache.Init()

	createdAt, _ := time.Parse(time.RFC3339, "2021-01-01T00:00:00Z")
	m := psTime.NewMockITime(ctrl)
	m.EXPECT().Now().Return(createdAt).AnyTimes()
	psTime.Time = m

	ha1 := mw.EthAddr{0x11, 0x22, 0x33, 0x44, 0x55, 0x66}
	pa := mw.V4Addr{192, 168, 1, 1}
	state := resolved
	_ = cache.Create(ha1, pa, state)

	ha2 := mw.EthAddr{0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f}
	_ = cache.Renew(pa, ha2, state)

	want := &arpCacheEntry{
		Status:    resolved,
		CreatedAt: createdAt,
		HA:        mw.EthAddr{0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f},
		PA:        mw.V4Addr{192, 168, 1, 1},
	}

	got := cache.GetEntry(pa)
	if d := cmp.Diff(got, want); d != "" {
		t.Errorf("Renew() differs: (-got +want)\n%s", d)
	}
}

func TestCache_Renew_2(t *testing.T) {
	defer cache.Init()

	ha1 := mw.EthAddr{0x11, 0x22, 0x33, 0x44, 0x55, 0x66}
	pa := mw.V4Addr{192, 168, 1, 1}
	state := resolved
	_ = cache.Create(ha1, pa, state)

	want := psErr.NotFound
	got := cache.Renew(mw.V4Addr{192, 0, 2, 1}, mw.EthAddr{}, resolved)

	if got != want {
		t.Errorf("Renew() = %s; want %s", got, want)
	}
}

func TestCache_GetEntry_1(t *testing.T) {
	defer cache.Init()

	ha := mw.EthAddr{0x11, 0x22, 0x33, 0x44, 0x55, 0x66}
	pa := mw.V4Addr{192, 168, 1, 1}
	state := resolved
	_ = cache.Create(ha, pa, state)

	entry := cache.GetEntry(pa)
	if entry == nil {
		t.Errorf("arp cache entry does not exist")
	}
}

// GetEntry() returns nil when entry to match does not exist.
func TestCache_GetEntry_2(t *testing.T) {
	defer cache.Init()

	ha := mw.EthAddr{0x11, 0x22, 0x33, 0x44, 0x55, 0x66}
	pa := mw.V4Addr{192, 168, 1, 1}
	state := resolved
	_ = cache.Create(ha, pa, state)

	entry := cache.GetEntry(mw.V4Addr{192, 0, 2, 1})

	if entry != nil {
		t.Errorf("unexpected arp cache entry exist")
	}
}

func TestTimer_1(t *testing.T) {
	ctrl, teardown := SetupCacheTest(t)
	defer teardown()
	defer cache.Init()

	createdAt, _ := time.Parse(time.RFC3339, "2021-01-01T00:00:00Z")
	m := psTime.NewMockITime(ctrl)
	m.EXPECT().Now().Return(createdAt)
	psTime.Time = m

	pa := mw.V4Addr{192, 168, 1, 1}
	_ = cache.Create(mw.EthAddr{0x11, 0x12, 0x13, 0x14, 0x15, 0x16}, pa, resolved)

	var wg sync.WaitGroup
	_ = StartService(&wg)
	<-MonitorCh
	<-MonitorCh
	<-MonitorCh
	StopService()
	wg.Wait()

	got := cache.GetEntry(pa)
	if got != nil {
		t.Errorf("ARP cache is not expired")
	}
}

func TestTimer_2(t *testing.T) {
	_, teardown := SetupCacheTest(t)
	defer teardown()
	defer cache.Init()

	var wg sync.WaitGroup
	_ = StartService(&wg)
	<-MonitorCh
	<-MonitorCh
	<-MonitorCh
	StopService()
	wg.Wait()

	got := cache.GetEntry(mw.V4Addr{192, 168, 0, 1})
	if got != nil {
		t.Errorf("ARP cache exists")
	}
}

func SetupCacheTest(t *testing.T) (ctrl *gomock.Controller, teardown func()) {
	psLog.DisableOutput()
	backup := psTime.Time
	ctrl = gomock.NewController(t)

	teardown = func() {
		psLog.EnableOutput()
		psTime.Time = backup
		ctrl.Finish()
	}

	return
}
