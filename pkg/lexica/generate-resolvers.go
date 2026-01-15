package lexica

import (
	"fmt"
	"strings"
)

func (lexicon *Lexicon) resolveUnionType(defname, propname string) string {
	return capitalize(defname) + "_" + capitalize(propname)
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
	return "/* FIXME unsupported items type: " + items.Type + " */ string"
}

func (lexicon *Lexicon) resolveRefType(ref string) string {
	if ref[0] == '#' {
		parts := strings.Split(lexicon.Id, ".")
		if len(parts) != 4 {
			return "/* FIXME: i can't parse this " + lexicon.Id + " */ string"
		}
		typename := capitalize(parts[2]) + capitalize(parts[3]) + "_" + capitalize(ref[1:])
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
			refLexicon := Lookup(id)
			if refLexicon != nil {
				refDef := refLexicon.Lookup(tag)
				if refDef != nil {
					refType = refDef.Type
				}
			}

			idparts := strings.Split(id, ".")
			if len(idparts) != 4 {
				return "/* FIXME " + fmt.Sprintf("%+v", ref) + " */ string"
			}
			name := capitalize(idparts[2]) + capitalize(idparts[3])
			if tag != "main" {
				name += "_" + capitalize(tag)
			}
			// is the ref target in the same package as the lexicon?
			// if not, we need to add the package name prefix

			if !strings.HasPrefix(lexicon.Id, idparts[0]+"."+idparts[1]+".") {
				prefix := idparts[0] + "_" + idparts[1]
				name = prefix + "." + name
			}

			if refType == "array" {
				return "[]" + name + "_Elem"
			}

			return "*" + name
		} else {
			return "/* FIXME ref " + fmt.Sprintf("%+v", ref) + " */ string"
		}
	}
}
