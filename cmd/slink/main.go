package main

import (
	"os"
	"strings"

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
		Long: strings.Join(
			[]string{
				``,
				`"Perhaps we’ve shaken him off at last, the miserable slinker!"`,
				``,
				`A tool for working with the AT Protocol.`,
			}, "\n"),
	}
	cmd.AddCommand(call.Cmd())
	return cmd
}
