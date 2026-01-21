package check

import (
	"github.com/agentio/slink/pkg/lexica"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	var input string
	var output string
	var logLevel string
	var cmd = &cobra.Command{
		Use:   "check",
		Short: "Generate a command-line interface to check records in a directory of lexicons",
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			ll, err := log.ParseLevel(logLevel)
			if err != nil {
				return err
			}
			log.SetLevel(ll)
			catalog := lexica.NewCatalog()
			if err = catalog.Load(input, false /* skip lint */); err != nil {
				return err
			}
			err = catalog.GenerateCheckCommands(output)
			if err != nil {
				return err
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&input, "input", "i", "lexicons", "input directory")
	cmd.Flags().StringVarP(&output, "output", "o", "gen/check", "output directory")
	cmd.Flags().StringVarP(&logLevel, "log-level", "l", "info", "log level (debug, info, warn, error, fatal)")
	return cmd
}
