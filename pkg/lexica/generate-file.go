package lexica

import (
	"encoding/json"
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
	s := "// " + lexicon.Id + "\n\n"
	s += "package " + packagename + "\n\n"

	s += `import "github.com/agentio/slink/pkg/xrpc"` + "\n"

	prefix := codeprefix(lexicon.Id)
	for name, def := range lexicon.Defs {
		s += lexicon.generateDef(def, name, prefix)
	}

	if true { // append lexicon source to generated file
		filter := func(s string) string {
			return strings.ReplaceAll(s, "*/*", "[ANY]")
		}
		b, _ := json.MarshalIndent(lexicon, "", "  ")
		s += "/*\n"
		s += filter(string(b)) + "\n"
		s += "*/\n"
	}

	formatted, err := imports.Process(filename, []byte(s), nil)
	if err != nil {
		log.Fatalf("failed to run goimports: %v\n%s", err, s)
	}
	return os.WriteFile(filename, []byte(formatted), 0644)
}

func codeprefix(id string) string {
	parts := strings.Split(id, ".")
	if len(parts) != 4 {
		return ""
	}
	return capitalize(parts[0]) + capitalize(parts[1]) + capitalize(parts[2]) + capitalize(parts[3])
}
