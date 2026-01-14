package lexica

import (
	"fmt"
	"slices"
	"sort"
	"strings"
)

func parseParameters(parameters *Parameters) (string, bool) {
	var parms []string
	var parameterNames []string
	for parameterName := range parameters.Properties {
		parameterNames = append(parameterNames, parameterName)
	}
	sort.Strings(parameterNames)
	for _, parameterName := range parameterNames {
		parameterValue := parameters.Properties[parameterName]
		declaration := parameterName + " "
		switch parameterValue.Type {
		case "integer":
			declaration += "int64"
		case "string":
			declaration += "string"
		case "boolean":
			declaration += "bool"
		case "array":
			if parameterValue.Items.Type == "string" {
				declaration += "[]string"
			} else {
				return "/* FIXME */", false
			}
		default:
			return "/* FIXME */", false
		}
		parms = append(parms, declaration)
	}
	return ", " + strings.Join(parms, ", "), true
}

func codeprefix(id string) string {
	parts := strings.Split(id, ".")
	if len(parts) != 4 {
		return ""
	}
	return capitalize(parts[2]) + capitalize(parts[3])
}

func capitalize(s string) string {
	return strings.ToUpper(s[0:1]) + s[1:]
}

func renderProperties(lexicon *Lexicon, defname string, properties map[string]Property, required []string) string {
	var s string

	var propnames []string
	for propname, _ := range properties {
		propnames = append(propnames, propname)
	}
	sort.Strings(propnames)

	for _, propname := range propnames {
		property := properties[propname]
		required := slices.Contains(required, propname)
		switch property.Type {
		case "boolean":
			if required {
				s += capitalize(propname) + " bool `json:" + `"` + propname + `"` + "`\n"
			} else {
				s += capitalize(propname) + " *bool `json:" + `"` + propname + `,omitempty"` + "`\n"
			}
		case "integer":
			if required {
				s += capitalize(propname) + " int64 `json:" + `"` + propname + `,omitempty"` + "`\n"
			} else {
				s += capitalize(propname) + " *int64 `json:" + `"` + propname + `,omitempty"` + "`\n"
			}
		case "string":
			if required {
				s += capitalize(propname) + " string `json:" + `"` + propname + `,omitempty"` + "`\n"
			} else {
				s += capitalize(propname) + " *string `json:" + `"` + propname + `,omitempty"` + "`\n"
			}
		case "array":
			itemstype := resolveItemsType(lexicon, defname, propname, property.Items)
			if required {
				s += capitalize(propname) + " []" + itemstype + " `json:" + `"` + propname + `,omitempty"` + "`\n"
			} else {
				s += capitalize(propname) + " []" + itemstype + " `json:" + `"` + propname + `,omitempty"` + "`\n"
			}
		case "ref":
			reftype := resolveRefType(lexicon, property.Ref)
			s += capitalize(propname) + reftype + " `json:" + `"` + propname + `,omitempty"` + "`\n"
		case "unknown":
			if required {
				s += capitalize(propname) + " interface{} `json:" + `"` + propname + `,omitempty"` + "`\n"
			} else {
				s += capitalize(propname) + " *interface{} `json:" + `"` + propname + `,omitempty"` + "`\n"
			}
		case "blob":
			if required {
				s += capitalize(propname) + " []byte `json:" + `"` + propname + `"` + "`\n"
			} else {
				s += capitalize(propname) + " *[]byte `json:" + `"` + propname + `,omitempty"` + "`\n"
			}
		case "union":
			uniontype := resolveUnionType(lexicon, defname, propname)
			s += capitalize(propname) + " " + uniontype + " `json:" + `"` + propname + `,omitempty"` + "`\n"
		case "bytes":
			if required {
				s += capitalize(propname) + " []byte `json:" + `"` + propname + `"` + "`\n"
			} else {
				s += capitalize(propname) + " *[]byte `json:" + `"` + propname + `,omitempty"` + "`\n"
			}
		case "cid-link":
			if required {
				s += capitalize(propname) + " string `json:" + `"` + propname + `,omitempty"` + "`\n"
			} else {
				s += capitalize(propname) + " *string `json:" + `"` + propname + `,omitempty"` + "`\n"
			}
		default:
			s += "// FIXME: unsupported property type " + propname + " " + property.Type + " " + fmt.Sprintf("required=%t %+v", required, property) + "\n"
		}
	}
	return s
}

