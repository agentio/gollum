package lexica

import (
	"fmt"
	"slices"
	"strings"
)

func (lexicon *Lexicon) generateStruct(s *strings.Builder, defname, description string, properties map[string]Property, required []string, isRecord bool) {
	fmt.Fprintf(s, "// %s\n", description)
	lexicon.renderStruct(s, defname, properties, required, isRecord)
	lexicon.renderDependencies(s, defname, properties, required)
}

func (lexicon *Lexicon) renderStruct(s *strings.Builder, defname string, properties map[string]Property, required []string, isRecord bool) {
	fmt.Fprintf(s, "type %s struct {\n", defname)
	if isRecord {
		fmt.Fprintf(s, "LexiconTypeID string `json:\"$type,omitempty\"`\n")
	}
	for _, propertyName := range sortedPropertyNames(properties) {
		required := slices.Contains(required, propertyName)
		property := properties[propertyName]
		switch property.Type {
		case "boolean":
			if required {
				fmt.Fprintf(s, "%s bool `json:\"%s\"`\n", capitalize(propertyName), propertyName)
			} else {
				fmt.Fprintf(s, "%s *bool `json:\"%s,omitempty\"`\n", capitalize(propertyName), propertyName)
			}
		case "integer":
			if required {
				fmt.Fprintf(s, capitalize(propertyName)+" int64 `json:"+`"`+propertyName+`,omitempty"`+"`\n")
			} else {
				fmt.Fprintf(s, capitalize(propertyName)+" *int64 `json:"+`"`+propertyName+`,omitempty"`+"`\n")
			}
		case "string":
			if required {
				fmt.Fprintf(s, capitalize(propertyName)+" string `json:"+`"`+propertyName+`,omitempty"`+"`\n")
			} else {
				fmt.Fprintf(s, capitalize(propertyName)+" *string `json:"+`"`+propertyName+`,omitempty"`+"`\n")
			}
		case "array":
			itemstype := lexicon.resolveItemsType(defname, propertyName, property.Items)
			fmt.Fprintf(s, capitalize(propertyName)+" []"+itemstype+" `json:"+`"`+propertyName+`,omitempty"`+"`\n")
		case "ref":
			reftype := lexicon.resolveRefType(property.Ref)
			fmt.Fprintf(s, capitalize(propertyName)+reftype+" `json:"+`"`+propertyName+`,omitempty"`+"`\n")
		case "unknown":
			if required {
				fmt.Fprintf(s, capitalize(propertyName)+" any `json:"+`"`+propertyName+`,omitempty"`+"`\n")
			} else {
				fmt.Fprintf(s, capitalize(propertyName)+" *any `json:"+`"`+propertyName+`,omitempty"`+"`\n")
			}
		case "blob":
			if required {
				fmt.Fprintf(s, capitalize(propertyName)+" *common.Blob `json:"+`"`+propertyName+`"`+"`\n")
			} else {
				fmt.Fprintf(s, capitalize(propertyName)+" *common.Blob `json:"+`"`+propertyName+`,omitempty"`+"`\n")
			}
		case "union":
			uniontype := lexicon.resolveUnionFieldType(defname, propertyName)
			fmt.Fprintf(s, capitalize(propertyName)+" *"+uniontype+" `json:"+`"`+propertyName+`,omitempty"`+"`\n")
		case "bytes":
			if required {
				fmt.Fprintf(s, capitalize(propertyName)+" []byte `json:"+`"`+propertyName+`"`+"`\n")
			} else {
				fmt.Fprintf(s, capitalize(propertyName)+" *[]byte `json:"+`"`+propertyName+`,omitempty"`+"`\n")
			}
		case "cid-link":
			if required {
				fmt.Fprintf(s, capitalize(propertyName)+" string `json:"+`"`+propertyName+`,omitempty"`+"`\n")
			} else {
				fmt.Fprintf(s, capitalize(propertyName)+" *string `json:"+`"`+propertyName+`,omitempty"`+"`\n")
			}
		default:
			fmt.Fprintf(s, "// FIXME skipping unsupported property type "+propertyName+" "+property.Type+" "+fmt.Sprintf("required=%t %+v", required, property)+"\n")
		}
	}
	fmt.Fprintf(s, "}\n\n")
}

