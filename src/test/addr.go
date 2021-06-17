package test

import "github.com/42milez/ProtocolStack/src/mw"

var EthAddrBuilder ethAddrBuilder

type ethAddrBuilder struct{}

func (ethAddrBuilder) Default() mw.EthAddr {
	return mw.EthAddr{11, 12, 13, 14, 15, 16}
}

func init() {
	EthAddrBuilder = ethAddrBuilder{}
}
