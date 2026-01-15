package lexica

import (
	"strings"

	"github.com/charmbracelet/log"
)

func (lexicon *Lexicon) generateDef(def *Def, name string, prefix string) string {
	var defname string
	if name == "main" {
		defname = prefix
	} else {
		defname = prefix + "_" + capitalize(name)
	}
	var s strings.Builder
	switch def.Type {
	case "query":
		s.WriteString(lexicon.generateQuery(defname, def))
	case "procedure":
		s.WriteString(lexicon.generateProcedure(defname, def))
	case "object":
		s.WriteString(lexicon.generateStruct(defname, def.Description, def.Properties, def.Required))
	case "string":
		s.WriteString("type " + defname + " string\n")
	case "record":
		s.WriteString(lexicon.generateStruct(defname, def.Description, def.Properties, def.Required))
	case "array":
		s.WriteString("type " + defname + "_Elem struct {\n")
		s.WriteString("}\n\n")
	case "token":
		s.WriteString("// " + def.Description + "\n")
		s.WriteString("const " + defname + " string = " + `"` + name + `"` + "\n\n")
	default:
		log.Warnf("skipping %s.%s (type %s)", lexicon.Id, name, def.Type)
	}
	return s.String()
}
