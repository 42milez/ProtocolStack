package cli

import (
	"fmt"
	psErr "github.com/42milez/ProtocolStack/src/error"
	psLog "github.com/42milez/ProtocolStack/src/log"
	"github.com/spf13/cobra"
	"syscall"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "a simple echo server",
	Run: func(cmd *cobra.Command, args []string) {
		if err := setup(); err != psErr.OK {
			psLog.F("initialization failed")
		}
		for {
			sig := <-sigCh
			psLog.I(fmt.Sprintf("signal: %s", sig))
			if sig == syscall.SIGINT || sig == syscall.SIGTERM {
				stopServices()
				return
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
