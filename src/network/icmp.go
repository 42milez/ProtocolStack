package network

import (
	psErr "github.com/42milez/ProtocolStack/src/error"
	psLog "github.com/42milez/ProtocolStack/src/log"
)

func IcmpInputHandler() psErr.E {
	psLog.I("ICMP handler was called")
	return psErr.OK
}
