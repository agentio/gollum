package lexica

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/iancoleman/strcase"
)

func (lexicon *Lexicon) generateLeafCommand(root string) {

	hasCode := false
	for _, def := range lexicon.Defs {
		if def.Type == "procedure" || def.Type == "query" {
			hasCode = true
		}
	}
	if !hasCode {
		return
	}

	id := strings.Replace(lexicon.Id, ".", "-", 1)
	parts := strings.Split(id, ".")
	lastpart := parts[len(parts)-1]
	dirname := root + "/" + strings.ReplaceAll(id, ".", "/")
	dirname = strings.ToLower(dirname)
	os.MkdirAll(dirname, 0755)

	filename := dirname + "/cmd.go"

	s := &strings.Builder{}
	fmt.Fprintf(s, "package %s\n", strings.ReplaceAll(strings.ToLower(lastpart), "-", "_"))
	fmt.Fprintf(s, "import \"github.com/spf13/cobra\"\n")
	fmt.Fprintf(s, "func Cmd() *cobra.Command {\n")
	fmt.Fprintf(s, "cmd := &cobra.Command{\n")
	fmt.Fprintf(s, "Use: \"%s\",\n", strcase.ToKebab(lastpart))
	fmt.Fprintf(s, "RunE: func(cmd *cobra.Command, args []string) error {\n")
	fmt.Fprintf(s, "return errors.New(\"unimplemented\")")
	fmt.Fprintf(s, "},\n")
	fmt.Fprintf(s, "}\n")
	fmt.Fprintf(s, "return cmd\n")
	fmt.Fprintf(s, "}\n")

	/*

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
	*/
	if err := writeFormattedFile(filename, s.String()); err != nil {
		log.Errorf("error writing file %s %s", filename, err)
	}
}
