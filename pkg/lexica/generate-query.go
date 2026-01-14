package lexica

import "fmt"

func (lexicon *Lexicon) generateQuery(defname string, def *Def) string {
	s := ""
	if def.Output != nil && def.Output.Encoding == "application/json" {
		// output
		s += "type " + defname + "_Output struct {\n"
		s += renderProperties(lexicon, defname+"_Output", def.Output.Schema.Properties, def.Output.Schema.Required)
		s += "}\n\n"
		s += renderDependentTypes(lexicon, defname+"_Output", def.Output.Schema.Properties, def.Output.Schema.Required)
		// parameters
		params := ""
		paramsok := false
		if def.Parameters != nil && def.Parameters.Type == "params" {
			params, paramsok = parseParameters(def.Parameters)
		}
		// func
		s += "// " + def.Description + "\n"
		s += "func " + defname + "(ctx context.Context, c xrpc.Client" + params + ") (*" + defname + "_Output" + ", error) {\n"
		s += "  var output " + defname + "_Output" + "\n"
		s += "params := map[string]interface{}{\n"
		if paramsok {
			for parameterName, _ := range def.Parameters.Properties {
				s += `"` + parameterName + `":` + parameterName + ",\n"
			}
		}
		s += "}\n"
		s += `if err := c.Do(ctx, xrpc.Query, "", "` + lexicon.Id + `", params, nil, &output); err != nil {` + "\n"
		s += "return nil, err\n"
		s += "}\n"
		s += "  return &output, nil\n"
		s += "}\n\n"
	} else if def.Output != nil {
		s += fmt.Sprintf("// FIXME (query, output type %s)\n", def.Output.Encoding)
	} else {
		s += fmt.Sprintf("// FIXME (query, no output) %+v\n", def)
	}

	return s
}
