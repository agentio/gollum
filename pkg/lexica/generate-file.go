package lexica

import (
	"encoding/json"
	"os"
	"strings"
	"sync"

	"github.com/charmbracelet/log"
	"golang.org/x/tools/imports"
)

func (lexica *Lexica) GenerateCode(root string) error {
	os.RemoveAll(root)
	var wg sync.WaitGroup
	for _, lexicon := range lexica.Lexicons {
		wg.Go(func() {
			lexicon.generateLexiconSourceFile(root)
		})
	}
	wg.Wait()
	return nil
}

func (lexicon *Lexicon) generateLexiconSourceFile(root string) {
	filename := lexiconFileName(root, lexicon.Id)
	if filename == "" {
		return
	}
	packagename := lexiconPackageName(lexicon.Id)
	if packagename == "" {
		return
	}
	lexicon.generateFile(filename, packagename)
}

func lexiconFileName(root, id string) string {
	parts := strings.Split(id, ".")
	if len(parts) != 4 {
		log.Warnf("skipping three-segment name %s", id)
		return ""
	}
	d := root + "/" + parts[0] + "_" + parts[1]
	os.MkdirAll(d, 0755)
	filename := d + "/" + parts[2] + "_" + parts[3] + ".go"
	return filename
}

func lexiconPackageName(id string) string {
	parts := strings.Split(id, ".")
	if len(parts) != 4 {
		log.Warnf("skipping three-segment name %s", id)
		return ""
	}
	packagename := parts[0] + "_" + parts[1]
	return packagename
}

func (lexicon *Lexicon) generateFile(filename, packagename string) error {
	s := "// " + lexicon.Id + "\n\n"
	s += "package " + packagename + "\n\n"

	s += `import "github.com/agentio/slink/pkg/xrpc"` + "\n"
	s += `import "github.com/agentio/slink/gen/com_atproto"` + "\n"
	s += `import "github.com/agentio/slink/gen/app_bsky"` + "\n"

	prefix := codeprefix(lexicon.Id)
	for name, def := range lexicon.Defs {
		s += lexicon.generateDef(def, name, prefix)
	}

	if false { // append lexicon source to generated file
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
	return capitalize(parts[2]) + capitalize(parts[3])
}
