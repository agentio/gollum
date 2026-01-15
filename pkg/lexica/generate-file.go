package lexica

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/log"
	"golang.org/x/tools/imports"
)

func (lexicon *Lexicon) generateLexiconSourceFile(root string) {
	filename := lexiconFileName(root, lexicon.Id)
	if filename == "" {
		return
	}
	packagename := lexiconPackageName(root, lexicon.Id)
	if packagename == "" {
		return
	}
	lexicon.generateFile(filename, packagename)
}

func lexiconFileName(root, id string) string {
	base := strings.ReplaceAll(id, ".", "-")
	os.MkdirAll(root, 0755)
	filename := root + "/" + base + ".go"
	return filename
}

func lexiconPackageName(root, id string) string {
	return root
}

func (lexicon *Lexicon) generateFile(filename, packagename string) error {
	s := &strings.Builder{}
	fmt.Fprintf(s, "package %s // %s\n\n", packagename, lexicon.Id)
	s.WriteString(`import "github.com/agentio/slink/pkg/xrpc"` + "\n\n")
	prefix := idPrefix(lexicon.Id)
	for name, def := range lexicon.Defs {
		lexicon.generateDef(s, def, name, prefix)
	}
	if true { // append lexicon source to generated file
		filter := func(s string) string {
			return strings.ReplaceAll(s, "*/*", "[ANY]")
		}
		b, _ := json.MarshalIndent(lexicon, "", "  ")
		s.WriteString("/*\n")
		s.WriteString(filter(string(b)) + "\n")
		s.WriteString("*/\n")
	}
	formatted, err := imports.Process(filename, []byte(s.String()), nil)
	if err != nil {
		log.Fatalf("failed to run goimports: %v\n%s", err, s.String())
	}
	return os.WriteFile(filename, []byte(formatted), 0644)
}
