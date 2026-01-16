package lexica

import (
	"fmt"
	"strings"
)

func (lexicon *Lexicon) generateQuery(s *strings.Builder, defname string, def *Def) {
	s.WriteString("const " + defname + "_Description = " + `"` + def.Description + `"` + "\n\n")
	if def.Output != nil && def.Output.Encoding == "application/json" {
		lexicon.generateStruct(s, defname+"_Output", "", def.Output.Schema.Properties, def.Output.Schema.Required)
		params := ""
		paramsok := false
		if def.Parameters != nil && def.Parameters.Type == "params" {
			params, paramsok = parseQueryParameters(def.Parameters)
		}
		s.WriteString("// " + def.Description + "\n")
		s.WriteString("func " + defname + "(ctx context.Context, c xrpc.Client" + params + ") (*" + defname + "_Output" + ", error) {\n")
		s.WriteString("var output " + defname + "_Output" + "\n")
		s.WriteString("params := map[string]interface{}{\n")
		if paramsok {
			for parameterName := range def.Parameters.Properties {
				s.WriteString(`"` + parameterName + `":` + parameterName + ",\n")
			}
		}
		s.WriteString("}\n")
		s.WriteString(`if err := c.Do(ctx, xrpc.Query, "", "` + lexicon.Id + `", params, nil, &output); err != nil {` + "\n")
		s.WriteString("return nil, err\n")
		s.WriteString("}\n")
		s.WriteString("return &output, nil\n")
		s.WriteString("}\n\n")
	} else if def.Output != nil {
		s.WriteString(fmt.Sprintf("// FIXME skipping query with output type %s\n", def.Output.Encoding))
		params := ""
		paramsok := false
		if def.Parameters != nil && def.Parameters.Type == "params" {
			params, paramsok = parseQueryParameters(def.Parameters)
		}
		s.WriteString("// " + def.Description + "\n")
		s.WriteString("func " + defname + "(ctx context.Context, c xrpc.Client" + params + ") ([]byte, error) {\n")
		s.WriteString("params := map[string]interface{}{\n")
		if paramsok {
			for parameterName := range def.Parameters.Properties {
				s.WriteString(`"` + parameterName + `":` + parameterName + ",\n")
			}
		}
		s.WriteString("}\n")
		s.WriteString("var output []byte\n")
		s.WriteString(`if err := c.Do(ctx, xrpc.Query, "", "` + lexicon.Id + `", params, nil, &output); err != nil {` + "\n")
		s.WriteString("return nil, err\n")
		s.WriteString("}\n")
		s.WriteString("return output, nil\n")
		s.WriteString("}\n\n")
	} else {
		s.WriteString(fmt.Sprintf("// FIXME skipping query with no output %+v\n", def))
	}
}

func parseQueryParameters(parameters *Parameters) (string, bool) {
	var parms []string
	propertyNames := sortedPropertyNames(parameters.Properties)
	for _, propertyName := range propertyNames {
		propertyValue := parameters.Properties[propertyName]
		declaration := propertyName + " "
		switch propertyValue.Type {
		case "integer":
			declaration += "int64"
		case "string":
			declaration += "string"
		case "boolean":
			declaration += "bool"
		case "array":
			if propertyValue.Items.Type == "string" {
				declaration += "[]string"
			} else {
				return "/* FIXME failing on unsupported parameter array value type: " + propertyValue.Items.Type + " */", false
			}
		default:
			return "/* FIXME failing on unsupported parameter value type: " + propertyValue.Type + "*/", false
		}
		parms = append(parms, declaration)
	}
	return ", " + strings.Join(parms, ", "), true
}
