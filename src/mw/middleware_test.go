package mw

import (
	"math/rand"
	"testing"
)

func TestRandU8(t *testing.T) {
	rand.Seed(0)
	want := uint8(250)
	got := RandU8()
	if got != want {
		t.Errorf("RandU8() = %d; want %d", got, want)
	}
}

func TestRandU16(t *testing.T) {
	rand.Seed(0)
	want := uint16(12282)
	got := RandU16()
	if got != want {
		t.Errorf("RandU16() = %d; want %d", got, want)
	}
}

func TestRandU32(t *testing.T) {
	rand.Seed(0)
	want := uint32(4059586549)
	got := RandU32()
	if got != want {
		t.Errorf("RandU32() = %d; want %d", got, want)
	}
}
