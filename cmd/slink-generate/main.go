// The slink-generate command contains subcommands that can be used to regenerate slink.
package main

import (
	"os"

	"github.com/agentio/slink/cmd/internal/generate/call"
	"github.com/agentio/slink/cmd/internal/generate/check"
	"github.com/agentio/slink/cmd/internal/generate/lint"
	"github.com/agentio/slink/cmd/internal/generate/xrpc"
	"github.com/spf13/cobra"
)

func main() {
	if err := cmd().Execute(); err != nil {
		os.Exit(1)
	}
}

func cmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "slink-generate",
		Short: "The code generation subcommands of slink.",
	}
	cmd.AddCommand(lint.Cmd())
	cmd.AddCommand(call.Cmd())
	cmd.AddCommand(check.Cmd())
	cmd.AddCommand(xrpc.Cmd())
	return cmd
}
