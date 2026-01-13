package main

import (
	"os"

	"github.com/agentio/gollum/cmd/gollum/atproto"
	"github.com/agentio/gollum/cmd/gollum/generate"
	"github.com/spf13/cobra"
)

func main() {
	if err := cmd().Execute(); err != nil {
		os.Exit(1)
	}
}
func cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gollum",
		Short: "my precious",
	}
	cmd.AddCommand(atproto.Cmd())
	cmd.AddCommand(generate.Cmd())
	return cmd
}
