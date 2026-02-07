package now

import (
	"fmt"

	"github.com/agentio/slink/pkg/slink"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "now",
		Short: "Get the current time in ATProto timestamp format",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintf(cmd.OutOrStdout(), "%s\n", slink.Now())
			return nil
		},
	}
	return cmd
}
