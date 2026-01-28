package pds

import (
	"errors"
	"fmt"

	"github.com/agentio/slink/pkg/resolve"
	"github.com/agentio/slink/pkg/slink"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	var loglevel string
	cmd := &cobra.Command{
		Use:   "pds HANDLE",
		Short: "Lookup the PDS host for a handle",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := slink.SetLogLevel(loglevel); err != nil {
				return err
			}
			did, err := resolve.Handle(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			d, err := resolve.Did(cmd.Context(), did)
			if err != nil {
				return err
			}
			for _, s := range d.Service {
				if s.ID == "#atproto_pds" {
					fmt.Fprintf(cmd.OutOrStdout(), "%s\n", s.ServiceEndpoint)
					return nil
				}
			}
			return errors.New("handle has no #atproto_pds service")
		},
	}
	cmd.Flags().StringVarP(&loglevel, "log", "l", "warn", "log level (debug, info, warn, error, fatal)")
	return cmd
}
