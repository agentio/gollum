package xrpc

import (
	"github.com/agentio/slink/pkg/lexica"
	"github.com/agentio/slink/pkg/slink"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	var input string
	var output string
	var manifest string
	var _loglevel string
	var cmd = &cobra.Command{
		Use:   "xrpc",
		Short: "Generate xrpc handlers and structs for a directory of Lexicon files",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := slink.SetLogLevel(_loglevel); err != nil {
				return err
			}
			catalog := lexica.NewCatalog()
			if err := catalog.Load(input, false /* skip lint */); err != nil {
				return err
			}
			if manifest != "" {
				m, err := lexica.ReadManifest(manifest)
				if err != nil {
					return err
				}
				if err = m.Expand(); err != nil {
					return err
				}
			}
			if err := catalog.GenerateXRPCHandlers(output); err != nil {
				return err
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&input, "input", "i", "lexicons", "input directory")
	cmd.Flags().StringVarP(&output, "output", "o", "gen/xrpc", "output directory")
	cmd.Flags().StringVarP(&manifest, "manifest", "m", "", "manifest")
	cmd.Flags().StringVarP(&_loglevel, "log-level", "l", "warn", "log level (debug, info, warn, error, fatal)")
	return cmd
}
