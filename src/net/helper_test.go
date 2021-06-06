package net

import (
	"github.com/42milez/ProtocolStack/src/mw"
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestGenIface(t *testing.T) {
	want := &mw.Iface{
		Family:    mw.V4AddrFamily,
		Unicast:   mw.ParseIP(mw.LoopbackIpAddr),
		Netmask:   mw.ParseIP(mw.LoopbackNetmask),
		Broadcast: mw.ParseIP(mw.LoopbackBroadcast),
	}
	got := GenIface(mw.LoopbackIpAddr, mw.LoopbackNetmask, mw.LoopbackBroadcast)
	if d := cmp.Diff(got, want); d != "" {
		t.Errorf("GenIface() differs: (-got +want)\n%s", d)
	}
}
