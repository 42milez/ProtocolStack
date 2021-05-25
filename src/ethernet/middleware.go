package ethernet

var RxCh chan *Packet // channel for receiving packets
var TxCh chan *Packet // channel for sending packets

const RxChBufSize = 10
const TxChBufSize = 10

type Packet struct {
	Type    EthType
	Payload []byte
}

func init() {
	RxCh = make(chan *Packet, RxChBufSize)
	TxCh = make(chan *Packet, TxChBufSize)
}
