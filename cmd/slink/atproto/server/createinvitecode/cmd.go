package createinvitecode

import (
	"github.com/agentio/slink/api"
	"github.com/agentio/slink/pkg/common"
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
			client := common.NewClient()
			response, err := api.ServerCreateInviteCode(
				cmd.Context(),
				client,
				&api.ServerCreateInviteCode_Input{
					ForAccount: common.StringPointerOrNil(forAccount),
					UseCount:   useCount,
				},
			)
			if err != nil {
				return err
			}
			return common.Write(cmd.OutOrStdout(), response)
		},
	}
	cmd.Flags().StringVar(&forAccount, "for-account", "", "for account")
	cmd.Flags().Int64Var(&useCount, "use-count", 1, "use count")
	return cmd
}
