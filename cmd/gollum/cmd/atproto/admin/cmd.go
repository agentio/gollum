package admin

import (
	"github.com/agentio/gollum/cmd/atproto/admin/getinvitecodes"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "admin",
		Short: "admin subcommands",
	}
	cmd.AddCommand(getinvitecodes.Cmd())
	return cmd
}
