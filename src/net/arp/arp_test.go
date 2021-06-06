package arp

import (
	"bytes"
	"encoding/binary"
	psErr "github.com/42milez/ProtocolStack/src/error"
	psLog "github.com/42milez/ProtocolStack/src/log"
	"github.com/42milez/ProtocolStack/src/mw"
	eth2 "github.com/42milez/ProtocolStack/src/net/eth"
	"github.com/42milez/ProtocolStack/src/repo"
	psTime "github.com/42milez/ProtocolStack/src/time"
	"github.com/golang/mock/gomock"
	"testing"
)

func TestReceive_1(t *testing.T) {
	ctrl, teardown := SetupReceiveTest(t)
	defer teardown()

	ethAddr := mw.EthAddr{0x11, 0x12, 0x13, 0x14, 0x15, 0x16}
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

	packet := builder.Default()
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.BigEndian, packet)
	dev := &eth2.TapDevice{
		Device: mw.Device{
			Type_: mw.EthernetDevice,
			Name_: "net0",
			Addr_: ethAddr,
			Flag_: mw.BroadcastFlag | mw.NeedArpFlag,
			MTU_:  mw.EthPayloadLenMax,
			Priv_: mw.Privilege{
				FD:   3,
				Name: "tap0",
			},
		},
	}

	want := psErr.OK
	got := Receive(buf.Bytes(), dev)
	if got != want {
		t.Errorf("Receive() = %s; want %s", got, want)
	}
}

// Fail when ARP packet is too short.
func TestReceive_2(t *testing.T) {
	_, teardown := SetupReceiveTest(t)
	defer teardown()

	var packet []byte
	dev := &eth2.TapDevice{}

	want := psErr.InvalidPacket
	got := Receive(packet, dev)
	if got != want {
		t.Errorf("Receive() = %s; want %s", got, want)
	}
}

// Fail when hardware type is invalid.
// Fail when protocol type is invalid.
func TestReceive_3(t *testing.T) {
	_, teardown := SetupReceiveTest(t)
	defer teardown()

	packet := builder.CustomHT(HwType(0xffff))
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.BigEndian, packet)
	dev := &eth2.TapDevice{}

	want := psErr.InvalidPacket
	got := Receive(buf.Bytes(), dev)
	if got != want {
		t.Errorf("Receive() = %s; want %s", got, want)
	}

	packet = builder.CustomPT(mw.EthType(0xffff))
	buf = new(bytes.Buffer)
	_ = binary.Write(buf, binary.BigEndian, packet)
	dev = &eth2.TapDevice{}

	want = psErr.InvalidPacket
	got = Receive(buf.Bytes(), dev)
	if got != want {
		t.Errorf("Receive() = %s; want %s", got, want)
	}
}

// Fail when Iface is not found.
func TestReceive_4(t *testing.T) {
	ctrl, teardown := SetupReceiveTest(t)
	defer teardown()

	mockIfaceRepo := repo.NewMockIIfaceRepo(ctrl)
	mockIfaceRepo.EXPECT().Lookup(any, any).Return(nil)
	repo.IfaceRepo = mockIfaceRepo

	packet := builder.Default()
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.BigEndian, packet)
	dev := &eth2.TapDevice{}

	want := psErr.InterfaceNotFound
	got := Receive(buf.Bytes(), dev)
	if got != want {
		t.Errorf("Receive() = %s; want %s", got, want)
	}
}

// Fail when SendReply returns error.
func TestReceive_5(t *testing.T) {
	ctrl, teardown := SetupReceiveTest(t)
	defer teardown()

	ethAddr := mw.EthAddr{0x11, 0x12, 0x13, 0x14, 0x15, 0x16}
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

	packet := builder.Default()
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.BigEndian, packet)
	dev := &eth2.TapDevice{}

	want := psErr.Error
	got := Receive(buf.Bytes(), dev)
	if got != want {
		t.Errorf("Receive() = %s; want %s", got, want)
	}
}

// Success when an ARP packet sent to other destination arrives.
func TestReceive_6(t *testing.T) {
	ctrl, teardown := SetupReceiveTest(t)
	defer teardown()

	mockIfaceRepo := repo.NewMockIIfaceRepo(ctrl)
	mockIfaceRepo.EXPECT().Lookup(any, any).Return(&mw.Iface{
		Family:    mw.V4AddrFamily,
		Unicast:   mw.ParseIP("192.0.2.2"),
		Netmask:   mw.ParseIP("255.255.255.0"),
		Broadcast: mw.ParseIP("192.0.2.255"),
		Dev:       &eth2.TapDevice{},
	})
	repo.IfaceRepo = mockIfaceRepo

	packet := builder.CustomTPA(mw.V4Addr{192, 168, 2, 3})
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.BigEndian, packet)
	dev := &eth2.TapDevice{
		Device: mw.Device{
			Name_: "net0",
		},
	}

	want := psErr.OK
	got := Receive(buf.Bytes(), dev)
	if got != want {
		t.Errorf("Receive() = %s; want %s", got, want)
	}
}

func TestHwType_String(t *testing.T) {
	want := hwTypes[Ethernet]
	got := Ethernet.String()
	if got != want {
		t.Errorf("HwType.String() = %s; want %s", got, want)
	}
}

func TestOpcode_String(t *testing.T) {
	want := opCodes[Request]
	got := Request.String()
	if got != want {
		t.Errorf("Opcode.String() = %s; want %s", got, want)
	}
}

var any = gomock.Any()
var builder = &PacketBuilder{}

type PacketBuilder struct{}

func (v PacketBuilder) Default() *Packet {
	return &Packet{
		Hdr: Hdr{
			HT:     Ethernet,
			PT:     mw.IPv4,
			HAL:    mw.EthAddrLen,
			PAL:    mw.V4AddrLen,
			Opcode: Request,
		},
		SHA: mw.EthAddr{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff},
		SPA: mw.V4Addr{192, 0, 2, 1},
		THA: mw.EthAddr{0x11, 0x12, 0x13, 0x14, 0x15, 0x16},
		TPA: mw.V4Addr{192, 0, 2, 2},
	}
}

func (v PacketBuilder) CustomHT(ht HwType) (packet *Packet) {
	packet = v.Default()
	packet.HT = ht
	return
}

func (v PacketBuilder) CustomPT(pt mw.EthType) (packet *Packet) {
	packet = v.Default()
	packet.PT = pt
	return
}

func (v PacketBuilder) CustomTPA(tpa mw.V4Addr) (packet *Packet) {
	packet = v.Default()
	packet.TPA = tpa
	return
}

func SetupReceiveTest(t *testing.T) (ctrl *gomock.Controller, teardown func()) {
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