func renderDependentTypes(lexicon *Lexicon, defname string, properties map[string]Property, required []string) string {
	var s string

	var propnames []string
	for propname := range properties {
		propnames = append(propnames, propname)
	}
	sort.Strings(propnames)
	for _, propname := range propnames {
		property := properties[propname]
		switch property.Type {
		case "union":
			uniontype := resolveUnionType(lexicon, defname, propname)
			s += "type " + uniontype + " struct {\n"
			for _, ref := range property.Refs {
				fieldname := unionfieldname(lexicon, ref)
				fieldtype := unionfieldtype(lexicon, ref)
				s += fieldname + " " + fieldtype + "\n"
			}
			s += "}\n\n"
		case "array":
			if property.Items.Type == "union" {
				uniontype := resolveUnionType(lexicon, defname, propname) + "_Elem"
				s += "type " + uniontype + " struct {\n"
				for _, ref := range property.Items.Refs {
					fieldname := unionfieldname(lexicon, ref)
					fieldtype := unionfieldtype(lexicon, ref)
					s += fieldname + " " + fieldtype + "\n"
				}
				s += "}\n\n"
			}
		}
	}
	return s
}

func unionfieldname(lexicon *Lexicon, ref string) string {
	parts := strings.Split(ref, "#")
	if len(parts) == 2 || len(parts) == 1 {
		var id string
		var tag string
		if len(parts) == 2 {
			id = parts[0]
			tag = parts[1]
		} else {
			id = parts[0]
			tag = "main"
		}
		if id == "" {
			id = lexicon.Id
		}

		var reftype string
		reflexicon := Lookup(id)
		if reflexicon != nil {
			refdef := reflexicon.Lookup(tag)
			if refdef != nil {
				reftype = refdef.Type
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

		if reftype == "array" {
			return "[]" + name + "_Elem"
		}

		return name
	} else {
		return "/* FIXME union field ref " + fmt.Sprintf("%+v", ref) + " */ string"
	}
}

func unionfieldtype(lexicon *Lexicon, ref string) string {
	parts := strings.Split(ref, "#")
	if len(parts) == 2 || len(parts) == 1 {
		var id string
		var tag string
		if len(parts) == 2 {
			id = parts[0]
			tag = parts[1]
		} else {
			id = parts[0]
			tag = "main"
		}
		if id == "" {
			id = lexicon.Id
		}
		var reftype string
		reflexicon := Lookup(id)
		if reflexicon != nil {
			refdef := reflexicon.Lookup(tag)
			if refdef != nil {
				reftype = refdef.Type
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

		if reftype == "array" {
			return "[]" + name + "_Elem"
		}

		return "*" + name
	} else {
		return "/* FIXME union field ref " + fmt.Sprintf("%+v", ref) + " */ string"
	}
}

func resolveUnionType(lexicon *Lexicon, defname, propname string) string {
	return capitalize(defname) + "_" + capitalize(propname) // "string"
}

func resolveItemsType(lexicon *Lexicon, defname, propname string, items *Items) string {
	switch items.Type {
	case "string":
		return "string"
	case "unknown":
		return "interface{}"
	case "cid-link":
		return "string"
	case "ref":
		ref := items.Ref
		if ref[0] == '#' {
			parts := strings.Split(lexicon.Id, ".")
			if len(parts) != 4 {
				return "/* FIXME: i can't parse this " + lexicon.Id + " */ string"
			}
			typename := capitalize(parts[2]) + capitalize(parts[3]) + "_" + capitalize(ref[1:])
			return "*" + typename
		} else {
			parts := strings.Split(ref, "#")
			if len(parts) != 2 && len(parts) != 1 {
				return "/* FIXME not enough parts" + fmt.Sprintf("%+v", ref) + " */ string"
			}
			id := parts[0]
			var tag string
			if len(parts) == 2 {
				tag = parts[1]
			} else {
				tag = "main"
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

			return "*" + name
		}
	case "union":
		return "*" + capitalize(defname) + "_" + capitalize(propname+"_Elem")
	default:
	}
	return "/* FIXME unsupported items type: " + items.Type + " */ string"
}

func resolveRefType(lexicon *Lexicon, ref string) string {
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
			var id string
			var tag string
			if len(parts) == 2 {
				id = parts[0]
				tag = parts[1]
			} else {
				id = parts[0]
				tag = "main"
			}

			var reftype string
			reflexicon := Lookup(id)
			if reflexicon != nil {
				refdef := reflexicon.Lookup(tag)
				if refdef != nil {
					reftype = refdef.Type
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

			if reftype == "array" {
				return "[]" + name + "_Elem"
			}

			return "*" + name
		} else {
			return "/* FIXME ref " + fmt.Sprintf("%+v", ref) + " */ string"
		}
	}
}
