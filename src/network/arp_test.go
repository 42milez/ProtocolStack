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

var Any = gomock.Any()

func setupArpTest(t *testing.T) (ctrl *gomock.Controller, teardown func()) {
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

func TestArpInputHandler(t *testing.T) {
	ctrl, teardown := setupArpTest(t)
	defer teardown()

	ethAddr := ethernet.EthAddr{11, 22, 33, 44, 55, 66}

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

	packet := &ArpPacket{
		ArpHdr: ArpHdr{
			HT:     ArpHwTypeEthernet,
			PT:     ethernet.EthTypeIpv4,
			HAL:    ethernet.EthAddrLen,
			PAL:    V4AddrLen,
			Opcode: ArpOpRequest,
		},
		THA: ethAddr,
		TPA: ArpProtoAddr{192, 0, 2, 2},
	}

	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.BigEndian, packet)

	want := psErr.OK
	got := ArpInputHandler(buf.Bytes(), dev)
	if got != want {
		t.Errorf("ArpInputHandler() = %s; want %s", got, want)
	}
}

func TestRunArpTimer_1(t *testing.T) {
	ctrl, teardown := setupArpTest(t)
	defer teardown()
	defer cache.Init()

	createdAt, _ := time.Parse(time.RFC3339, "2021-01-01T00:00:00Z")
	m := psTime.NewMockITime(ctrl)
	m.EXPECT().Now().Return(createdAt)
	psTime.Time = m

	ha := ethernet.EthAddr{0x11, 0x22, 0x33, 0x44, 0x55, 0x66}
	pa := ArpProtoAddr{192, 168, 1, 1}
	state := ArpCacheStateResolved
	_ = cache.Add(ha, pa, state)

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
