package cli

import (
	"errors"
	"fmt"
	psErr "github.com/42milez/ProtocolStack/src/error"
	psLog "github.com/42milez/ProtocolStack/src/log"
	"github.com/spf13/cobra"
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
		fmt.Printf("send icmp request to %s\n", dst)
	},
}

func init() {
	rootCmd.AddCommand(pingCmd)
	pingCmd.PersistentFlags().Uint32VarP(&count, "count", "c", 0, "stop after <count> replies")
}
