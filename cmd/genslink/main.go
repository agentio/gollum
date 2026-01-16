package main

import (
	"os"

	"github.com/agentio/slink/cli"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := cli.Cmd()
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.SetHelpCommand(&cobra.Command{
		Use:    "",
		Hidden: true,
	})
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
