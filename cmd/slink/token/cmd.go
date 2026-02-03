package token

import (
	"github.com/agentio/slink/cmd/slink/token/generate"
	"github.com/agentio/slink/cmd/slink/token/verify"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "token",
		Short: "Operate on tokens",
	}
	cmd.AddCommand(generate.Cmd())
	cmd.AddCommand(verify.Cmd())
	return cmd
}
