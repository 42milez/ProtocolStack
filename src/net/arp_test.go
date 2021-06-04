package net

import (
	"bytes"
	"encoding/binary"
	psErr "github.com/42milez/ProtocolStack/src/error"
	"github.com/42milez/ProtocolStack/src/eth"
	psLog "github.com/42milez/ProtocolStack/src/log"
	psTime "github.com/42milez/ProtocolStack/src/time"
	"github.com/golang/mock/gomock"
	"sync"
	"testing"
	"time"
)

func TestArpInputHandler_1(t *testing.T) {
	ctrl, teardown := SetupArpInputHandlerTest(t)
	defer teardown()

	ethAddr := [eth.EthAddrLen]byte{0x11, 0x12, 0x13, 0x14, 0x15, 0x16}
	mockDev := eth.NewMockIDevice(ctrl)
	mockDev.EXPECT().Addr().Return(ethAddr)
	mockDev.EXPECT().Transmit(Any, Any, Any).Return(psErr.OK)

	mockIfaceRepo := NewMockIIfaceRepo(ctrl)
	mockIfaceRepo.EXPECT().Lookup(Any, Any).Return(&Iface{
		Family:    V4AddrFamily,
		Unicast:   ParseIP("192.0.2.2"),
		Netmask:   ParseIP("255.255.255.0"),
		Broadcast: ParseIP("192.0.2.255"),
		Dev:       mockDev,
	})
	IfaceRepo = mockIfaceRepo

	packet := Builder.Default()
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.BigEndian, packet)
	dev := &eth.TapDevice{
		Device: eth.Device{
			Type_: eth.DevTypeEthernet,
			Name_: "net0",
			Addr_: ethAddr,
			Flag_: eth.DevFlagBroadcast | eth.DevFlagNeedArp,
			MTU_:  eth.EthPayloadLenMax,
			Priv_: eth.Privilege{
				FD:   3,
				Name: "tap0",
			},
		},
	}

	want := psErr.OK
	got := ARP.Receive(buf.Bytes(), dev)
	if got != want {
		t.Errorf("ArpInputHandler() = %s; want %s", got, want)
	}
}

// Fail when ARP packet is too short.
func TestArpInputHandler_2(t *testing.T) {
	_, teardown := SetupArpInputHandlerTest(t)
	defer teardown()

	var packet []byte
	dev := &eth.TapDevice{}

	want := psErr.InvalidPacket
	got := ARP.Receive(packet, dev)
	if got != want {
		t.Errorf("ArpInputHandler() = %s; want %s", got, want)
	}
}

// Fail when hardware type of ARP packet is invalid.
// Fail when protocol type of ARP packet is invalid.
func TestArpInputHandler_3(t *testing.T) {
	_, teardown := SetupArpInputHandlerTest(t)
	defer teardown()

	packet := Builder.CustomHT(ArpHwType(0xffff))
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.BigEndian, packet)
	dev := &eth.TapDevice{}

	want := psErr.InvalidPacket
	got := ARP.Receive(buf.Bytes(), dev)
	if got != want {
		t.Errorf("ArpInputHandler() = %s; want %s", got, want)
	}

	packet = Builder.CustomPT(eth.EthType(0xffff))
	buf = new(bytes.Buffer)
	_ = binary.Write(buf, binary.BigEndian, packet)
	dev = &eth.TapDevice{}

	want = psErr.InvalidPacket
	got = ARP.Receive(buf.Bytes(), dev)
	if got != want {
		t.Errorf("ArpInputHandler() = %s; want %s", got, want)
	}
}

// Fail when Iface is not found.
func TestArpInputHandler_4(t *testing.T) {
	ctrl, teardown := SetupArpInputHandlerTest(t)
	defer teardown()

	mockIfaceRepo := NewMockIIfaceRepo(ctrl)
	mockIfaceRepo.EXPECT().Lookup(Any, Any).Return(nil)
	IfaceRepo = mockIfaceRepo

	packet := Builder.Default()
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.BigEndian, packet)
	dev := &eth.TapDevice{}

	want := psErr.InterfaceNotFound
	got := ARP.Receive(buf.Bytes(), dev)
	if got != want {
		t.Errorf("ArpInputHandler() = %s; want %s", got, want)
	}
}

// Fail when arpReply() returns error.
func TestArpInputHandler_5(t *testing.T) {
	ctrl, teardown := SetupArpInputHandlerTest(t)
	defer teardown()

	ethAddr := [eth.EthAddrLen]byte{0x11, 0x12, 0x13, 0x14, 0x15, 0x16}
	mockDev := eth.NewMockIDevice(ctrl)
	mockDev.EXPECT().Addr().Return(ethAddr)
	mockDev.EXPECT().Transmit(Any, Any, Any).Return(psErr.Error)

	mockIfaceRepo := NewMockIIfaceRepo(ctrl)
	mockIfaceRepo.EXPECT().Lookup(Any, Any).Return(&Iface{
		Family:    V4AddrFamily,
		Unicast:   ParseIP("192.0.2.2"),
		Netmask:   ParseIP("255.255.255.0"),
		Broadcast: ParseIP("192.0.2.255"),
		Dev:       mockDev,
	})
	IfaceRepo = mockIfaceRepo

	packet := Builder.Default()
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.BigEndian, packet)
	dev := &eth.TapDevice{}

	want := psErr.Error
	got := ARP.Receive(buf.Bytes(), dev)
	if got != want {
		t.Errorf("ArpInputHandler() = %s; want %s", got, want)
	}
}

