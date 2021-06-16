package ip

import (
	"bytes"
	"encoding/binary"
	psErr "github.com/42milez/ProtocolStack/src/error"
	psLog "github.com/42milez/ProtocolStack/src/log"
	"github.com/42milez/ProtocolStack/src/mw"
	"github.com/42milez/ProtocolStack/src/net/arp"
	"github.com/42milez/ProtocolStack/src/net/eth"
	"github.com/42milez/ProtocolStack/src/repo"
	"github.com/42milez/ProtocolStack/src/worker"
	"github.com/golang/mock/gomock"
	"sync"
	"testing"
)

func TestPacketID_Next(t *testing.T) {
	want := uint16(0)
	got := id.Next()
	if got != want {
		t.Errorf("PacketID.Next() = %d; want %d", got, want)
	}
}

func TestReceive(t *testing.T) {
	_, teardown := setupIpTest(t)
	defer teardown()

	dev := createValidTapDevice()
	iface := createValidIface()
	_ = repo.IfaceRepo.Register(iface, dev)
	packet := createValidIpPacket()

	want := psErr.OK
	got := Receive(packet, dev)

	if got != want {
		t.Errorf("Receive() = %s; want %s", got, want)
	}
}

func TestSend(t *testing.T) {
	ctrl, teardown := setupIpTest(t)
	defer teardown()

	devMock := mw.NewMockIDevice(ctrl)
	devMock.EXPECT().IsUp().Return(true)
	devMock.EXPECT().Name().Return("net0")
	devMock.EXPECT().Flag().Return(mw.BroadcastFlag | mw.NeedArpFlag)
	devMock.EXPECT().MTU().Return(uint16(mw.EthPayloadLenMax)).AnyTimes()
	devMock.EXPECT().Priv().Return(mw.Privilege{FD: 3, Name: "tap0"})
	devMock.EXPECT().Transmit(any, any, any).Return(psErr.OK)

	iface := createValidIface()
	_ = repo.IfaceRepo.Register(iface, devMock)

	arpMock := arp.NewMockIResolver(ctrl)
	arpMock.EXPECT().Resolve(any, any).Return(mw.EthAddr{11, 12, 13, 14, 15, 16}, arp.Complete)
	arp.Resolver = arpMock

	payload := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	src := mw.IP{192, 168, 0, 1}
	dst := mw.IP{192, 168, 0, 2}

	want := psErr.OK
	got := Send(ICMP, payload, src, dst)

	if got != want {
		t.Errorf("Send() = %s; want %s", got, want)
	}
}

func TestStart(t *testing.T) {
	_, teardown := setupIpTest(t)
	defer teardown()

	var wg sync.WaitGroup
	_ = Start(&wg)
	rcvMonMsg := <-rcvMonCh
	sndMonMsg := <-sndMonCh

	if rcvMonMsg.Current != worker.Running || sndMonMsg.Current != worker.Running {
		t.Errorf("Start() failed")
	}
}

func TestStop(t *testing.T) {
	_, teardown := setupIpTest(t)
	defer teardown()

	var wg sync.WaitGroup
	_ = Start(&wg)
	<-rcvMonCh
	<-sndMonCh
	Stop()
	rcvMonMsg := <-rcvMonCh
	sndMonMsg := <-sndMonCh

	if rcvMonMsg.Current != worker.Stopped || sndMonMsg.Current != worker.Stopped {
		t.Errorf("Stop() failed")
	}
}

var any = gomock.Any()

func createValidIface() *mw.Iface {
	return &mw.Iface{
		Family:    mw.V4AddrFamily,
		Unicast:   mw.IP{192, 168, 0, 1},
		Netmask:   mw.IP{255, 255, 255, 0},
		Broadcast: mw.IP{192, 168, 0, 255},
	}
}

func createValidIpPacket() []byte {
	payload := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	hdr := &mw.IpHdr{}
	hdr.VHL = uint8(ipv4<<4) | uint8(HdrLenMin/4)
	hdr.TotalLen = uint16(HdrLenMin + len(payload))
	hdr.ID = 0
	hdr.TTL = 0xff
	hdr.Protocol = ICMP
	hdr.Src = mw.V4Addr{192, 168, 0, 1}
	hdr.Dst = mw.V4Addr{192, 168, 1, 1}

	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, hdr); err != nil {
		return nil
	}
	if err := binary.Write(buf, binary.BigEndian, &payload); err != nil {
		return nil
	}
	packet := buf.Bytes()

	csum := mw.Checksum(packet, 0)
	packet[10] = uint8((csum & 0xff00) >> 8)
	packet[11] = uint8(csum & 0x00ff)

	return packet
}

func createValidTapDevice() *eth.TapDevice {
	return &eth.TapDevice{
		Device: mw.Device{
			Type_: mw.EthernetDevice,
			Name_: "net0",
			MTU_:  mw.EthPayloadLenMax,
			Flag_: mw.BroadcastFlag | mw.NeedArpFlag,
			Addr_: mw.EthAddr{11, 12, 13, 14, 15, 16},
			Priv_: mw.Privilege{FD: 3, Name: "tap0"},
		},
	}
}

func setupIpTest(t *testing.T) (ctrl *gomock.Controller, teardown func()) {
	ctrl = gomock.NewController(t)
	psLog.DisableOutput()
	reset := func() {
		psLog.EnableOutput()
		repo.IfaceRepo.Init()
	}
	teardown = func() {
		ctrl.Finish()
		reset()
	}
	return
}
