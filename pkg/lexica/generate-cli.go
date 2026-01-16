package lexica

import (
	"io/fs"
	"os"
	"path/filepath"
	"sync"
)

func (catalog *Catalog) GenerateCLI(root string) error {
	os.RemoveAll(root)
	var wg sync.WaitGroup
	for _, lexicon := range catalog.Lexicons {
		wg.Go(func() {
			lexicon.generateLeafCommand(root)
		})
	}
	wg.Wait()
	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if d.Type().IsDir() {
			return catalog.generateInternalCommand(path)
		}
		return nil
	})
	return nil
}
