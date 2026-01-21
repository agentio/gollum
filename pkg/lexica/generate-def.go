package lexica

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/log"
)

func (lexicon *Lexicon) generateDef(s *strings.Builder, def *Def, name string, prefix string) {
	var defname string
	if name == "main" {
		defname = prefix
	} else {
		defname = prefix + "_" + capitalize(name)
	}
	switch def.Type {
	case "query":
		lexicon.generateQuery(s, defname, def)
	case "procedure":
		lexicon.generateProcedure(s, defname, def)
	case "object":
		lexicon.generateStruct(s, defname, def.Description, def.Properties, def.Required, true)
	case "string":
		s.WriteString("type " + defname + " string\n")
	case "record":
		fmt.Fprintf(s, "const %s_Description = \"%s\"\n", defname, def.Description)
		lexicon.generateStruct(s, defname, def.Description, def.Record.Properties, def.Record.Required, true)
	case "array":
		if def.Items.Type == "union" {
			uniontype := defname + "_Elem"
			s.WriteString("type " + defname + "_Elem struct {\n")
			for _, ref := range def.Items.Refs {
				fieldname := lexicon.unionFieldName(ref)
				fieldtype := lexicon.unionFieldType(ref)
				s.WriteString(fieldname + " " + fieldtype + "\n")
			}
			s.WriteString("}\n\n")
			fmt.Fprintf(s, "func (m *%s) UnmarshalJSON(data []byte) error {\n", uniontype)
			fmt.Fprintf(s, "recordType := common.LexiconTypeFromJSONBytes(data)\n")
			fmt.Fprintf(s, "switch recordType {\n")
			for _, ref := range def.Items.Refs {
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
			for _, ref := range def.Items.Refs {
				fieldname := lexicon.unionFieldName(ref)
				fmt.Fprintf(s, "if m.%s != nil {\n", fieldname)
				fmt.Fprintf(s, "return json.Marshal(m.%s)\n", fieldname)
				fmt.Fprintf(s, "} else ")
			}
			fmt.Fprintf(s, "{ return []byte(\"{}\"), nil }\n")
			fmt.Fprintf(s, "}\n\n")
		} else {
			fmt.Fprintf(s, "// FIXME: ungenerated array %+v\n", def)
		}
	case "token":
		s.WriteString("// " + def.Description + "\n")
		s.WriteString("const " + defname + " string = " + `"` + name + `"` + "\n\n")
	case "permission-set":
		s.WriteString("// CHECKME skipping permission set " + defname + "\n")
	case "subscription":
		s.WriteString("// CHECKME skipping subscription " + defname + "\n")
	default:
		log.Warnf("skipping %s.%s (type %s)", lexicon.Id, name, def.Type)
	}
}
