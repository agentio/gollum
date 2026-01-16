package lexica

import (
	"os"

	"github.com/charmbracelet/log"
	"golang.org/x/tools/imports"
)

func writeFormattedFile(filename string, body string) error {
	formatted, err := imports.Process(filename, []byte(body), nil)
	if err != nil {
		log.Errorf("failed to run goimports: %v\n%s", err, body)
		os.WriteFile(filename, []byte(body), 0644)
		return nil
	}
	return os.WriteFile(filename, []byte(formatted), 0644)
}
