package ethernet

import (
	"fmt"
	"strconv"
	"strings"
)

const EthAddrLen = 6

const EthHeaderSize = 14
const EthFrameSizeMin = 60
const EthFrameSizeMax = 1514
const EthPayloadSizeMin = EthFrameSizeMin - EthHeaderSize
const EthPayloadSizeMax = EthFrameSizeMax - EthHeaderSize

var EthAddrAny = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
var EthAddrBroadcast = []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}

type MAC string

func (mac MAC) Byte() ([]byte, error) {
	t := strings.Split(string(mac), ":")
	p := make([]byte, EthAddrLen)
	for i := 0; i < EthAddrLen; i++ {
		var n uint64
		var err error
		if n, err = strconv.ParseUint(t[i], 16, 8); err != nil {
			return nil, err
		}
		if n > 0xff {
			return nil, fmt.Errorf("invalid MAC address")
		}
		p[i] = byte(n)
	}
	return p, nil
}