func (lexicon *Lexicon) renderDependencies(s *strings.Builder, defname string, properties map[string]Property, required []string) {
	propertyNames := sortedPropertyNames(properties)
	for _, propertyName := range propertyNames {
		property := properties[propertyName]
		switch property.Type {
		case "union":
			uniontype := lexicon.resolveUnionFieldType(defname, propertyName)
			fmt.Fprintf(s, "type %s struct {\n", uniontype)
			for _, ref := range property.Refs {
				fieldname := lexicon.unionFieldName(ref)
				fieldtype := lexicon.unionFieldType(ref)
				fmt.Fprintf(s, "%s %s\n", fieldname, fieldtype)
			}
			fmt.Fprintf(s, "}\n\n")
			fmt.Fprintf(s, "func (m *%s) UnmarshalJSON(data []byte) error {\n", uniontype)
			fmt.Fprintf(s, "recordType := common.LexiconTypeFromJSONBytes(data)\n")
			fmt.Fprintf(s, "switch recordType {\n")
			for _, ref := range property.Refs {
				fieldname := lexicon.unionFieldName(ref)
				fieldtype := lexicon.unionFieldType(ref)[1:] // strip leading *
				fmt.Fprintf(s, "case \"%s\":\n", ref)
				fmt.Fprintf(s, "m.%s = &%s{}\n", fieldname, fieldtype)
				fmt.Fprintf(s, "json.Unmarshal(data, m.%s)\n", fieldname)

			}
			fmt.Fprintf(s, "}\n")
			fmt.Fprintf(s, "return nil\n")
			fmt.Fprintf(s, "}\n\n")
			fmt.Fprintf(s, "func (m %s) MarshalJSON() ([]byte, error) {\n", uniontype)
			for _, ref := range property.Refs {
				fieldname := lexicon.unionFieldName(ref)
				fmt.Fprintf(s, "if m.%s != nil {\n", fieldname)
				fmt.Fprintf(s, "return json.Marshal(m.%s)\n", fieldname)
				fmt.Fprintf(s, "} else ")
			}
			fmt.Fprintf(s, "{ return []byte(\"{}\"), nil }\n")
			fmt.Fprintf(s, "}\n\n")
		case "array":
			if property.Items.Type == "union" {
				uniontype := lexicon.resolveUnionFieldType(defname, propertyName) + "_Elem"
				fmt.Fprintf(s, "type %s struct {\n", uniontype)
				for _, ref := range property.Items.Refs {
					fieldname := lexicon.unionFieldName(ref)
					fieldtype := lexicon.unionFieldType(ref)
					fmt.Fprintf(s, "%s %s\n", fieldname, fieldtype)
				}
				fmt.Fprintf(s, "}\n\n")

				//fmt.Fprintf(s, "/*\n")
				fmt.Fprintf(s, "func (m *%s) UnmarshalJSON(data []byte) error {\n", uniontype)
				fmt.Fprintf(s, "recordType := common.LexiconTypeFromJSONBytes(data)\n")
				fmt.Fprintf(s, "switch recordType {\n")
				for _, ref := range property.Items.Refs {
					fieldname := lexicon.unionFieldName(ref)
					fieldtype := lexicon.unionFieldType(ref)[1:] // strip leading *
					refType := ref
					if refType[0] == '#' {
						refType = lexicon.Id + refType
					}
					fmt.Fprintf(s, "case \"%s\":\n", refType)
					fmt.Fprintf(s, "m.%s = &%s{}\n", fieldname, fieldtype)
					fmt.Fprintf(s, "json.Unmarshal(data, m.%s)\n", fieldname)

				}
				fmt.Fprintf(s, "}\n")
				fmt.Fprintf(s, "return nil\n")
				fmt.Fprintf(s, "}\n\n")
				fmt.Fprintf(s, "func (m %s) MarshalJSON() ([]byte, error) {\n", uniontype)
				for _, ref := range property.Items.Refs {
					fieldname := lexicon.unionFieldName(ref)
					fmt.Fprintf(s, "if m.%s != nil {\n", fieldname)
					fmt.Fprintf(s, "return json.Marshal(m.%s)\n", fieldname)
					fmt.Fprintf(s, "} else ")
				}
				fmt.Fprintf(s, "{ return []byte(\"{}\"), nil }\n")
				fmt.Fprintf(s, "}\n\n")
				//fmt.Fprintf(s, "*/\n")
			}
		}
	}
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
		reflexicon := LookupLexicon(id)
		if reflexicon != nil {
			refdef := reflexicon.Lookup(tag)
			if refdef != nil {
				reftype = refdef.Type
			}
		}
		idparts := strings.Split(id, ".")
		if len(idparts) != 4 {
			return "/* FIXME skipping union field with invalid id " + fmt.Sprintf("%+v", ref) + " */ string"
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
		return "/* FIXME defaulting on unparseable union field ref " + fmt.Sprintf("%+v", ref) + " */ string"
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
		reflexicon := LookupLexicon(id)
		if reflexicon != nil {
			refdef := reflexicon.Lookup(tag)
			if refdef != nil {
				reftype = refdef.Type
			}
		}
		name := symbolForID(id)
		if tag != "main" {
			name += "_" + capitalize(tag)
		}
		if reftype == "array" {
			return "[]" + name + "_Elem"
		}
		return "*" + name
	} else {
		return "/* FIXME defaulting on unparsable union field ref " + fmt.Sprintf("%+v", ref) + " */ string"
	}
}
