package lexica

import (
	"fmt"
	"strings"
)

func (lexicon *Lexicon) generateSubscribe(s *strings.Builder, defname string, def *Def) {
	fmt.Fprintf(s, "const %s_Description = \"%s\"\n\n", defname, def.Description)
	if def.Message != nil {
		params := ""
		paramsok := false
		if def.Parameters != nil && def.Parameters.Type == "params" {
			params, paramsok = parseQueryParameters(def.Parameters)
		}
		fmt.Fprintf(s, "type %s_Handler func(m *%s_Message ) error\n\n", defname, defname)
		params += fmt.Sprintf(", fn %s_Handler", defname)
		fmt.Fprintf(s, "// %s\n", def.Description)
		fmt.Fprintf(s, "func %s(ctx context.Context, c slink.Client%s) (error) {\n", defname, params)
		fmt.Fprintf(s, "log.Error(\"FIXME: %s doesn't work yet\")\n", defname)
		fmt.Fprintf(s, "params := map[string]any{\n")
		if paramsok {
			for _, parameterName := range sortedPropertyNames(def.Parameters.Properties) {
				fmt.Fprintf(s, "\"%s\":%s,\n", parameterName, parameterName)
			}
		}
		fmt.Fprintf(s, "}\n")
		fmt.Fprintf(s, "var output []byte\n")
		fmt.Fprintf(s, "if err := c.Do(ctx, slink.Query, \"\", \"%s\", params, nil, &output); err != nil {\n", lexicon.Id)
		fmt.Fprintf(s, "return err\n")
		fmt.Fprintf(s, "}\n")
		fmt.Fprintf(s, "_ = output\n")
		fmt.Fprintf(s, "return nil\n")
		fmt.Fprintf(s, "}\n\n")
		lexicon.generateUnion(s, defname+"_Message", def.Message.Schema.Refs)
	} else {
		fmt.Fprintf(s, "// FIXME skipping subscription with no message %+v\n", def)
	}
}
