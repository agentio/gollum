package main

import (
	"os"

	"github.com/agentio/slink/cmd/bootstrap/generate"
	"github.com/agentio/slink/cmd/slink/atproto"
	"github.com/spf13/cobra"
)

func main() {
	if err := cmd().Execute(); err != nil {
		os.Exit(1)
	}
}
func cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "slink",
		Short: "my precious",
	}
	cmd.AddCommand(atproto.Cmd())
	cmd.AddCommand(generate.Cmd())
	return cmd
}
