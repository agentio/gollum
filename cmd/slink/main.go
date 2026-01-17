package main

import (
	"os"

	"github.com/agentio/slink/gen/call"
	"github.com/spf13/cobra"
)

func main() {
	if err := cmd().Execute(); err != nil {
		os.Exit(1)
	}
}

func cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "slink",
	}
	cmd.AddCommand(call.Cmd())
	return cmd
}
