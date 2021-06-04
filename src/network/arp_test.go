package network

import (
	"bytes"
	"encoding/binary"
	psErr "github.com/42milez/ProtocolStack/src/error"
	"github.com/42milez/ProtocolStack/src/ethernet"
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

	ethAddr := ethernet.EthAddr{0x11, 0x12, 0x13, 0x14, 0x15, 0x16}
	mockDev := ethernet.NewMockIDevice(ctrl)
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
	dev := &ethernet.TapDevice{
		Device: ethernet.Device{
			Type_: ethernet.DevTypeEthernet,
			Name_: "net0",
			Addr_: ethAddr,
			Flag_: ethernet.DevFlagBroadcast | ethernet.DevFlagNeedArp,
			MTU_:  ethernet.EthPayloadLenMax,
			Priv_: ethernet.Privilege{
				FD:   3,
				Name: "tap0",
			},
		},
	}

	want := psErr.OK
	got := ArpInputHandler(buf.Bytes(), dev)
	if got != want {
		t.Errorf("ArpInputHandler() = %s; want %s", got, want)
	}
}

// Fail when ARP packet is too short.
func TestArpInputHandler_2(t *testing.T) {
	_, teardown := SetupArpInputHandlerTest(t)
	defer teardown()

	var packet []byte
	dev := &ethernet.TapDevice{}

	want := psErr.InvalidPacket
	got := ArpInputHandler(packet, dev)
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
	dev := &ethernet.TapDevice{}

	want := psErr.InvalidPacket
	got := ArpInputHandler(buf.Bytes(), dev)
	if got != want {
		t.Errorf("ArpInputHandler() = %s; want %s", got, want)
	}

	packet = Builder.CustomPT(ethernet.EthType(0xffff))
	buf = new(bytes.Buffer)
	_ = binary.Write(buf, binary.BigEndian, packet)
	dev = &ethernet.TapDevice{}

	want = psErr.InvalidPacket
	got = ArpInputHandler(buf.Bytes(), dev)
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
	dev := &ethernet.TapDevice{}

	want := psErr.InterfaceNotFound
	got := ArpInputHandler(buf.Bytes(), dev)
	if got != want {
		t.Errorf("ArpInputHandler() = %s; want %s", got, want)
	}
}

// Fail when arpReply() returns error.
func TestArpInputHandler_5(t *testing.T) {
	ctrl, teardown := SetupArpInputHandlerTest(t)
	defer teardown()

	ethAddr := ethernet.EthAddr{0x11, 0x12, 0x13, 0x14, 0x15, 0x16}
	mockDev := ethernet.NewMockIDevice(ctrl)
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
	dev := &ethernet.TapDevice{}

	want := psErr.Error
	got := ArpInputHandler(buf.Bytes(), dev)
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
		Dev:       &ethernet.TapDevice{},
	})
	IfaceRepo = mockIfaceRepo

	packet := Builder.CustomTPA(ArpProtoAddr{192, 168, 2, 3})
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.BigEndian, packet)
	dev := &ethernet.TapDevice{
		Device: ethernet.Device{
			Name_: "net0",
		},
	}

	want := psErr.OK
	got := ArpInputHandler(buf.Bytes(), dev)
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
	_ = cache.Add(ethernet.EthAddr{0x11, 0x12, 0x13, 0x14, 0x15, 0x16}, pa, ArpCacheStateResolved)

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
	_, teardown := SetupRunArpTimerTest(t)
	defer teardown()
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

var Any = gomock.Any()
var Builder = ArpPacketBuilder{}

type ArpPacketBuilder struct {}

func (v ArpPacketBuilder) Default() *ArpPacket {
	return &ArpPacket{
		ArpHdr: ArpHdr{
			HT:     ArpHwTypeEthernet,
			PT:     ethernet.EthTypeIpv4,
			HAL:    ethernet.EthAddrLen,
			PAL:    V4AddrLen,
			Opcode: ArpOpRequest,
		},
		SHA: ethernet.EthAddr{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff},
		SPA: ArpProtoAddr{192, 0, 2, 1},
		THA: ethernet.EthAddr{0x11, 0x12, 0x13, 0x14, 0x15, 0x16},
		TPA: ArpProtoAddr{192, 0, 2, 2},
	}
}

func (v ArpPacketBuilder) CustomHT(ht ArpHwType) (packet *ArpPacket) {
	packet = v.Default()
	packet.HT = ht
	return
}

func (v ArpPacketBuilder) CustomPT(pt ethernet.EthType) (packet *ArpPacket) {
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
