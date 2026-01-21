package lexica

import (
	"os"
	"sort"
	"strings"

	"github.com/charmbracelet/log"
	"golang.org/x/tools/imports"
)

func capitalize(s string) string {
	return strings.ToUpper(s[0:1]) + s[1:]
}

func sortedPropertyNames(properties map[string]Property) []string {
	var propnames []string
	for propname := range properties {
		propnames = append(propnames, propname)
	}
	sort.Strings(propnames)
	return propnames
}

func symbolForID(id string) string {
	id = strings.TrimPrefix(id, "com.atproto.") // put these symbols in the top-level namespace
	var s strings.Builder
	for _, part := range strings.Split(id, ".") {
		s.WriteString(capitalize(part))
	}
	return s.String()
}
func writeFormattedFile(filename string, body string) error {
	formatted, err := imports.Process(filename, []byte(body), nil)
	if err != nil {
		log.Errorf("failed to run goimports: %v\n%s", err, body)
		os.WriteFile(filename, []byte(body), 0644)
		return nil
	}
	return os.WriteFile(filename, []byte(formatted), 0644)
}
