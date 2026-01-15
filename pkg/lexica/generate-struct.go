package lexica

import (
	"fmt"
	"slices"
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
	propertyNames := sortedPropertyNames(properties)
	for _, propertyName := range propertyNames {
		property := properties[propertyName]
		required := slices.Contains(required, propertyName)
		switch property.Type {
		case "boolean":
			if required {
				s += capitalize(propertyName) + " bool `json:" + `"` + propertyName + `"` + "`\n"
			} else {
				s += capitalize(propertyName) + " *bool `json:" + `"` + propertyName + `,omitempty"` + "`\n"
			}
		case "integer":
			if required {
				s += capitalize(propertyName) + " int64 `json:" + `"` + propertyName + `,omitempty"` + "`\n"
			} else {
				s += capitalize(propertyName) + " *int64 `json:" + `"` + propertyName + `,omitempty"` + "`\n"
			}
		case "string":
			if required {
				s += capitalize(propertyName) + " string `json:" + `"` + propertyName + `,omitempty"` + "`\n"
			} else {
				s += capitalize(propertyName) + " *string `json:" + `"` + propertyName + `,omitempty"` + "`\n"
			}
		case "array":
			itemstype := lexicon.resolveItemsType(defname, propertyName, property.Items)
			if required {
				s += capitalize(propertyName) + " []" + itemstype + " `json:" + `"` + propertyName + `,omitempty"` + "`\n"
			} else {
				s += capitalize(propertyName) + " []" + itemstype + " `json:" + `"` + propertyName + `,omitempty"` + "`\n"
			}
		case "ref":
			reftype := lexicon.resolveRefType(property.Ref)
			s += capitalize(propertyName) + reftype + " `json:" + `"` + propertyName + `,omitempty"` + "`\n"
		case "unknown":
			if required {
				s += capitalize(propertyName) + " interface{} `json:" + `"` + propertyName + `,omitempty"` + "`\n"
			} else {
				s += capitalize(propertyName) + " *interface{} `json:" + `"` + propertyName + `,omitempty"` + "`\n"
			}
		case "blob":
			if required {
				s += capitalize(propertyName) + " []byte `json:" + `"` + propertyName + `"` + "`\n"
			} else {
				s += capitalize(propertyName) + " *[]byte `json:" + `"` + propertyName + `,omitempty"` + "`\n"
			}
		case "union":
			uniontype := lexicon.resolveUnionType(defname, propertyName)
			s += capitalize(propertyName) + " " + uniontype + " `json:" + `"` + propertyName + `,omitempty"` + "`\n"
		case "bytes":
			if required {
				s += capitalize(propertyName) + " []byte `json:" + `"` + propertyName + `"` + "`\n"
			} else {
				s += capitalize(propertyName) + " *[]byte `json:" + `"` + propertyName + `,omitempty"` + "`\n"
			}
		case "cid-link":
			if required {
				s += capitalize(propertyName) + " string `json:" + `"` + propertyName + `,omitempty"` + "`\n"
			} else {
				s += capitalize(propertyName) + " *string `json:" + `"` + propertyName + `,omitempty"` + "`\n"
			}
		default:
			s += "// FIXME: unsupported property type " + propertyName + " " + property.Type + " " + fmt.Sprintf("required=%t %+v", required, property) + "\n"
		}
	}
	s += "}\n\n"
	return s
}

func (lexicon *Lexicon) renderDependencies(defname string, properties map[string]Property, required []string) string {
	var s string
	propertyNames := sortedPropertyNames(properties)
	for _, propertyName := range propertyNames {
		property := properties[propertyName]
		switch property.Type {
		case "union":
			uniontype := lexicon.resolveUnionType(defname, propertyName)
			s += "type " + uniontype + " struct {\n"
			for _, ref := range property.Refs {
				fieldname := lexicon.unionFieldName(ref)
				fieldtype := lexicon.unionFieldType(ref)
				s += fieldname + " " + fieldtype + "\n"
			}
			s += "}\n\n"
		case "array":
			if property.Items.Type == "union" {
				uniontype := lexicon.resolveUnionType(defname, propertyName) + "_Elem"
				s += "type " + uniontype + " struct {\n"
				for _, ref := range property.Items.Refs {
					fieldname := lexicon.unionFieldName(ref)
					fieldtype := lexicon.unionFieldType(ref)
					s += fieldname + " " + fieldtype + "\n"
				}
				s += "}\n\n"
			}
		}
	}
	return s
}

func (lexicon *Lexicon) unionFieldName(ref string) string {
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

func (lexicon *Lexicon) unionFieldType(ref string) string {
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
