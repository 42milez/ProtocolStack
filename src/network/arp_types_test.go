package network

import (
	psErr "github.com/42milez/ProtocolStack/src/error"
	"github.com/42milez/ProtocolStack/src/ethernet"
	psTime "github.com/42milez/ProtocolStack/src/time"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"testing"
	"time"
)

func isEthAddrAllZero(addr [ethernet.EthAddrLen]byte) bool {
	for _, v := range addr {
		if v != 0 {
			return false
		}
	}
	return true
}

func setup(t *testing.T) (ctrl *gomock.Controller, teardown func()) {
	ctrl = gomock.NewController(t)
	teardown = func() {
		ctrl.Finish()
	}
	return
}

func TestArpCache_Add_1(t *testing.T) {
	defer cache.Init()

	ha := ethernet.EthAddr{0x11, 0x22, 0x33, 0x44, 0x55, 0x66}
	pa := ArpProtoAddr{192, 168, 1, 1}
	state := ArpCacheStateResolved

	want := psErr.OK
	got := cache.Add(ha, pa, state)
	if got != want {
		t.Errorf("ArpCache.Add() = %s; want %s", got, want)
	}
}

func TestArpCache_Add_2(t *testing.T) {
	defer cache.Init()

	ha := ethernet.EthAddr{0x11, 0x22, 0x33, 0x44, 0x55, 0x66}
	pa := ArpProtoAddr{192, 168, 1, 1}
	state := ArpCacheStateResolved

	_ = cache.Add(ha, pa, state)

	want := psErr.Exist
	got := cache.Add(ha, pa, state)
	if got != want {
		t.Errorf("ArpCache.Add() = %s; want %s", got, want)
	}
}

func TestArpCache_Add_3(t *testing.T) {
	defer cache.Init()

	want := psErr.OK
	var got psErr.E

	for i := 0; i < ArpCacheSize+1; i++ {
		ha := ethernet.EthAddr{0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
		pa := ArpProtoAddr{192, 168, byte(i), 1}
		state := ArpCacheStateResolved
		got = cache.Add(ha, pa, state)
	}

	if got != want {
		t.Errorf("ArpCache.Add() = %s; want %s", got, want)
	}
}

func TestArpCache_Renew_1(t *testing.T) {
	ctrl, teardown := setup(t)
	defer teardown()
	defer cache.Init()

	createdAt, _ := time.Parse(time.RFC3339, "2021-01-01T00:00:00Z")
	m := psTime.NewMockITime(ctrl)
	m.EXPECT().Now().Return(createdAt).AnyTimes()
	psTime.Time = m

	ha1 := ethernet.EthAddr{0x11, 0x22, 0x33, 0x44, 0x55, 0x66}
	pa := ArpProtoAddr{192, 168, 1, 1}
	state := ArpCacheStateResolved
	_ = cache.Add(ha1, pa, state)

	ha2 := ethernet.EthAddr{0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f}
	_ = cache.Renew(pa, ha2, state)

	want := &ArpCacheEntry{
		State:     ArpCacheStateResolved,
		CreatedAt: createdAt,
		HA:        ethernet.EthAddr{0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f},
		PA:        ArpProtoAddr{192, 168, 1, 1},
	}

	got := cache.get(pa)
	if d := cmp.Diff(got, want); d != "" {
		t.Errorf("Renew() differs: (-got +want)\n%s", d)
	}
}

func TestArpCache_Renew_2(t *testing.T) {
	defer cache.Init()

	ha1 := ethernet.EthAddr{0x11, 0x22, 0x33, 0x44, 0x55, 0x66}
	pa := ArpProtoAddr{192, 168, 1, 1}
	state := ArpCacheStateResolved
	_ = cache.Add(ha1, pa, state)

	want := psErr.NotFound
	got := cache.Renew(ArpProtoAddr{192, 0, 2, 1}, ethernet.EthAddr{}, ArpCacheStateResolved)

	if got != want {
		t.Errorf("Renew() = %s; want %s", got, want)
	}
}

func TestArpCache_EthAddr_1(t *testing.T) {
	defer cache.Init()

	ha := ethernet.EthAddr{0x11, 0x22, 0x33, 0x44, 0x55, 0x66}
	pa := ArpProtoAddr{192, 168, 1, 1}
	state := ArpCacheStateResolved
	_ = cache.Add(ha, pa, state)

	ethAddr, found := cache.EthAddr(pa)

	if ethAddr != ha || !found {
		t.Errorf("ArpCache.EthAddr() = %v, %t; want %v, %t", ethAddr, found, ha, true)
	}
}

func TestArpCache_EthAddr_2(t *testing.T) {
	defer cache.Init()

	ha := ethernet.EthAddr{0x11, 0x22, 0x33, 0x44, 0x55, 0x66}
	pa := ArpProtoAddr{192, 168, 1, 1}
	state := ArpCacheStateResolved
	_ = cache.Add(ha, pa, state)

	ethAddr, found := cache.EthAddr(ArpProtoAddr{192, 0, 2, 1})

	if !isEthAddrAllZero(ethAddr) || found {
		t.Errorf("ArpCache.EthAddr() = %v, %t; want %v, %t", ethAddr, found, ha, true)
	}
}
