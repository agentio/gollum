package lexica

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
)

type Catalog struct {
	Lexicons []*Lexicon
}

var _catalog *Catalog

func NewCatalog() *Catalog {
	_catalog = &Catalog{}
	return _catalog
}

func (catalog *Catalog) Load(root string) error {
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && filepath.Ext(path) == ".json" {
			if err := catalog.LoadLexicon(path); err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func (catalog *Catalog) LoadLexicon(path string) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var lexicon Lexicon
	err = json.Unmarshal(b, &lexicon)
	if err != nil {
		return err
	}
	lexicon.Validate(path)
	catalog.Lexicons = append(catalog.Lexicons, &lexicon)
	return nil
}

func LookupLexicon(id string) *Lexicon {
	for _, lexicon := range _catalog.Lexicons {
		if lexicon.Id == id {
			return lexicon
		}
	}
	return nil
}
