package lexica

import (
	"fmt"
	"strings"
)

func (lexicon *Lexicon) generateProcedure(s *strings.Builder, defname string, def *Def) {
	s.WriteString("const " + defname + "_Description = " + `"` + def.Description + `"` + "\n\n")
	if def.Input != nil && def.Input.Encoding == "application/json" &&
		def.Output != nil && def.Output.Encoding == "application/json" {
		lexicon.generateStruct(s, defname+"_Input", "", def.Input.Schema.Properties, def.Input.Schema.Required, false)

		if def.Output.Schema.Type == "ref" {
			if def.Output.Schema.Ref[0] == '#' {
				fmt.Fprintf(s, "type %s = %s\n", defname+"_Output", defname+"_"+capitalize(def.Output.Schema.Ref[1:]))
			} else {
				parts := strings.Split(def.Output.Schema.Ref, "#")
				fmt.Fprintf(s, "type %s = %s\n", defname+"_Output", idPrefix(parts[0])+"_"+capitalize(parts[1]))
			}
		} else {
			lexicon.generateStruct(s, defname+"_Output", "", def.Output.Schema.Properties, def.Output.Schema.Required, false)
		}
		s.WriteString("// " + def.Description + "\n")
		s.WriteString("func " + defname + "(ctx context.Context, c common.Client, input *" + defname + "_Input) (*" + defname + "_Output" + ", error) {\n")
		s.WriteString("var output " + defname + "_Output" + "\n")
		s.WriteString(`if err := c.Do(ctx, common.Procedure, "` + def.Input.Encoding + `", "` + lexicon.Id + `", nil, input, &output); err != nil {` + "\n")
		s.WriteString("return nil, err\n")
		s.WriteString("}\n")
		s.WriteString("return &output, nil\n")
		s.WriteString("}\n\n")
	} else if def.Input == nil &&
		def.Output != nil && def.Output.Encoding == "application/json" {
		lexicon.generateStruct(s, defname+"_Output", "", def.Output.Schema.Properties, def.Output.Schema.Required, false)
		s.WriteString("// " + def.Description + "\n")
		s.WriteString("func " + defname + "(ctx context.Context, c common.Client) (*" + defname + "_Output" + ", error) {\n")
		s.WriteString("var output " + defname + "_Output" + "\n")
		s.WriteString(`if err := c.Do(ctx, common.Procedure, "", "` + lexicon.Id + `", nil, nil, &output); err != nil {` + "\n")
		s.WriteString("return nil, err\n")
		s.WriteString("}\n")
		s.WriteString("return &output, nil\n")
		s.WriteString("}\n\n")
	} else if def.Input != nil && def.Input.Encoding == "application/json" &&
		def.Output == nil {
		lexicon.generateStruct(s, defname+"_Input", "", def.Input.Schema.Properties, def.Input.Schema.Required, false)
		s.WriteString("// " + def.Description + "\n")
		s.WriteString("func " + defname + "(ctx context.Context, c common.Client, input *" + defname + "_Input) error {\n")
		s.WriteString(`return c.Do(ctx, common.Procedure, "` + def.Input.Encoding + `", "` + lexicon.Id + `", nil, input, nil)` + "\n")
		s.WriteString("}\n\n")
	} else if def.Input == nil && def.Output == nil {
		s.WriteString("// " + def.Description + "\n")
		s.WriteString("func " + defname + "(ctx context.Context, c common.Client) error {\n")
		s.WriteString(`return c.Do(ctx, common.Procedure, "", "` + lexicon.Id + `", nil, nil, nil)` + "\n")
		s.WriteString("}\n\n")
	} else if def.Input != nil &&
		def.Output != nil && def.Output.Encoding == "application/json" {
		if def.Output.Schema.Type == "ref" {
			if def.Output.Schema.Ref[0] == '#' {
				fmt.Fprintf(s, "type %s = %s\n", defname+"_Output", defname+"_"+capitalize(def.Output.Schema.Ref[1:]))
			} else {
				parts := strings.Split(def.Output.Schema.Ref, "#")
				fmt.Fprintf(s, "type %s = %s\n", defname+"_Output", idPrefix(parts[0])+"_"+capitalize(parts[1]))
			}
		} else {
			lexicon.generateStruct(s, defname+"_Output", "", def.Output.Schema.Properties, def.Output.Schema.Required, false)
		}
		s.WriteString("// " + def.Description + "\n")
		s.WriteString("func " + defname + "(ctx context.Context, c common.Client, input io.Reader) (*" + defname + "_Output" + ", error) {\n")
		s.WriteString("var output " + defname + "_Output" + "\n")
		s.WriteString(`if err := c.Do(ctx, common.Procedure, "` + def.Input.Encoding + `", "` + lexicon.Id + `", nil, input, &output); err != nil {` + "\n")
		s.WriteString("return nil, err\n")
		s.WriteString("}\n")
		s.WriteString("return &output, nil\n")
		s.WriteString("}\n\n")
	} else if def.Input != nil && def.Output == nil {
		s.WriteString("// " + def.Description + "\n")
		s.WriteString("func " + defname + "(ctx context.Context, c common.Client, input io.Reader) error {\n")
		s.WriteString(`if err := c.Do(ctx, common.Procedure, "` + def.Input.Encoding + `", "` + lexicon.Id + `", nil, input, nil); err != nil {` + "\n")
		s.WriteString("return err\n")
		s.WriteString("}\n")
		s.WriteString("return nil\n")
		s.WriteString("}\n\n")
	} else {
		s.WriteString("// FIXME skipping procedure with unhandled types\n")
	}
}
