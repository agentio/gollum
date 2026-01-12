package lexica

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
)

type Lexica struct {
}

func (lexica *Lexica) LoadTree(root string) error {
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
	log.Infof("%s (%d)", path, len(b))
	var lexicon Lexicon
	err = json.Unmarshal(b, &lexicon)
	if err != nil {
		return err
	}
	log.Debugf("%+v", lexicon)
	lexicon.Validate(path)
	return nil
}

// https://atproto.com/specs/lexicon#lexicon-files
type Lexicon struct {
	Lexicon     int            `json:"lexicon"`
	Id          string         `json:"id"`
	Description string         `json:"description"`
	Defs        map[string]Def `json:"defs"`
}

func (lexicon *Lexicon) Validate(path string) error {
	if lexicon.Lexicon != 1 {
		log.Warnf("%s unexpected value for lexicon version: %d", path, lexicon.Lexicon)
	}
	expected := strings.ReplaceAll(lexicon.Id, ".", "/") + ".json"
	if !strings.HasSuffix(path, expected) {
		log.Warnf("%s does not match %s", path, expected)
	}
	defCount := len(lexicon.Defs)
	if defCount == 0 {
		log.Warnf("%s has no defs", path)
	}
	log.Infof("%s has %d defs", path, defCount)
	return nil
}

type Def struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}
