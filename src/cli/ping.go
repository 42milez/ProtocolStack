package cli

import (
	"errors"
	"fmt"
	psErr "github.com/42milez/ProtocolStack/src/error"
	psLog "github.com/42milez/ProtocolStack/src/log"
	"github.com/42milez/ProtocolStack/src/mw"
	"github.com/42milez/ProtocolStack/src/net/icmp"
	"github.com/spf13/cobra"
	"syscall"
)

var count uint32
var dst string

var pingCmd = &cobra.Command{
	Use:   "ping <destination> [flags]",
	Short: "a simple ping command",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires a destination")
		}
		dst = args[0]
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if err := setup(); err != psErr.OK {
			psLog.F("initialization failed")
		}
		send()
		for {
			sig := <-sigCh
			psLog.I(fmt.Sprintf("signal: %s", sig))
			if sig == syscall.SIGINT || sig == syscall.SIGTERM {
				stopServices()
				break
			}
		}
	},
}

func send()  {
	id := 0
	seq := 0
	msg := &mw.IcmpTxMessage{
		Type: icmp.Echo,
		Code: 0,
		Content: uint32(id)<<16 | uint32(seq),
		Payload: []byte{1,2,3,4},
		Src: mw.IP{0, 0, 0, 0},
		Dst: mw.ParseIP(dst),
	}
	mw.IcmpTxCh <- msg
}

func receive() {

}

func init() {
	rootCmd.AddCommand(pingCmd)
	pingCmd.PersistentFlags().Uint32VarP(&count, "count", "c", 0, "stop after <count> replies")
}
