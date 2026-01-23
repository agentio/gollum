package lexica

import (
	"fmt"
	"os"
	"strings"

	"github.com/iancoleman/strcase"
)

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

func (catalog *Catalog) generateInternalCommand(path, prompt string) error {
	filename := path + "/cmd.go"
	_, err := os.Stat(filename)
	if err == nil {
		return nil // a command already exists
	}
	parts := strings.Split(path, "/")
	lastpart := parts[len(parts)-1]
	packagename := strings.ReplaceAll(strings.ToLower(lastpart), "-", "_")
	short := prompt
	if len(parts) > 2 {
		short += " under " + strings.ReplaceAll(lastpart, "-", ".")
	}
	subdirectories := getsubdirs(path)

	s := &strings.Builder{}
	packageComment(s, packagename)
	fmt.Fprintf(s, "package %s\n", packagename)
	fmt.Fprintf(s, "import (\n")
	fmt.Fprintf(s, "\"github.com/spf13/cobra\"\n")
	for _, subdir := range subdirectories {
		packagename := strings.ReplaceAll(subdir, "-", "_")
		fmt.Fprintf(s, "%s \"github.com/agentio/slink/%s/%s\"\n", packagename, path, subdir)
	}
	fmt.Fprintf(s, ")\n")
	fmt.Fprintf(s, "func Cmd() *cobra.Command {\n")
	fmt.Fprintf(s, "cmd := &cobra.Command{\n")
	fmt.Fprintf(s, "Use: \"%s\",\n", strings.ReplaceAll(strcase.ToKebab(lastpart), "-", "."))
	fmt.Fprintf(s, "Short: \"%s\",\n", short)
	fmt.Fprintf(s, "}\n")
	for _, subdir := range subdirectories {
		fmt.Fprintf(s, "cmd.AddCommand(%s.Cmd())\n", strings.ReplaceAll(strings.ToLower(subdir), "-", "_"))
	}
	fmt.Fprintf(s, "return cmd\n")
	fmt.Fprintf(s, "}\n")
	return writeFormattedFile(filename, s.String())
}
