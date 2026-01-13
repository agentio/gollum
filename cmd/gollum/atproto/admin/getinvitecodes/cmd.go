package getinvitecodes

import (
	"encoding/json"

	"github.com/agentio/gollum/gen/com_atproto"
	xrpc_local "github.com/agentio/gollum/pkg/xrpc/local"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	var cursor string
	var limit int64
	var sort string
	cmd := &cobra.Command{
		Use:   "get-invite-codes",
		Short: "Get invite codes",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			client := xrpc_local.NewClient()
			client.Host = "http://localhost:5050"
			response, err := com_atproto.AdminGetInviteCodes(cmd.Context(),
				client,
				cursor,
				limit,
				sort,
			)
			if err != nil {
				return err
			}
			b, err := json.MarshalIndent(response, "", "  ")
			cmd.OutOrStdout().Write(b)
			cmd.OutOrStdout().Write([]byte("\n"))
			return nil
		},
	}
	cmd.Flags().StringVar(&cursor, "cursor", "", "")
	cmd.Flags().Int64Var(&limit, "limit", 100, "")
	cmd.Flags().StringVar(&sort, "sort", "", "")
	return cmd
}
