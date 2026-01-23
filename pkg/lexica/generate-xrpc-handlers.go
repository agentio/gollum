package lexica

import (
	"os"
	"sync"
)

func (catalog *Catalog) GenerateXRPCHandlers(root string) error {
	os.RemoveAll(root)
	var wg sync.WaitGroup
	for _, lexicon := range catalog.Lexicons {
		wg.Go(func() {
			lexicon.generateXRPCSourceFile(root)
		})
	}
	wg.Wait()
	return nil
}
