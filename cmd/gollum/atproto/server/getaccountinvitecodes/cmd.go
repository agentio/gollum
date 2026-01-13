package getaccountinvitecodes

import (
	"log"

	"github.com/agentio/gollum/api/com_atproto"
	xrpc_local "github.com/agentio/gollum/pkg/xrpc/local"
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
			client.Host = "http://localhost:5050"
			response, err := com_atproto.ServerGetAccountInviteCodes(cmd.Context(),
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
