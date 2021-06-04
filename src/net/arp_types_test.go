package net

import (
	"fmt"
	psErr "github.com/42milez/ProtocolStack/src/error"
	"github.com/42milez/ProtocolStack/src/eth"
	psLog "github.com/42milez/ProtocolStack/src/log"
	psTime "github.com/42milez/ProtocolStack/src/time"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"testing"
	"time"
)

func isEthAddrAllZero(addr [eth.EthAddrLen]byte) bool {
	for _, v := range addr {
		if v != 0 {
			return false
		}
	}
	return true
}

func setupArpTypesTest(t *testing.T) (ctrl *gomock.Controller, teardown func()) {
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

func TestArpCache_Add_1(t *testing.T) {
	defer ARP.cache.Init()

	ha := eth.EthAddr{0x11, 0x22, 0x33, 0x44, 0x55, 0x66}
	pa := ArpProtoAddr{192, 168, 1, 1}
	state := ArpCacheStateResolved

	want := psErr.OK
	got := ARP.cache.Create(ha, pa, state)
	if got != want {
		t.Errorf("ArpCache.Add() = %s; want %s", got, want)
	}
}

func TestArpCache_Add_2(t *testing.T) {
	defer ARP.cache.Init()

	ha := eth.EthAddr{0x11, 0x22, 0x33, 0x44, 0x55, 0x66}
	pa := ArpProtoAddr{192, 168, 1, 1}
	state := ArpCacheStateResolved

	_ = ARP.cache.Create(ha, pa, state)

	want := psErr.Exist
	got := ARP.cache.Create(ha, pa, state)
	if got != want {
		t.Errorf("ArpCache.Add() = %s; want %s", got, want)
	}
}

func TestArpCache_Add_3(t *testing.T) {
	ctrl, teardown := setupArpTypesTest(t)
	defer teardown()
	defer ARP.cache.Init()

	want := psErr.OK
	var got psErr.E

	m := psTime.NewMockITime(ctrl)
	psTime.Time = m
	for i := ArpCacheSize; i >= 0; i-- {
		ha := eth.EthAddr{0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
		pa := ArpProtoAddr{192, 168, byte(i), 1}
		state := ArpCacheStateResolved
		createdAt, _ := time.Parse(time.RFC3339, fmt.Sprintf("2021-01-01T00:%02d:00Z", i))
		m.EXPECT().Now().Return(createdAt)
		got = ARP.cache.Create(ha, pa, state)
	}

	if got != want {
		t.Errorf("ArpCache.Add() = %s; want %s", got, want)
	}
}

func TestArpCache_Renew_1(t *testing.T) {
	ctrl, teardown := setupArpTypesTest(t)
	defer teardown()
	defer ARP.cache.Init()

	createdAt, _ := time.Parse(time.RFC3339, "2021-01-01T00:00:00Z")
	m := psTime.NewMockITime(ctrl)
	m.EXPECT().Now().Return(createdAt).AnyTimes()
	psTime.Time = m

	ha1 := eth.EthAddr{0x11, 0x22, 0x33, 0x44, 0x55, 0x66}
	pa := ArpProtoAddr{192, 168, 1, 1}
	state := ArpCacheStateResolved
	_ = ARP.cache.Create(ha1, pa, state)

	ha2 := eth.EthAddr{0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f}
	_ = ARP.cache.Renew(pa, ha2, state)

	want := &arpCacheEntry{
		State:     ArpCacheStateResolved,
		CreatedAt: createdAt,
		HA:        eth.EthAddr{0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f},
		PA:        ArpProtoAddr{192, 168, 1, 1},
	}

	got := ARP.cache.GetEntry(pa)
	if d := cmp.Diff(got, want); d != "" {
		t.Errorf("Renew() differs: (-got +want)\n%s", d)
	}
}

func TestArpCache_Renew_2(t *testing.T) {
	defer ARP.cache.Init()

	ha1 := eth.EthAddr{0x11, 0x22, 0x33, 0x44, 0x55, 0x66}
	pa := ArpProtoAddr{192, 168, 1, 1}
	state := ArpCacheStateResolved
	_ = ARP.cache.Create(ha1, pa, state)

	want := psErr.NotFound
	got := ARP.cache.Renew(ArpProtoAddr{192, 0, 2, 1}, eth.EthAddr{}, ArpCacheStateResolved)

	if got != want {
		t.Errorf("Renew() = %s; want %s", got, want)
	}
}

func TestArpCache_EthAddr_1(t *testing.T) {
	defer ARP.cache.Init()

	ha := eth.EthAddr{0x11, 0x22, 0x33, 0x44, 0x55, 0x66}
	pa := ArpProtoAddr{192, 168, 1, 1}
	state := ArpCacheStateResolved
	_ = ARP.cache.Create(ha, pa, state)

	ethAddr, found := ARP.cache.EthAddr(pa)

	if ethAddr != ha || !found {
		t.Errorf("ArpCache.EthAddr() = %v, %t; want %v, %t", ethAddr, found, ha, true)
	}
}

func TestArpCache_EthAddr_2(t *testing.T) {
	defer ARP.cache.Init()

	ha := eth.EthAddr{0x11, 0x22, 0x33, 0x44, 0x55, 0x66}
	pa := ArpProtoAddr{192, 168, 1, 1}
	state := ArpCacheStateResolved
	_ = ARP.cache.Create(ha, pa, state)

	ethAddr, found := ARP.cache.EthAddr(ArpProtoAddr{192, 0, 2, 1})

	if !isEthAddrAllZero(ethAddr) || found {
		t.Errorf("ArpCache.EthAddr() = %v, %t; want %v, %t", ethAddr, found, ha, true)
	}
}

func TestArpHwType_String(t *testing.T) {
	want := arpHwTypes[ArpHwTypeEthernet]
	got := ArpHwTypeEthernet.String()
	if got != want {
		t.Errorf("ArpHwType.String() = %s; want %s", got, want)
	}
}

func TestArpOpcode_String(t *testing.T) {
	want := arpOpCodes[ArpOpRequest]
	got := ArpOpRequest.String()
	if got != want {
		t.Errorf("ArpOpcode.String() = %s; want %s", got, want)
	}
}

func TestArpProtoAddr_String(t *testing.T) {
	want := "192.168.1.1"
	got := ArpProtoAddr{192, 168, 1, 1}.String()
	if got != want {
		t.Errorf("ArpProtoAddr.String() = %s; want %s", got, want)
	}
}