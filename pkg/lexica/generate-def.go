package lexica

import (
	"github.com/charmbracelet/log"
)

func (lexicon *Lexicon) generateDef(def *Def, name string, prefix string) string {

	var defname string
	if name == "main" {
		defname = prefix
	} else {
		defname = prefix + "_" + capitalize(name)
	}

	s := ""
	switch def.Type {
	case "query":
		s += lexicon.generateQuery(defname, def)
	case "procedure":
		s += lexicon.generateProcedure(defname, def)
	case "object":
		s += "// " + def.Description + "\n"
		s += "type " + defname + " struct {\n"
		s += renderProperties(lexicon, defname, def.Properties, def.Required)
		s += "}\n\n"
		s += renderDependentTypes(lexicon, defname, def.Properties, def.Required)
	case "string":
		s += "type " + defname + " string\n"
	case "record":
		s += "// " + def.Description + "\n"
		s += "type " + defname + " struct {\n"
		s += renderProperties(lexicon, defname, def.Properties, def.Required)
		s += "}\n\n"
		s += renderDependentTypes(lexicon, defname, def.Properties, def.Required)
	case "array":
		s += "type " + defname + "_Elem struct {\n"
		s += "}\n\n"
	case "token":
		s += "// " + def.Description + "\n"
		s += "const " + defname + " string = " + `"` + name + `"` + "\n\n"
	default:
		log.Warnf("skipping %s.%s (type %s)", lexicon.Id, name, def.Type)
	}
	return s
}
