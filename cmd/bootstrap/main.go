package main

import (
	"os"

	"github.com/agentio/slink/cmd/bootstrap/call"
	"github.com/agentio/slink/cmd/bootstrap/check"
	"github.com/agentio/slink/cmd/bootstrap/lint"
	"github.com/agentio/slink/cmd/bootstrap/xrpc"
	"github.com/spf13/cobra"
)

func main() {
	if err := cmd().Execute(); err != nil {
		os.Exit(1)
	}
}

func cmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use: "bootstrap",
	}
	cmd.AddCommand(lint.Cmd())
	cmd.AddCommand(call.Cmd())
	cmd.AddCommand(check.Cmd())
	cmd.AddCommand(xrpc.Cmd())
	return cmd
}
