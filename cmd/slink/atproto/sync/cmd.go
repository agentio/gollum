package sync

import (
	"github.com/agentio/slink/cmd/slink/atproto/sync/listrepos"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "sync subcommands",
	}
	cmd.AddCommand(listrepos.Cmd())
	return cmd
}
