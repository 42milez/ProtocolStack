package arp

import (
	"bytes"
	"encoding/binary"
	psErr "github.com/42milez/ProtocolStack/src/error"
	"github.com/42milez/ProtocolStack/src/eth"
	psLog "github.com/42milez/ProtocolStack/src/log"
	"github.com/42milez/ProtocolStack/src/mw"
	"github.com/42milez/ProtocolStack/src/repo"
	psTime "github.com/42milez/ProtocolStack/src/time"
	"github.com/golang/mock/gomock"
	"sync"
	"testing"
	"time"
)

func TestArpInputHandler_1(t *testing.T) {
	ctrl, teardown := SetupArpInputHandlerTest(t)
	defer teardown()

	ethAddr := mw.Addr{0x11, 0x12, 0x13, 0x14, 0x15, 0x16}
	mockDev := mw.NewMockIDevice(ctrl)
	mockDev.EXPECT().Addr().Return(ethAddr)
	mockDev.EXPECT().Transmit(any, any, any).Return(psErr.OK)

	mockIfaceRepo := repo.NewMockIIfaceRepo(ctrl)
	mockIfaceRepo.EXPECT().Lookup(any, any).Return(&mw.Iface{
		Family:    mw.V4AddrFamily,
		Unicast:   mw.ParseIP("192.0.2.2"),
		Netmask:   mw.ParseIP("255.255.255.0"),
		Broadcast: mw.ParseIP("192.0.2.255"),
		Dev:       mockDev,
	})
	repo.IfaceRepo = mockIfaceRepo

	packet := Builder.Default()
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.BigEndian, packet)
	dev := &eth.TapDevice{
		Device: mw.Device{
			Type_: mw.DevTypeEthernet,
			Name_: "net0",
			Addr_: ethAddr,
			Flag_: mw.DevFlagBroadcast | mw.DevFlagNeedArp,
			MTU_:  mw.PayloadLenMax,
			Priv_: mw.Privilege{
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

	packet := Builder.CustomHT(HwType(0xffff))
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.BigEndian, packet)
	dev := &eth.TapDevice{}

	want := psErr.InvalidPacket
	got := ARP.Receive(buf.Bytes(), dev)
	if got != want {
		t.Errorf("ArpInputHandler() = %s; want %s", got, want)
	}

	packet = Builder.CustomPT(mw.EthType(0xffff))
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

	mockIfaceRepo := repo.NewMockIIfaceRepo(ctrl)
	mockIfaceRepo.EXPECT().Lookup(any, any).Return(nil)
	repo.IfaceRepo = mockIfaceRepo

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

	ethAddr := mw.Addr{0x11, 0x12, 0x13, 0x14, 0x15, 0x16}
	mockDev := mw.NewMockIDevice(ctrl)
	mockDev.EXPECT().Addr().Return(ethAddr)
	mockDev.EXPECT().Transmit(any, any, any).Return(psErr.Error)

	mockIfaceRepo := repo.NewMockIIfaceRepo(ctrl)
	mockIfaceRepo.EXPECT().Lookup(any, any).Return(&mw.Iface{
		Family:    mw.V4AddrFamily,
		Unicast:   mw.ParseIP("192.0.2.2"),
		Netmask:   mw.ParseIP("255.255.255.0"),
		Broadcast: mw.ParseIP("192.0.2.255"),
		Dev:       mockDev,
	})
	repo.IfaceRepo = mockIfaceRepo

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

	mockIfaceRepo := repo.NewMockIIfaceRepo(ctrl)
	mockIfaceRepo.EXPECT().Lookup(any, any).Return(&mw.Iface{
		Family:    mw.V4AddrFamily,
		Unicast:   mw.ParseIP("192.0.2.2"),
		Netmask:   mw.ParseIP("255.255.255.0"),
		Broadcast: mw.ParseIP("192.0.2.255"),
		Dev:       &eth.TapDevice{},
	})
	repo.IfaceRepo = mockIfaceRepo

	packet := Builder.CustomTPA(ArpProtoAddr{192, 168, 2, 3})
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.BigEndian, packet)
	dev := &eth.TapDevice{
		Device: mw.Device{
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
	defer cache.Init()

	createdAt, _ := time.Parse(time.RFC3339, "2021-01-01T00:00:00Z")
	m := psTime.NewMockITime(ctrl)
	m.EXPECT().Now().Return(createdAt)
	psTime.Time = m

	pa := ArpProtoAddr{192, 168, 1, 1}
	_ = cache.Create(mw.Addr{0x11, 0x12, 0x13, 0x14, 0x15, 0x16}, pa, cacheStatusResolved)

	var wg sync.WaitGroup
	ARP.RunTimer(&wg)
	<-ArpCondCh
	ARP.StopTimer()
	wg.Wait()

	got := cache.GetEntry(pa)
	if got != nil {
		t.Errorf("ARP cache is not expired")
	}
}

func TestRunArpTimer_2(t *testing.T) {
	_, teardown := SetupRunArpTimerTest(t)
	defer teardown()
	defer cache.Init()

	var wg sync.WaitGroup
	ARP.RunTimer(&wg)
	<-ArpCondCh
	ARP.StopTimer()
	wg.Wait()

	got := cache.GetEntry(ArpProtoAddr{192, 168, 0, 1})
	if got != nil {
		t.Errorf("ARP cache exists")
	}
}

var Builder = ArpPacketBuilder{}
var any = gomock.Any()

type ArpPacketBuilder struct{}

func (v ArpPacketBuilder) Default() *Packet {
	return &Packet{
		Hdr: Hdr{
			HT:     ArpHwTypeEthernet,
			PT:     mw.IPv4,
			HAL:    mw.AddrLen,
			PAL:    mw.V4AddrLen,
			Opcode: ArpOpRequest,
		},
		SHA: mw.Addr{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff},
		SPA: ArpProtoAddr{192, 0, 2, 1},
		THA: mw.Addr{0x11, 0x12, 0x13, 0x14, 0x15, 0x16},
		TPA: ArpProtoAddr{192, 0, 2, 2},
	}
}

func (v ArpPacketBuilder) CustomHT(ht HwType) (packet *Packet) {
	packet = v.Default()
	packet.HT = ht
	return
}

func (v ArpPacketBuilder) CustomPT(pt mw.EthType) (packet *Packet) {
	packet = v.Default()
	packet.PT = pt
	return
}

func (v ArpPacketBuilder) CustomTPA(tpa ArpProtoAddr) (packet *Packet) {
	packet = v.Default()
	packet.TPA = tpa
	return
}

func SetupArpInputHandlerTest(t *testing.T) (ctrl *gomock.Controller, teardown func()) {
	psLog.DisableOutput()
	ctrl = gomock.NewController(t)
	backupIfaceRepo := repo.IfaceRepo
	backupTime := psTime.Time

	teardown = func() {
		repo.IfaceRepo = backupIfaceRepo
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
