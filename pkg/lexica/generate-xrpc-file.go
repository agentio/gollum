package lexica

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
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
	if err := lexicon.generateFile(filename, packagename); err != nil {
		log.Errorf("error writing file %s %s", filename, err)
	}
}

func lexiconFileName(root, id string) string {
	base := strings.ReplaceAll(id, ".", "-")
	os.MkdirAll(root, 0755)
	filename := root + "/" + base + ".go"
	return filename
}

func lexiconPackageName(root, id string) string {
	return filepath.Base(root)
}

func (lexicon *Lexicon) generateFile(filename, packagename string) error {
	log.Debugf("generating %s", filename)
	s := &strings.Builder{}
	fmt.Fprintf(s, "package %s // %s\n\n", packagename, lexicon.Id)
	fmt.Fprintf(s, "import \"github.com/agentio/slink/pkg/common\"\n\n")
	prefix := symbolForID(lexicon.Id)
	for name, def := range lexicon.Defs {
		lexicon.generateDef(s, def, name, prefix)
	}
	if true { // append lexicon source to generated file
		lexicon.generateSourceComment(s)
	}
	return writeFormattedFile(filename, s.String())
}
