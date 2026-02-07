package generate

import (
	"github.com/agentio/slink/cmd/internal/generate/call"
	"github.com/agentio/slink/cmd/internal/generate/check"
	"github.com/agentio/slink/cmd/internal/generate/lint"
	"github.com/agentio/slink/cmd/internal/generate/manifest"
	"github.com/agentio/slink/cmd/internal/generate/xrpc"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate slink code from a directory of Lexicon files",
	}
	cmd.AddCommand(call.Cmd())
	cmd.AddCommand(check.Cmd())
	cmd.AddCommand(lint.Cmd())
	cmd.AddCommand(manifest.Cmd())
	cmd.AddCommand(xrpc.Cmd())
	return cmd
}
