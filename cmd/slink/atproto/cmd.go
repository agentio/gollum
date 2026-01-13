package atproto

import (
	"github.com/agentio/slink/cmd/slink/atproto/admin"
	"github.com/agentio/slink/cmd/slink/atproto/server"
	"github.com/agentio/slink/cmd/slink/atproto/sync"
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
