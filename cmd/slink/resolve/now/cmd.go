package now

import (
	"fmt"

	"github.com/agentio/slink/pkg/common"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "now",
		Short: "Print the current time in ATProto timestamp format",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintf(cmd.OutOrStdout(), "%s\n", common.Now())
			return nil
		},
	}
	return cmd
}
