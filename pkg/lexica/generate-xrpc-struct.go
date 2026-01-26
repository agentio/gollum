package lexica

import (
	"fmt"
	"slices"
	"strings"
)

func (lexicon *Lexicon) generateStructAndDependencies(s *strings.Builder, defname, description string, properties map[string]Property, required []string, isRecord bool, name string) {
	fmt.Fprintf(s, "const %s_Description = \"%s\"\n", defname, description)
	lexicon.generateStruct(s, defname, properties, required, isRecord, name)
	lexicon.generateDependencies(s, defname, properties, required)
}

func (lexicon *Lexicon) generateStruct(s *strings.Builder, defname string, properties map[string]Property, required []string, isRecord bool, name string) {
	if isRecord {
		fmt.Fprintf(s, "// %s is a record with Lexicon type %s#%s\n", defname, lexicon.Id, name)
	}
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
				fmt.Fprintf(s, "%s int64 `json:\"%s\"`\n", capitalize(propertyName), propertyName)
			} else {
				fmt.Fprintf(s, "%s *int64 `json:\"%s,omitempty\"`\n", capitalize(propertyName), propertyName)
			}
		case "string":
			if required {
				fmt.Fprintf(s, "%s string `json:\"%s\"`\n", capitalize(propertyName), propertyName)
			} else {
				fmt.Fprintf(s, "%s *string `json:\"%s,omitempty\"`\n", capitalize(propertyName), propertyName)
			}
		case "array":
			itemstype := lexicon.resolveItemsType(defname, propertyName, property.Items)
			fmt.Fprintf(s, "%s []%s `json:\"%s,omitempty\"`\n", capitalize(propertyName), itemstype, propertyName)
		case "ref":
			reftype := lexicon.resolveRefType(property.Ref)
			fmt.Fprintf(s, "%s %s `json:\"%s,omitempty\"`\n", capitalize(propertyName), reftype, propertyName)
		case "unknown":
			if required {
				fmt.Fprintf(s, "%s any `json:\"%s\"`\n", capitalize(propertyName), propertyName)
			} else {
				fmt.Fprintf(s, "%s *any `json:\"%s,omitempty\"`\n", capitalize(propertyName), propertyName)
			}
		case "blob":
			if required {
				fmt.Fprintf(s, "%s *slink.Blob `json:\"%s\"`\n", capitalize(propertyName), propertyName)
			} else {
				fmt.Fprintf(s, "%s *slink.Blob `json:\"%s,omitempty\"`\n", capitalize(propertyName), propertyName)
			}
		case "union":
			uniontype := lexicon.resolveUnionFieldType(defname, propertyName)
			fmt.Fprintf(s, "%s *%s `json:\"%s,omitempty\"`\n", capitalize(propertyName), uniontype, propertyName)
		case "bytes":
			if required {
				fmt.Fprintf(s, "%s []byte `json:\"%s\"`\n", capitalize(propertyName), propertyName)
			} else {
				fmt.Fprintf(s, "%s *[]byte `json:\"%s,omitempty\"`\n", capitalize(propertyName), propertyName)
			}
		case "cid-link":
			if required {
				fmt.Fprintf(s, "%s string `json:\"%s\"`\n", capitalize(propertyName), propertyName)
			} else {
				fmt.Fprintf(s, "%s *string `json:\"%s,omitempty\"`\n", capitalize(propertyName), propertyName)
			}
		default:
			fmt.Fprintf(s, "// FIXME skipping unsupported property type %s %s required=%t %+v\n", propertyName, property.Type, required, property)
		}
	}
	fmt.Fprintf(s, "}\n\n")
}

func (lexicon *Lexicon) generateDependencies(s *strings.Builder, defname string, properties map[string]Property, required []string) {
	propertyNames := sortedPropertyNames(properties)
	for _, propertyName := range propertyNames {
		property := properties[propertyName]
		switch property.Type {
		case "union":
			uniontype := lexicon.resolveUnionFieldType(defname, propertyName)
			lexicon.generateUnion(s, uniontype, property.Refs)
		case "array":
			if property.Items.Type == "union" {
				uniontype := lexicon.resolveUnionFieldType(defname, propertyName) + "_Elem"
				lexicon.generateUnion(s, uniontype, property.Items.Refs)
			}
		}
	}
}

