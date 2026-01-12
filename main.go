package main

import (
	"fmt"
	"os"

	"github.com/agentio/gollum/pkg/lexica"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

func main() {
	if err := cmd().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func cmd() *cobra.Command {
	var logLevel string
	var cmd = &cobra.Command{
		Use: "gollum",
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			ll, err := log.ParseLevel(logLevel)
			if err != nil {
				return err
			}
			log.SetLevel(ll)
			lex := &lexica.Lexica{}
			rootDir := "lexicons"
			err = lex.LoadTree(rootDir)
			if err != nil {
				return err
			}
			err = lex.Generate("api")
			if err != nil {
				return err
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&logLevel, "log-level", "l", "info", "log level (debug, info, warn, error, fatal)")
	return cmd
}
