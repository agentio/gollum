package getaccountinvitecodes

import (
	"log"

	"github.com/agentio/slink/api"
	xrpc_sidecar "github.com/agentio/slink/pkg/xrpc/sidecar"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	var createAvailable bool
	var includeUsed bool
	cmd := &cobra.Command{
		Use:   "get-account-invite-codes",
		Short: api.ServerGetAccountInviteCodes_Description,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			client := xrpc_sidecar.NewClient()
			response, err := api.ServerGetAccountInviteCodes(cmd.Context(),
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
