package verify

import (
	"encoding/json"
	"fmt"

	"github.com/agentio/slink/pkg/slink"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "verify TOKEN",
		Short: "Verify an authorization token",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			t, err := slink.VerifyAuthToken(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			b, _ := json.MarshalIndent(t, "", "  ")
			fmt.Fprintf(cmd.OutOrStdout(), "%s\n", string(b))
			return nil
		},
	}
	return cmd
}
