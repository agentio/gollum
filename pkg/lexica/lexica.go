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
	Lexicons []Lexicon
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
	var lexicon Lexicon
	err = json.Unmarshal(b, &lexicon)
	if err != nil {
		return err
	}
	log.Debugf("%s %+v", path, lexicon)
	lexicon.Validate(path)
	lexica.Lexicons = append(lexica.Lexicons, lexicon)
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
	for k, v := range lexicon.Defs {
		v.Validate(path + ":" + k)
	}
	return nil
}

type Def struct {
	Type        string `json:"type"`
	Description string `json:"description"`

	// object
	Required   []string            `json:"required"`
	Properties map[string]Property `json:"properties"`

	// query
	Parameters Parameters `json:"parameters"`
	Output     Output     `json:"output"`

	// procedure
	Input Input `json:"input"`
}

func (def *Def) Validate(path string) error {
	log.Debugf("%s has type %s", path, def.Type)
	switch def.Type {
	case "boolean":
	case "integer":
	case "string":
	case "bytes":
	case "cid-link":
	case "blob":
	case "array":
	case "object":
	case "params":
	case "permission":
	case "token":
	case "ref":
	case "union":
	case "unknown":
	case "record":
	case "query":
	case "procedure":
	case "subscription":
	case "permission-set":
	default:
		log.Warnf("%s has unrecognized type: %s", path, def.Type)
	}
	return nil
}

type Parameters struct {
	Type       string              `json:"type"`
	Properties map[string]Property `json:"properties"`
}

type Output struct {
	Encoding string `json:"encoding"`
	Schema   Schema `json:"schema"`
}

type Input struct {
	Encoding string `json:"encoding"`
	Schema   Schema `json:"schema"`
}

type Schema struct {
	Type       string              `json:"type"`
	Required   []string            `json:"required"`
	Properties map[string]Property `json:"properties"`
}

type Property struct {
	Type  string `json:"type"`
	Ref   string `json:"ref"`
	Items Items  `json:"items"`
}

type Items struct {
	Type string `json:"type"`
	Ref  string `json:"ref"`
}
