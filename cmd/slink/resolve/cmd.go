package resolve

import (
	"github.com/agentio/slink/cmd/slink/resolve/did"
	"github.com/agentio/slink/cmd/slink/resolve/handle"
	"github.com/agentio/slink/cmd/slink/resolve/now"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "resolve",
		Short: "Resolve atproto identifiers",
	}
	cmd.AddCommand(did.Cmd())
	cmd.AddCommand(handle.Cmd())
	cmd.AddCommand(now.Cmd())
	return cmd
}
