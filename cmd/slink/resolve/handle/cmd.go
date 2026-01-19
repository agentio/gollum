package handle

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "handle HANDLE",
		Short: "Lookup the DID for a handle",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			handle := args[0]
			// first check DNS
			records, err := net.LookupTXT(fmt.Sprintf("_atproto.%s", handle))
			if err == nil {
				for _, rec := range records {
					if strings.HasPrefix(rec, "did=") {
						did := strings.Split(rec, "did=")[1]
						if did != "" {
							fmt.Fprintf(cmd.OutOrStdout(), "%s\n", did)
							return nil
						}
					}
				}
			}
			// if that didn't work, check the .well-known/atproto-did path
			req, err := http.NewRequestWithContext(
				cmd.Context(),
				"GET",
				fmt.Sprintf("https://%s/.well-known/atproto-did", handle),
				nil,
			)
			if err != nil {
				return err
			}
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				_, _ = io.Copy(io.Discard, resp.Body)
				return fmt.Errorf("unable to resolve %s", handle)
			}
			b, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "%s\n", string(b))
			return nil
		},
	}
	return cmd
}
