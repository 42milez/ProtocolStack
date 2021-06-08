package main

import (
	"github.com/42milez/ProtocolStack/src/cli"
	psLog "github.com/42milez/ProtocolStack/src/log"
)

var release string

func main() {
	cli.Execute()
}

func init() {
	if release == "true" {
		psLog.DisableDebug()
	}
}