// Success when an ARP packet sent to other destination arrives.
func TestArpInputHandler_6(t *testing.T) {
	ctrl, teardown := SetupArpInputHandlerTest(t)
	defer teardown()

	mockIfaceRepo := NewMockIIfaceRepo(ctrl)
	mockIfaceRepo.EXPECT().Lookup(Any, Any).Return(&Iface{
		Family:    V4AddrFamily,
		Unicast:   ParseIP("192.0.2.2"),
		Netmask:   ParseIP("255.255.255.0"),
		Broadcast: ParseIP("192.0.2.255"),
		Dev:       &eth.TapDevice{},
	})
	IfaceRepo = mockIfaceRepo

	packet := Builder.CustomTPA(ArpProtoAddr{192, 168, 2, 3})
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.BigEndian, packet)
	dev := &eth.TapDevice{
		Device: eth.Device{
			Name_: "net0",
		},
	}

	want := psErr.OK
	got := ARP.Receive(buf.Bytes(), dev)
	if got != want {
		t.Errorf("ArpInputHandler() = %s; want %s", got, want)
	}
}

func TestRunArpTimer_1(t *testing.T) {
	ctrl, teardown := SetupRunArpTimerTest(t)
	defer teardown()
	defer ARP.cache.Init()

	createdAt, _ := time.Parse(time.RFC3339, "2021-01-01T00:00:00Z")
	m := psTime.NewMockITime(ctrl)
	m.EXPECT().Now().Return(createdAt)
	psTime.Time = m

	pa := ArpProtoAddr{192, 168, 1, 1}
	_ = ARP.cache.Create([eth.EthAddrLen]byte{0x11, 0x12, 0x13, 0x14, 0x15, 0x16}, pa, ArpCacheStateResolved)

	var wg sync.WaitGroup
	ARP.RunTimer(&wg)
	<-ArpCondCh
	ARP.StopTimer()
	wg.Wait()

	got := ARP.cache.GetEntry(pa)
	if got != nil {
		t.Errorf("ARP cache is not expired")
	}
}

func TestRunArpTimer_2(t *testing.T) {
	_, teardown := SetupRunArpTimerTest(t)
	defer teardown()
	defer ARP.cache.Init()

	var wg sync.WaitGroup
	ARP.RunTimer(&wg)
	<-ArpCondCh
	ARP.StopTimer()
	wg.Wait()

	got := ARP.cache.GetEntry(ArpProtoAddr{192, 168, 0, 1})
	if got != nil {
		t.Errorf("ARP cache exists")
	}
}

var Any = gomock.Any()
var Builder = ArpPacketBuilder{}

type ArpPacketBuilder struct{}

func (v ArpPacketBuilder) Default() *ArpPacket {
	return &ArpPacket{
		ArpHdr: ArpHdr{
			HT:     ArpHwTypeEthernet,
			PT:     eth.EthTypeIpv4,
			HAL:    eth.EthAddrLen,
			PAL:    V4AddrLen,
			Opcode: ArpOpRequest,
		},
		SHA: [eth.EthAddrLen]byte{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff},
		SPA: ArpProtoAddr{192, 0, 2, 1},
		THA: [eth.EthAddrLen]byte{0x11, 0x12, 0x13, 0x14, 0x15, 0x16},
		TPA: ArpProtoAddr{192, 0, 2, 2},
	}
}

func (v ArpPacketBuilder) CustomHT(ht ArpHwType) (packet *ArpPacket) {
	packet = v.Default()
	packet.HT = ht
	return
}

func (v ArpPacketBuilder) CustomPT(pt eth.EthType) (packet *ArpPacket) {
	packet = v.Default()
	packet.PT = pt
	return
}

func (v ArpPacketBuilder) CustomTPA(tpa ArpProtoAddr) (packet *ArpPacket) {
	packet = v.Default()
	packet.TPA = tpa
	return
}

func SetupArpInputHandlerTest(t *testing.T) (ctrl *gomock.Controller, teardown func()) {
	psLog.DisableOutput()
	ctrl = gomock.NewController(t)
	backupIfaceRepo := IfaceRepo
	backupTime := psTime.Time

	teardown = func() {
		IfaceRepo = backupIfaceRepo
		psTime.Time = backupTime
		ctrl.Finish()
		psLog.EnableOutput()
	}

	return
}

func SetupRunArpTimerTest(t *testing.T) (ctrl *gomock.Controller, teardown func()) {
	psLog.DisableOutput()
	ctrl = gomock.NewController(t)
	backupTime := psTime.Time

	teardown = func() {
		psTime.Time = backupTime
		ctrl.Finish()
		psLog.EnableOutput()
	}

	return
}
