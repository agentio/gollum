package cmd

import (
	"github.com/agentio/gollum/cmd/atproto"
	"github.com/agentio/gollum/cmd/generate"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gollum",
		Short: "my precious",
	}
	cmd.AddCommand(atproto.Cmd())
	cmd.AddCommand(generate.Cmd())
	return cmd
}
