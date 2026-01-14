package lexica

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
)

type Lexica struct {
	Lexicons []*Lexicon
}

var _lexica *Lexica

func NewLexica() *Lexica {
	_lexica = &Lexica{}
	return _lexica
}

func (lexica *Lexica) LoadSources(root string) error {
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && filepath.Ext(path) == ".json" {
			if err := lexica.LoadFile(path); err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func (lexica *Lexica) LoadFile(path string) error {
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
	lexica.Lexicons = append(lexica.Lexicons, &lexicon)
	return nil
}

func Lookup(id string) *Lexicon {
	for _, lexicon := range _lexica.Lexicons {
		if lexicon.Id == id {
			return lexicon
		}
	}
	return nil
}