func (lexicon *Lexicon) generateUnion(s *strings.Builder, uniontype string, refs []string) {
	uniontypeinner := uniontype + "_Inner"
	fmt.Fprintf(s, "// %s is a union with these possible values:\n", uniontype)
	for _, ref := range refs {
		if ref == "#main" {
			ref = lexicon.Id
		}
		if ref[0] == '#' {
			ref = lexicon.Id + ref
		}
		fmt.Fprintf(s, "// - %s (%s)\n", lexicon.unionFieldElementType(ref), ref)
	}
	fmt.Fprintf(s, "type %s struct {\n", uniontype)
	fmt.Fprintf(s, "Inner %s", uniontypeinner)
	fmt.Fprintf(s, "}\n\n")

	fmt.Fprintf(s, "type %s interface {\n", uniontypeinner)
	fmt.Fprintf(s, "is%s()", uniontype)
	fmt.Fprintf(s, "}\n\n")

	for _, ref := range refs {
		if ref == "#main" {
			ref = lexicon.Id
		}
		if ref[0] == '#' {
			ref = lexicon.Id + ref
		}
		wrappertype := lexicon.wrapperTypeName(uniontype, ref)
		wrappedtype := lexicon.unionFieldType(ref)
		fmt.Fprintf(s, "type %s struct {\n", wrappertype)
		fmt.Fprintf(s, "Value %s\n", wrappedtype)
		fmt.Fprintf(s, "}\n\n")
		fmt.Fprintf(s, "func (t %s) is%s() {}\n\n", wrappertype, uniontype)
	}

	fmt.Fprintf(s, "func (u %s) MarshalJSON() ([]byte, error) {\n", uniontype)
	fmt.Fprintf(s, "switch v := u.Inner.(type) {\n")
	for _, ref := range refs {
		if ref == "#main" {
			ref = lexicon.Id
		}
		if ref[0] == '#' {
			ref = lexicon.Id + ref
		}
		wrappertype := lexicon.wrapperTypeName(uniontype, ref)
		fmt.Fprintf(s, "case %s:\n", wrappertype)
		fmt.Fprintf(s, "return slink.MarshalWithType(\"%s\", v.Value)\n", ref)
	}
	fmt.Fprintf(s, "default:\n")
	fmt.Fprintf(s, "	return nil, fmt.Errorf(\"unsupported type %%T\", v)\n")
	fmt.Fprintf(s, "}\n")
	fmt.Fprintf(s, "}\n")

	fmt.Fprintf(s, "func (u *%s) UnmarshalJSON(data []byte) error {\n", uniontype)
	fmt.Fprintf(s, "switch slink.LexiconTypeFromJSONBytes(data) {\n")
	for _, ref := range refs {
		if ref == "#main" {
			ref = lexicon.Id
		}
		if ref[0] == '#' {
			ref = lexicon.Id + ref
		}
		wrappertype := lexicon.wrapperTypeName(uniontype, ref)
		fmt.Fprintf(s, "case \"%s\":\n", ref)
		fmt.Fprintf(s, "var v %s\n", lexicon.unionFieldElementType(ref))
		fmt.Fprintf(s, "if err := json.Unmarshal(data, &v); err != nil {return err}\n")
		fmt.Fprintf(s, "u.Inner = %s{Value: &v}\n", wrappertype)
		fmt.Fprintf(s, "return nil\n")
	}
	fmt.Fprintf(s, "default:\n")
	fmt.Fprintf(s, "return nil\n")
	fmt.Fprintf(s, "}\n")
	fmt.Fprintf(s, "}\n\n")

}

func (lexicon *Lexicon) generateOldUnion(s *strings.Builder, uniontype string, refs []string) {
	fmt.Fprintf(s, "// union type, only one field must be set\n")
	fmt.Fprintf(s, "type %s struct {\n", uniontype)
	for _, ref := range refs {
		fieldname := lexicon.unionFieldName(ref)
		fieldtype := lexicon.unionFieldType(ref)
		fmt.Fprintf(s, "%s %s\n", fieldname, fieldtype)
	}
	fmt.Fprintf(s, "}\n\n")
	fmt.Fprintf(s, "func (m *%s) UnmarshalJSON(data []byte) error {\n", uniontype)
	fmt.Fprintf(s, "recordType := slink.LexiconTypeFromJSONBytes(data)\n")
	fmt.Fprintf(s, "switch recordType {\n")
	for _, ref := range refs {
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
	for _, ref := range refs {
		fieldname := lexicon.unionFieldName(ref)
		fmt.Fprintf(s, "if m.%s != nil {\n", fieldname)
		fmt.Fprintf(s, "return json.Marshal(m.%s)\n", fieldname)
		fmt.Fprintf(s, "} else ")
	}
	fmt.Fprintf(s, "{ return []byte(\"{}\"), nil }\n")
	fmt.Fprintf(s, "}\n\n")
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
		name := symbolForID(id)
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

func (lexicon *Lexicon) unionFieldElementType(ref string) string {
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
		return name
	} else {
		return "/* FIXME defaulting on unparsable union field ref " + fmt.Sprintf("%+v", ref) + " */ string"
	}
}

func (lexicon *Lexicon) wrapperTypeName(uniontype, ref string) string {
	return uniontype + "__" + lexicon.unionFieldElementType(ref)
}
