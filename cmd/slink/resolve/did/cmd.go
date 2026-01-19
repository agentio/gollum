package did

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "did DID",
		Short: "Fetch the DID document for a DID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			did := args[0]
			var url string
			if strings.HasPrefix(did, "did:plc:") {
				url = fmt.Sprintf("https://plc.directory/%s", did)
			} else if strings.HasPrefix(did, "did:web:") {
				url = fmt.Sprintf("https://%s/.well-known/did.json", strings.TrimPrefix(did, "did:web:"))
			} else {
				return fmt.Errorf("%s is not a valid did", did)
			}
			req, err := http.NewRequestWithContext(cmd.Context(), "GET", url, nil)
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
				return fmt.Errorf("%s is not in the PLC registry", did)
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
