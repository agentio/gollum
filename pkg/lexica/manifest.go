package lexica

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/charmbracelet/log"
)

var _manifest *Manifest

type Manifest struct {
	IDs []string `json:"ids"`
}

func BuildManifest(filename string) (*Manifest, error) {
	m, err := readManifest(filename)
	if err != nil {
		return m, err
	}
	if err = m.expand(); err != nil {
		return m, err
	}
	return m, nil
}

func readManifest(filename string) (*Manifest, error) {
	b, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var manifest Manifest
	err = json.Unmarshal(b, &manifest)
	if err != nil {
		return nil, err
	}
	return &manifest, nil
}

func parse(id string) (nsid string, name string) {
	parts := strings.Split(id, "#")
	nsid = parts[0]
	if len(parts) == 2 {
		name = parts[1]
	} else {
		name = "main"
	}
	return
}

func ManifestIncludes(nsid, name string) bool {
	if _manifest == nil {
		return true
	}
	if name == "main" {
		if slices.Contains(_manifest.IDs, nsid) {
			return true
		}
	}
	return slices.Contains(_manifest.IDs, nsid+"#"+name)
}

func (manifest *Manifest) expand() error {
	for i := 0; i < len(manifest.IDs); i++ {
		id := manifest.IDs[i]
		nsid, name := parse(id)
		l := LookupLexicon(nsid)
		if l == nil {
			return fmt.Errorf("can't find lexicon %s", id)
		}
		def := l.Lookup(name)
		if def == nil {
			return fmt.Errorf("can't find def %s", id)
		}
		err := manifest.AddDependencies(l, def)
		if err != nil {
			return err
		}
	}
	slices.Sort(manifest.IDs)
	_manifest = manifest
	return nil
}

func (manifest *Manifest) addID(l *Lexicon, name string) {
	var id string
	if name[0] == '#' {
		id = l.Id + name
	} else {
		id = name
	}
	if !slices.Contains(manifest.IDs, id) {
		manifest.IDs = append(manifest.IDs, id)
	}
}

func (manifest *Manifest) AddDependencies(l *Lexicon, def *Def) error {
	switch def.Type {
	case "permission-set", "string", "subscription", "token":
		return nil // these types have no dependencies
	case "query", "procedure", "object", "record":
		return manifest.AddDependenciesForDef(l, def)
	case "array":
		return fmt.Errorf("unsupported def type %s", def.Type)
	default:
		return fmt.Errorf("unsupported def type %s", def.Type)
	}
}

func (manifest *Manifest) AddDependenciesForDef(l *Lexicon, def *Def) error {
	if def.Input != nil && def.Input.Encoding == "application/json" {
		manifest.AddDependenciesForProperties(l, def.Input.Schema.Properties)
	}
	if def.Output != nil && def.Output.Encoding == "application/json" {
		manifest.AddDependenciesForProperties(l, def.Output.Schema.Properties)
	}
	if def.Properties != nil {
		manifest.AddDependenciesForProperties(l, def.Properties)
	}
	if def.Record != nil {
		manifest.AddDependenciesForProperties(l, def.Record.Properties)
	}
	return nil
}

func (manifest *Manifest) AddDependenciesForProperties(l *Lexicon, properties map[string]Property) error {
	for propertyName, propertyValue := range properties {
		switch propertyValue.Type {
		case "string", "integer", "boolean", "unknown", "bytes", "blob":
		case "union":
			for _, refName := range propertyValue.Refs {
				manifest.addID(l, refName)
			}
		case "ref":
			manifest.addID(l, propertyValue.Ref)
		case "array":
			switch propertyValue.Items.Type {
			case "ref":
				manifest.addID(l, propertyValue.Items.Ref)
			case "union":
				for _, refName := range propertyValue.Items.Refs {
					manifest.addID(l, refName)
				}
			case "string":
			default:
				log.Warnf("array items %+v", propertyValue.Items)
			}
		case "object":
			manifest.AddDependenciesForProperties(l, propertyValue.Properties)
		default:
			log.Warnf("%s %+v", propertyName, propertyValue)
		}
	}
	return nil
}
