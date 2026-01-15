package createinvitecode

import (
	"encoding/json"

	"github.com/agentio/slink/api"
	xrpc_sidecar "github.com/agentio/slink/pkg/xrpc/sidecar"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	var forAccount string
	var useCount int64
	cmd := &cobra.Command{
		Use:   "create-invite-code",
		Short: api.ServerCreateInviteCode_Description,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			client := xrpc_sidecar.NewClient()
			response, err := api.ServerCreateInviteCode(cmd.Context(),
				client,
				&api.ServerCreateInviteCode_Input{
					ForAccount: stringPointerOrNil(forAccount),
					UseCount:   useCount,
				})
			if err != nil {
				return err
			}
			b, err := json.MarshalIndent(response, "", "  ")
			cmd.OutOrStdout().Write(b)
			cmd.OutOrStdout().Write([]byte("\n"))
			return nil
		},
	}
	cmd.Flags().StringVar(&forAccount, "for-account", "", "for account")
	cmd.Flags().Int64Var(&useCount, "use-count", 1, "use count")
	return cmd
}

func stringPointerOrNil(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
