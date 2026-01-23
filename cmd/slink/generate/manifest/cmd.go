package manifest

import (
	"github.com/agentio/slink/pkg/lexica"
	"github.com/agentio/slink/pkg/slink"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	var input string
	var _loglevel string
	var cmd = &cobra.Command{
		Use:   "manifest MANIFEST",
		Short: "Generate a list of dependencies in a Lexicon",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := slink.SetLogLevel(_loglevel); err != nil {
				return err
			}
			catalog := lexica.NewCatalog()
			if err := catalog.Load(input, false /* lint */); err != nil {
				return err
			}

			manifest, err := lexica.ReadManifest(args[0])
			if err != nil {
				return err
			}
			if err = manifest.Expand(); err != nil {
				return err
			}
			slink.Write(cmd.OutOrStdout(), "-", manifest)
			return nil
		},
	}
	cmd.Flags().StringVarP(&input, "input", "i", "lexicons", "input directory")
	cmd.Flags().StringVarP(&_loglevel, "log-level", "l", "warn", "log level (debug, info, warn, error, fatal)")
	return cmd
}
