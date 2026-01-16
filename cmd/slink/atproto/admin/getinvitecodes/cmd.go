package getinvitecodes

import (
	"github.com/agentio/slink/api"
	"github.com/agentio/slink/pkg/common"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	var cursor string
	var limit int64
	var sort string
	cmd := &cobra.Command{
		Use:   "get-invite-codes",
		Short: api.AdminGetInviteCodes_Description,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			client := common.NewClient()
			response, err := api.AdminGetInviteCodes(
				cmd.Context(),
				client,
				cursor,
				limit,
				sort,
			)
			if err != nil {
				return err
			}
			return common.Write(cmd.OutOrStdout(), response)

		},
	}
	cmd.Flags().StringVar(&cursor, "cursor", "", "")
	cmd.Flags().Int64Var(&limit, "limit", 100, "")
	cmd.Flags().StringVar(&sort, "sort", "", "")
	return cmd
}
