package main

import (
	"os"

	"github.com/agentio/slink/gen/slink"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := slink.Cmd()
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.SetHelpCommand(&cobra.Command{
		Use:    "",
		Hidden: true,
	})
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
