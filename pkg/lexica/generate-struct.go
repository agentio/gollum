package lexica

import (
	"fmt"
	"slices"
	"sort"
	"strings"
)

func (lexicon *Lexicon) generateStruct(defname, description string, properties map[string]Property, required []string) string {
	var s string
	s += "// " + description + "\n"
	s += lexicon.renderStruct(defname, properties, required)
	s += lexicon.renderDependencies(defname, properties, required)
	return s
}

func (lexicon *Lexicon) renderStruct(defname string, properties map[string]Property, required []string) string {
	var s string
	s += "type " + defname + " struct {\n"
	var propnames []string
	for propname := range properties {
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
			itemstype := lexicon.resolveItemsType(defname, propname, property.Items)
			if required {
				s += capitalize(propname) + " []" + itemstype + " `json:" + `"` + propname + `,omitempty"` + "`\n"
			} else {
				s += capitalize(propname) + " []" + itemstype + " `json:" + `"` + propname + `,omitempty"` + "`\n"
			}
		case "ref":
			reftype := lexicon.resolveRefType(property.Ref)
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
			uniontype := lexicon.resolveUnionType(defname, propname)
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
	s += "}\n\n"
	return s
}

func (lexicon *Lexicon) renderDependencies(defname string, properties map[string]Property, required []string) string {
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
			uniontype := lexicon.resolveUnionType(defname, propname)
			s += "type " + uniontype + " struct {\n"
			for _, ref := range property.Refs {
				fieldname := lexicon.unionfieldname(ref)
				fieldtype := lexicon.unionfieldtype(ref)
				s += fieldname + " " + fieldtype + "\n"
			}
			s += "}\n\n"
		case "array":
			if property.Items.Type == "union" {
				uniontype := lexicon.resolveUnionType(defname, propname) + "_Elem"
				s += "type " + uniontype + " struct {\n"
				for _, ref := range property.Items.Refs {
					fieldname := lexicon.unionfieldname(ref)
					fieldtype := lexicon.unionfieldtype(ref)
					s += fieldname + " " + fieldtype + "\n"
				}
				s += "}\n\n"
			}
		}
	}
	return s
}

func (lexicon *Lexicon) unionfieldname(ref string) string {
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

func (lexicon *Lexicon) unionfieldtype(ref string) string {
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
