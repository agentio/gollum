package updateaccountpassword

import (
	"github.com/agentio/slink/api"
	xrpc_sidecar "github.com/agentio/slink/pkg/xrpc/sidecar"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	var did string
	var password string
	cmd := &cobra.Command{
		Use:   "update-account-password",
		Short: api.AdminUpdateAccountPassword_Description,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			client := xrpc_sidecar.NewClient()
			err := api.AdminUpdateAccountPassword(cmd.Context(),
				client,
				&api.AdminUpdateAccountPassword_Input{
					Did:      did,
					Password: password,
				},
			)
			if err != nil {
				return err
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&did, "did", "", "")
	cmd.Flags().StringVar(&password, "password", "", "")
	return cmd
}
