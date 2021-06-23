package cli

import (
	"bytes"
	"errors"
	"fmt"
	psErr "github.com/42milez/ProtocolStack/src/error"
	psLog "github.com/42milez/ProtocolStack/src/log"
	"github.com/42milez/ProtocolStack/src/mw"
	"github.com/42milez/ProtocolStack/src/net/icmp"
	"github.com/spf13/cobra"
	"syscall"
	"time"
)

var provider provider_
var count int
var dst string
var skipNextRequest bool

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
			if !skipNextRequest {
				id := provider.ID()
				seq := provider.SeqNum()
				send(id, seq, provider.Payload())
				psLog.I(fmt.Sprintf("icmp packet sent:     seq=%d, id=%d", seq, id))
			}

			select {
			case sig := <-sigCh:
				psLog.I(fmt.Sprintf("signal: %s", sig))
				if sig == syscall.SIGINT || sig == syscall.SIGTERM {
					stopServices()
					return
				}
			case reply := <-icmp.ReplyQueue:
				handleReply(reply)
				nReplied += 1
				skipNextRequest = false
			case letter := <-mw.IcmpDeadLetterQueue:
				time.Sleep(100 * time.Millisecond)
				handleDeadLetter(letter)
				skipNextRequest = true
			}

			if nReplied == count {
				return
			}

			time.Sleep(time.Second)
		}
	},
}

type provider_ struct {
	seqNumber uint16
}

func (p *provider_) ID() (id uint16) {
	id = mw.RandU16()
	return
}

func (p *provider_) SeqNum() (seq uint16) {
	seq = p.seqNumber
	p.seqNumber += 1
	return
}

func (p *provider_) Payload() (payload []byte) {
	var size = 56
	payload = make([]byte, size)
	for i := 0; i < size; i += 4 {
		payload[i] = mw.RandU8()
		payload[i+1] = mw.RandU8()
		payload[i+2] = mw.RandU8()
		payload[i+3] = mw.RandU8()
	}
	return
}

func send(id uint16, seq uint16, payload []byte) {
	msg := &mw.IcmpTxMessage{
		Type:    icmp.Echo,
		Code:    0,
		Content: uint32(id)<<16 | uint32(seq),
		Data:    payload,
		Src:     mw.IP{192, 0, 2, 2},
		Dst:     mw.ParseIP(dst),
	}
	mw.IcmpTxCh <- msg
}

func handleDeadLetter(letter *mw.IcmpQueueEntry) {
	hdr, _ := icmp.ReadHeader(bytes.NewBuffer(letter.Packet))
	id, seq := icmp.SplitContent(hdr.Content)
	send(id, seq, letter.Packet[icmp.HdrLen:])
	psLog.W(fmt.Sprintf("icmp packet sent (dead letter): id=%d, seq=%d", id, seq))
}

func handleReply(reply *icmp.Reply) {
	psLog.I(fmt.Sprintf("icmp packet received: seq=%d, id=%d", reply.Seq, reply.ID))
}

func init() {
	provider = provider_{}
	rootCmd.AddCommand(pingCmd)
	pingCmd.PersistentFlags().IntVarP(&count, "count", "c", 0, "stop after <count> replies")
}
