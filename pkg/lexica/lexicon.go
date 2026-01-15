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

type Def struct {
	Type        string `json:"type"`
	Description string `json:"description"`

	// object
	Required   []string            `json:"required,omitempty"`
	Properties map[string]Property `json:"properties,omitempty"`

	// query
	Parameters *Parameters `json:"parameters,omitempty"`
	Output     *Output     `json:"output,omitempty"`

	// procedure
	Input *Input `json:"input,omitempty"`

	// array
	Items *Items `json:"items,omitempty"`
}

func (def *Def) Validate(path string) error {
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
	Type  string   `json:"type,omitempty"`
	Ref   string   `json:"ref,omitempty"`
	Refs  []string `json:"refs,omitempty"`
	Items *Items   `json:"items,omitempty"`
}

type Items struct {
	Type string   `json:"type,omitempty"`
	Ref  string   `json:"ref,omitempty"`
	Refs []string `json:"refs,omitempty"`
}
