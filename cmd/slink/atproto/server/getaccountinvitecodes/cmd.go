package getaccountinvitecodes

import (
	"log"

	"github.com/agentio/slink/api"
	xrpc_local "github.com/agentio/slink/pkg/xrpc/local"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	var createAvailable bool
	var includeUsed bool
	cmd := &cobra.Command{
		Use:   "get-account-invite-codes",
		Short: "Get account invite codes",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			client := xrpc_local.NewClient()
			response, err := api.ComAtprotoServerGetAccountInviteCodes(cmd.Context(),
				client,
				createAvailable,
				includeUsed,
			)
			if err != nil {
				return err
			}
			log.Printf("%+v", response)
			return nil
		},
	}
	cmd.Flags().BoolVar(&createAvailable, "create-available", false, "")
	cmd.Flags().BoolVar(&includeUsed, "include-used", false, "")
	return cmd
}
