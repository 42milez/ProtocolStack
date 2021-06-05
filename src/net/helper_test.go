package net

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestGenIface(t *testing.T) {
	want := &Iface{
		Family:    V4AddrFamily,
		Unicast:   ParseIP(LoopbackIpAddr),
		Netmask:   ParseIP(LoopbackNetmask),
		Broadcast: ParseIP(LoopbackBroadcast),
	}
	got := GenIface(LoopbackIpAddr, LoopbackNetmask, LoopbackBroadcast)
	if d := cmp.Diff(got, want); d != "" {
		t.Errorf("GenIface() differs: (-got +want)\n%s", d)
	}
}
