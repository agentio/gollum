package lexica

import (
	"strings"

	"github.com/charmbracelet/log"
)

// https://atproto.com/specs/lexicon#lexicon-files
type Lexicon struct {
	Lexicon     int             `json:"lexicon"`
	Id          string          `json:"id"`
	Description string          `json:"description"`
	Defs        map[string]*Def `json:"defs"`
}

func (lexicon *Lexicon) Lookup(id string) *Def {
	d, ok := lexicon.Defs[id]
	if !ok {
		return nil
	}
	return d
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
	for k, v := range lexicon.Defs {
		v.Validate(path + ":" + k)
	}
	return nil
}
