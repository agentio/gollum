package cmd

import (
	"github.com/agentio/gollum/cmd/atproto"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pdsboss",
		Short: "PDS Boss",
	}
	cmd.AddCommand(atproto.Cmd())
	return cmd
}
