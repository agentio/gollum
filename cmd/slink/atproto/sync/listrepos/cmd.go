package listrepos

import (
	"github.com/agentio/slink/api"
	"github.com/agentio/slink/pkg/common"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	var cursor string
	var limit int64
	cmd := &cobra.Command{
		Use:   "list-repos",
		Short: api.SyncListRepos_Description,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			client := common.NewClient()
			response, err := api.SyncListRepos(
				cmd.Context(),
				client,
				cursor,
				limit)
			if err != nil {
				return err
			}
			return common.Write(cmd.OutOrStdout(), response)

		},
	}
	cmd.Flags().StringVar(&cursor, "cursor", "", "")
	cmd.Flags().Int64Var(&limit, "limit", 100, "")
	return cmd
}
