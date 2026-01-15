package getinvitecodes

import (
	"encoding/json"

	"github.com/agentio/slink/api"
	xrpc_sidecar "github.com/agentio/slink/pkg/xrpc/sidecar"
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
			client := xrpc_sidecar.NewClient()
			response, err := api.AdminGetInviteCodes(cmd.Context(),
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
