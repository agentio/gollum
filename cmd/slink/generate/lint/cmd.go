package lint

import (
	"github.com/agentio/slink/pkg/lexica"
	"github.com/agentio/slink/pkg/tool"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	var input string
	var _loglevel string
	var cmd = &cobra.Command{
		Use:   "lint",
		Short: "Check a directory of Lexicon files for possible problems",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := tool.SetLogLevel(_loglevel); err != nil {
				return err
			}
			catalog := lexica.NewCatalog()
			if err := catalog.Load(input, true /* lint */); err != nil {
				return err
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&input, "input", "i", "lexicons", "input directory")
	cmd.Flags().StringVarP(&_loglevel, "log-level", "l", "warn", "log level (debug, info, warn, error, fatal)")
	return cmd
}
