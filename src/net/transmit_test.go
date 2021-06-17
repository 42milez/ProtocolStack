package net

import (
	psErr "github.com/42milez/ProtocolStack/src/error"
	psLog "github.com/42milez/ProtocolStack/src/log"
	"github.com/42milez/ProtocolStack/src/mw"
	"github.com/42milez/ProtocolStack/src/test"
	"github.com/golang/mock/gomock"
	"testing"
)

func TestTransmit(t *testing.T) {
	ctrl, teardown := setupTransmitTest(t)
	defer teardown()

	devMock := mw.NewMockIDevice(ctrl)
	devMock.EXPECT().IsUp().Return(true)
	devMock.EXPECT().MTU().Return(uint16(mw.EthPayloadLenMax))
	devMock.EXPECT().Transmit(any, any, any).Return(psErr.OK)

	iface := test.IfaceBuilder.Default()
	iface.Dev = devMock

	payload := test.PayloadBuilder.Default()
	dstEthAddr := test.EthAddrBuilder.Default()

	want := psErr.OK
	got := Transmit(dstEthAddr, payload, mw.IPv4, iface)

	if got != want {
		t.Errorf("Transmit() = %s; want %s", got, want)
	}
}

var any = gomock.Any()

func setupTransmitTest(t *testing.T) (ctrl *gomock.Controller, teardown func()) {
	ctrl = gomock.NewController(t)
	psLog.DisableOutput()
	reset := func() {
		psLog.EnableOutput()
	}
	teardown = func() {
		ctrl.Finish()
		reset()
	}
	return
}
