package atproto

import (
	"github.com/agentio/gollum/cmd/gollum/atproto/admin"
	"github.com/agentio/gollum/cmd/gollum/atproto/server"
	"github.com/agentio/gollum/cmd/gollum/atproto/sync"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "atproto",
		Short: "atproto subcommands",
	}
	cmd.AddCommand(admin.Cmd())
	cmd.AddCommand(server.Cmd())
	cmd.AddCommand(sync.Cmd())
	return cmd
}
