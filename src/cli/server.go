package cli

import (
	"errors"
	"fmt"
	psErr "github.com/42milez/ProtocolStack/src/error"
	psLog "github.com/42milez/ProtocolStack/src/log"
	"github.com/42milez/ProtocolStack/src/mw"
	"github.com/42milez/ProtocolStack/src/net/tcp"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"syscall"
)

var host string
var port uint16

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "a simple echo server",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("requires host and port")
		}
		host = args[0]
		p, err := strconv.ParseInt(args[1], 10, 0)
		if err != nil {
			return err
		}
		port = uint16(p)
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if err := setup(); err != psErr.OK {
			psLog.F("initialization failed")
		}

		soc, err := tcp.Open()
		if err != psErr.OK {
			psLog.E(fmt.Sprintf("can't open socket: %s", err))
			os.Exit(1)
		}

		local := tcp.EndPoint{
			Addr: mw.ParseIP(host).ToV4(),
			Port: port,
		}
		if err := tcp.Bind(soc, local); err != psErr.OK {
			psLog.E(fmt.Sprintf("can't bind: %s", err))
			os.Exit(1)
		}

		if err := tcp.Listen(soc, 1); err != psErr.OK {
			psLog.E(fmt.Sprintf("can't listen: %s", err))
			os.Exit(1)
		}

		_, foreign, err := tcp.Accept(soc)
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
