package cli

import (
	"fmt"
	psErr "github.com/42milez/ProtocolStack/src/error"
	psLog "github.com/42milez/ProtocolStack/src/log"
	"github.com/42milez/ProtocolStack/src/mw"
	"github.com/42milez/ProtocolStack/src/net/tcp"
	"github.com/spf13/cobra"
	"os"
	"syscall"
)

const host = "192.0.2.2"
const port = 12345

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "a simple echo server",
	Run: func(cmd *cobra.Command, args []string) {
		if err := setup(); err != psErr.OK {
			psLog.F("initialization failed")
		}

		id, err := tcp.Open()
		if err != psErr.OK {
			psLog.E(fmt.Sprintf("can't open socket: %s", err))
			os.Exit(1)
		}

		local := tcp.EndPoint{
			Addr: mw.ParseIP(host).ToV4(),
			Port: port,
		}
		if err := tcp.Bind(id, local); err != psErr.OK {
			psLog.E(fmt.Sprintf("can't bind: %s", err))
			os.Exit(1)
		}

		if err := tcp.Listen(id, 1); err != psErr.OK {
			psLog.E(fmt.Sprintf("can't listen: %s", err))
			os.Exit(1)
		}

		_, foreign, err := tcp.Accept(id)
		if err != psErr.OK {
			psLog.E(fmt.Sprintf("can't accept: %s", err))
			os.Exit(1)
		}

		psLog.I("connection accepted",
			fmt.Sprintf("Host: %s", foreign.Addr.String()),
			fmt.Sprintf("Port: %d", foreign.Port))

		for {
			sig := <-sigCh
			psLog.I(fmt.Sprintf("signal: %s", sig))
			if sig == syscall.SIGINT || sig == syscall.SIGTERM {
				stopServices()
				break
			}
		}
		psLog.I("stopped")
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
