package cli

import (
	"errors"
	"fmt"
	psErr "github.com/42milez/ProtocolStack/src/error"
	psLog "github.com/42milez/ProtocolStack/src/log"
	"github.com/42milez/ProtocolStack/src/mw"
	"github.com/42milez/ProtocolStack/src/net/icmp"
	"github.com/spf13/cobra"
	"math/rand"
	"syscall"
	"time"
)

var count int
var dst string
var provider provider_

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
		var nReplied int
		for {
			send(provider.ID(), provider.SeqNumber(), provider.Payload())
			select {
			case sig := <-sigCh:
				psLog.I(fmt.Sprintf("signal: %s", sig))
				if sig == syscall.SIGINT || sig == syscall.SIGTERM {
					stopServices()
					return
				}
			case reply := <-icmp.ReplyQueue:
				psLog.I(fmt.Sprintf("id, seq: %d, %d", reply.ID, reply.Seq))
				nReplied += 1
				time.Sleep(time.Second)
			case letter := <-mw.IcmpDeadLetterQueue:
				psLog.W("dead letter detected")
				time.Sleep(100 * time.Millisecond)
				id := uint16(letter.Payload[5])<<8 | uint16(letter.Payload[6])
				seq := uint16(letter.Payload[7])<<8 | uint16(letter.Payload[8])
				send(id, seq, letter.Payload[icmp.HdrLen:])
			}

			if nReplied == count {
				return
			}
		}
	},
}

type provider_ struct {
	seqNumber uint16
}

func (p *provider_) ID() (id uint16) {
	id = uint16(rand.Int())
	return
}

func (p *provider_) SeqNumber() (seq uint16) {
	seq = p.seqNumber
	seq += 1
	return
}

func (p *provider_) Payload() (payload []byte) {
	var size = 56
	payload = make([]byte, size)
	for i := 0; i < size; i += 4 {
		v := rand.Uint32()
		payload[i] = uint8(v)
		payload[i+1] = uint8(v >> 8)
		payload[i+2] = uint8(v >> 16)
		payload[i+3] = uint8(v >> 24)
	}
	return
}

func send(id uint16, seq uint16, payload []byte) {
	msg := &mw.IcmpTxMessage{
		Type:    icmp.Echo,
		Code:    0,
		Content: uint32(id)<<16 | uint32(seq),
		Payload: payload,
		Src:     mw.IP{192, 0, 2, 2},
		Dst:     mw.ParseIP(dst),
	}
	mw.IcmpTxCh <- msg
}

func init() {
	provider = provider_{}
	rootCmd.AddCommand(pingCmd)
	pingCmd.PersistentFlags().IntVarP(&count, "count", "c", 0, "stop after <count> replies")
}
