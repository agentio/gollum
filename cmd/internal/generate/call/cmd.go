package call

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
		Use:   "call",
		Short: "Generate a command-line interface to call methods in a directory of Lexicon files",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := slink.SetLogLevel(_loglevel); err != nil {
				return err
			}
			catalog := lexica.NewCatalog()
			if err := catalog.Load(input, false /* skip lint */); err != nil {
				return err
			}
			if manifest != "" {
				_, err := lexica.BuildManifest(manifest)
				if err != nil {
					return err
				}
			}
			if err := catalog.GenerateCallCommands(output); err != nil {
				return err
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&input, "input", "i", "lexicons", "input directory")
	cmd.Flags().StringVarP(&output, "output", "o", "gen/call", "output directory")
	cmd.Flags().StringVarP(&manifest, "manifest", "m", "", "manifest")
	cmd.Flags().StringVarP(&_loglevel, "log", "l", "warn", "log level (debug, info, warn, error, fatal)")
	return cmd
}
