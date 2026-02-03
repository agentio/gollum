package main

import (
	"os"
	"strings"

	"github.com/agentio/slink/cmd/slink/generate"
	"github.com/agentio/slink/cmd/slink/resolve"
	"github.com/agentio/slink/cmd/slink/token"
	"github.com/agentio/slink/gen/call"
	"github.com/agentio/slink/gen/check"
	"github.com/spf13/cobra"
)

func main() {
	if err := cmd().Execute(); err != nil {
		os.Exit(1)
	}
}

func cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "slink",
		Long: strings.Join(
			[]string{
				``,
				`"Perhaps we’ve shaken him off at last, the miserable slinker!"`,
				``,
				`A tool for working with the AT Protocol.`,
				``,
				`Environment Variables:`,
				`  SLINK_HOST sets the target host (e.g. "https://public.api.bsky.app").`,
				`  SLINK_AUTH sets the authorization header (e.g. "Bearer XXXX").`,
				`  SLINK_ATPROTOPROXY sets the atproto-proxy header.`,
				`  SLINK_PROXYSESSION sets the proxy-session header (used by IO).`,
				`  SLINK_USERDID sets the user-did header (used by IO).`,
			}, "\n"),
	}
	cmd.AddCommand(call.Cmd())
	cmd.AddCommand(check.Cmd())
	cmd.AddCommand(generate.Cmd())
	cmd.AddCommand(resolve.Cmd())
	cmd.AddCommand(token.Cmd())
	return cmd
}
