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
		s.WriteString("type " + defname + "_Elem struct {\n")
		s.WriteString("}\n\n")
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
