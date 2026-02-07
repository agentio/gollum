package token

import (
	"github.com/agentio/slink/cmd/internal/token/generate"
	"github.com/agentio/slink/cmd/internal/token/verify"
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
