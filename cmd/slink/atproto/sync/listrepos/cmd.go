package listrepos

import (
	"encoding/json"

	"github.com/agentio/slink/gen/com_atproto"
	xrpc_local "github.com/agentio/slink/pkg/xrpc/local"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	var cursor string
	var limit int64
	cmd := &cobra.Command{
		Use:   "list-repos",
		Short: "List repos",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			client := xrpc_local.NewClient()
			response, err := com_atproto.SyncListRepos(cmd.Context(),
				client,
				cursor,
				limit)
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
	return cmd
}
