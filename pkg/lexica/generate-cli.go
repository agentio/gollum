package lexica

import (
	"io/fs"
	"os"
	"path/filepath"
	"sync"
)

func (catalog *Catalog) GenerateCLI(root string) error {
	os.RemoveAll(root)
	os.MkdirAll(root, 0755)
	var wg sync.WaitGroup
	for _, lexicon := range catalog.Lexicons {
		wg.Go(func() {
			lexicon.generateLeafCommands(root)
		})
	}
	wg.Wait()
	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if d.Type().IsDir() {
			wg.Go(func() {
				catalog.generateInternalCommand(path)
			})
		}
		return nil
	})
	wg.Wait()
	return nil
}
