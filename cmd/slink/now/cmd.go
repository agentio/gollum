package now

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

const ISO8601 = "2006-01-02T15:04:05.000Z"

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "now",
		Short: "Print the current time in ATProto timestamp format",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintf(cmd.OutOrStdout(), "%s", time.Now().UTC().Format(ISO8601))
			return nil
		},
	}
	return cmd
}
