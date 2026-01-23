package lexica

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
)

var _manifest *Manifest

type Manifest struct {
	IDs []string `json:"ids"`
}

func ReadManifest(filename string) (*Manifest, error) {
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

func (manifest *Manifest) Expand() error {
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
		err := manifest.AddDependencies(name, l, def)
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

func (manifest *Manifest) AddDependencies(name string, l *Lexicon, def *Def) error {
	switch def.Type {
	case "permission-set", "string", "subscription", "token":
		return nil // these types have no dependencies
	case "query", "procedure":
		return manifest.AddDependenciesForCallable(name, l, def)
	case "object", "record":
		return manifest.AddDependenciesForStruct(name, l, def)
	case "array":
		return fmt.Errorf("unsupported def type %s", def.Type)
	default:
		return fmt.Errorf("unsupported def type %s", def.Type)
	}
}

func (manifest *Manifest) AddDependenciesForCallable(name string, l *Lexicon, def *Def) error {
	if def.Input != nil && def.Input.Encoding == "application/json" {
		for paramName, paramValue := range def.Input.Schema.Properties {
			switch paramValue.Type {
			case "string", "integer", "boolean", "unknown", "bytes":
			case "ref":
				manifest.addID(l, paramValue.Ref)
			default:
				log.Printf("input %s %+v", paramName, paramValue)
			}
		}
	}
	if def.Output != nil && def.Output.Encoding == "application/json" {
		for paramName, paramValue := range def.Output.Schema.Properties {
			switch paramValue.Type {
			case "string", "integer", "boolean", "unknown", "bytes":
			case "ref":
				manifest.addID(l, paramValue.Ref)
			default:
				log.Printf("output %s %+v", paramName, paramValue)
			}
		}
	}
	return nil
}

func (manifest *Manifest) AddDependenciesForStruct(name string, l *Lexicon, def *Def) error {
	for propertyName, propertyValue := range def.Properties {
		switch propertyValue.Type {
		case "string", "integer", "boolean", "unknown", "bytes":
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
			default:
				log.Printf("array items %+v", propertyValue.Items)
			}
		default:
			log.Printf("%s %+v", propertyName, propertyValue)
		}
	}
	return nil
}
