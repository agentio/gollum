package createinvitecode

import (
	"encoding/json"

	"github.com/agentio/slink/gen/com_atproto"
	xrpc_local "github.com/agentio/slink/pkg/xrpc/local"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	var forAccount string
	var useCount int64
	cmd := &cobra.Command{
		Use:   "create-invite-code",
		Short: "Create invite code",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			client := xrpc_local.NewClient()
			client.Host = "http://localhost:5050"
			response, err := com_atproto.ServerCreateInviteCode(cmd.Context(),
				client,
				&com_atproto.ServerCreateInviteCode_Input{
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
