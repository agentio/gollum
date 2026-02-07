package resolve

import (
	"github.com/agentio/slink/cmd/internal/resolve/did"
	"github.com/agentio/slink/cmd/internal/resolve/doc"
	"github.com/agentio/slink/cmd/internal/resolve/now"
	"github.com/agentio/slink/cmd/internal/resolve/pds"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "resolve",
		Short: "Resolve atproto identifiers",
	}
	cmd.AddCommand(did.Cmd())
	cmd.AddCommand(doc.Cmd())
	cmd.AddCommand(now.Cmd())
	cmd.AddCommand(pds.Cmd())
	return cmd
}
