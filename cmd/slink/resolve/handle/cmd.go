package handle

import (
	"fmt"

	"github.com/agentio/slink/pkg/resolve"
	"github.com/agentio/slink/pkg/tool"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	var _loglevel string
	cmd := &cobra.Command{
		Use:   "handle HANDLE",
		Short: "Lookup the DID for a handle",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := tool.SetLogLevel(_loglevel); err != nil {
				return err
			}
			did, err := resolve.Handle(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "%s\n", did)
			return nil
		},
	}
	cmd.Flags().StringVarP(&_loglevel, "log-level", "l", "warn", "log level (debug, info, warn, error, fatal)")
	return cmd
}
