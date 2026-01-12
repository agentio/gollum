package lexica

import (
	"os"
	"strings"

	"github.com/charmbracelet/log"
)

func (lexica *Lexica) Generate(root string) error {

	os.RemoveAll("api")

	for _, lexicon := range lexica.Lexicons {
		log.Infof("%s %s", lexicon.Id, filename(lexicon.Id))

	}

	return nil
}

func filename(id string) string {
	parts := strings.Split(id, ".")
	if len(parts) != 4 {
		log.Errorf("wtf %s", id)
		return ""
	}

	d := "api" + "/" + parts[0] + "_" + parts[1]

	os.MkdirAll(d, 0755)

	f := d + "/" + parts[2] + "" + parts[3] + ".go"

	b := []byte("package " + parts[0] + "_" + parts[1] + "\n")

	os.WriteFile(f, b, 0644)

	return f
}
