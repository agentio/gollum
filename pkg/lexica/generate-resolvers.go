package lexica

import (
	"fmt"
	"strings"
)

func (lexicon *Lexicon) resolveUnionType(defname, propname string) string {
	return codePrefix(lexicon.Id) + capitalize(defname) + "_" + capitalize(propname)
}

func (lexicon *Lexicon) resolveItemsType(defname, propname string, items *Items) string {
	switch items.Type {
	case "string":
		return "string"
	case "unknown":
		return "interface{}"
	case "cid-link":
		return "string"
	case "ref":
		return lexicon.resolveRefType(items.Ref)
	case "union":
		return "*" + lexicon.resolveUnionType(defname, propname) + "_Elem"
	default:
	}
	return "/* FIXME defaulting on unsupported items type: " + items.Type + " */ string"
}

func (lexicon *Lexicon) resolveRefType(ref string) string {
	if ref[0] == '#' {
		typename := codePrefix(lexicon.Id) + "_" + capitalize(ref[1:])
		return "*" + typename
	} else {
		parts := strings.Split(ref, "#")
		if len(parts) == 2 || len(parts) == 1 {
			id := parts[0]
			var tag string
			if len(parts) == 2 {
				tag = parts[1]
			} else {
				tag = "main"
			}
			var refType string
			refLexicon := LookupLexicon(id)
			if refLexicon != nil {
				refDef := refLexicon.Lookup(tag)
				if refDef != nil {
					refType = refDef.Type
				}
			}
			name := codePrefix(id)
			if tag != "main" {
				name += "_" + capitalize(tag)
			}
			if refType == "array" {
				return "[]" + name + "_Elem"
			}
			return "*" + name
		} else {
			return "/* FIXME defaulting on unparseable ref " + fmt.Sprintf("%+v", ref) + " */ string"
		}
	}
}
