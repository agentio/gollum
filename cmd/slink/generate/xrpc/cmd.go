package xrpc

import (
	"github.com/agentio/slink/pkg/common"
	"github.com/agentio/slink/pkg/lexica"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	var input string
	var output string
	var _loglevel string
	var cmd = &cobra.Command{
		Use:   "xrpc",
		Short: "Generate xrpc handlers and structs for a directory of Lexicon files",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := common.SetLogLevel(_loglevel); err != nil {
				return err
			}
			catalog := lexica.NewCatalog()
			if err := catalog.Load(input, false /* skip lint */); err != nil {
				return err
			}
			if err := catalog.GenerateXRPCHandlers(output); err != nil {
				return err
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&input, "input", "i", "lexicons", "input directory")
	cmd.Flags().StringVarP(&output, "output", "o", "gen/xrpc", "output directory")
	cmd.Flags().StringVarP(&_loglevel, "log-level", "l", "warn", "log level (debug, info, warn, error, fatal)")
	return cmd
}
