package lexica

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/charmbracelet/log"
	"github.com/iancoleman/strcase"
	"golang.org/x/tools/imports"
)

func (catalog *Catalog) GenerateCLI(root string) error {
	os.RemoveAll(root)
	var wg sync.WaitGroup
	for _, lexicon := range catalog.Lexicons {
		wg.Go(func() {
			lexicon.generateLexiconCLI(root)
		})
	}
	wg.Wait()

	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if !d.Type().IsDir() {
			return nil
		}
		filename := path + "/cmd.go"
		_, err = os.Stat(filename)
		if err == nil {
			return nil
		}
		log.Printf("generating %s", filename)

		parts := strings.Split(path, "/")
		lastpart := parts[len(parts)-1]

		subdirectories := getsubdirs(path)

		s := &strings.Builder{}
		fmt.Fprintf(s, "package %s\n", strings.ReplaceAll(strings.ToLower(lastpart), "-", "_"))
		fmt.Fprintf(s, "import (\n")
		fmt.Fprintf(s, "\"github.com/spf13/cobra\"\n")
		for _, subdir := range subdirectories {
			packagename := strings.ReplaceAll(subdir, "-", "_")
			fmt.Fprintf(s, "%s \"github.com/agentio/slink/%s/%s\"\n", packagename, path, subdir)
		}
		fmt.Fprintf(s, ")\n")
		fmt.Fprintf(s, "func Cmd() *cobra.Command {\n")
		fmt.Fprintf(s, "cmd := &cobra.Command{\n")
		fmt.Fprintf(s, "Use: \"%s\",\n", strcase.ToKebab(lastpart))
		fmt.Fprintf(s, "}\n")
		for _, subdir := range subdirectories {
			fmt.Fprintf(s, "cmd.AddCommand(%s.Cmd())\n", strings.ReplaceAll(strings.ToLower(subdir), "-", "_"))
		}
		fmt.Fprintf(s, "return cmd\n")
		fmt.Fprintf(s, "}\n")
		formatted, err := imports.Process(filename, []byte(s.String()), nil)
		if err != nil {
			log.Errorf("failed to run goimports: %v\n%s", err, s.String())
			os.WriteFile(filename, []byte(s.String()), 0644)
			return nil
		}
		os.WriteFile(filename, []byte(formatted), 0644)
		return nil
	})

	return nil
}

func getsubdirs(path string) []string {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil
	}

	var subdirs []string
	for _, entry := range entries {
		if entry.IsDir() {
			subdirs = append(subdirs, entry.Name())
		}
	}
	return subdirs
}

func (lexicon *Lexicon) generateLexiconCLI(root string) {

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

	formatted, err := imports.Process(filename, []byte(s.String()), nil)
	if err != nil {
		log.Errorf("failed to run goimports: %v\n%s", err, s.String())
		os.WriteFile(filename, []byte(s.String()), 0644)
		return
	}
	os.WriteFile(filename, []byte(formatted), 0644)
}
