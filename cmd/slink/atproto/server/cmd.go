package server

import (
	"github.com/agentio/slink/cmd/slink/atproto/server/createinvitecode"
	"github.com/agentio/slink/cmd/slink/atproto/server/getaccountinvitecodes"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "server subcommands",
	}
	cmd.AddCommand(createinvitecode.Cmd())
	cmd.AddCommand(getaccountinvitecodes.Cmd())
	return cmd
}
