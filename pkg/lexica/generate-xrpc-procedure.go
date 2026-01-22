package lexica

import (
	"fmt"
	"strings"
)

func (lexicon *Lexicon) generateProcedure(s *strings.Builder, defname string, def *Def) {
	fmt.Fprintf(s, "const "+defname+"_Description = "+`"`+def.Description+`"`+"\n\n")
	if def.Input != nil && def.Input.Encoding == "application/json" &&
		def.Output != nil && def.Output.Encoding == "application/json" {
		lexicon.generateStruct(s, defname+"_Input", "", def.Input.Schema.Properties, def.Input.Schema.Required, false)

		if def.Output.Schema.Type == "ref" {
			if def.Output.Schema.Ref[0] == '#' {
				fmt.Fprintf(s, "type %s = %s\n", defname+"_Output", defname+"_"+capitalize(def.Output.Schema.Ref[1:]))
			} else {
				parts := strings.Split(def.Output.Schema.Ref, "#")
				fmt.Fprintf(s, "type %s = %s\n", defname+"_Output", symbolForID(parts[0])+"_"+capitalize(parts[1]))
			}
		} else {
			lexicon.generateStruct(s, defname+"_Output", "", def.Output.Schema.Properties, def.Output.Schema.Required, false)
		}
		fmt.Fprintf(s, "// %s\n", def.Description)
		fmt.Fprintf(s, "func "+defname+"(ctx context.Context, c common.Client, input *"+defname+"_Input) (*"+defname+"_Output"+", error) {\n")
		fmt.Fprintf(s, "var output "+defname+"_Output"+"\n")
		fmt.Fprintf(s, `if err := c.Do(ctx, common.Procedure, "`+def.Input.Encoding+`", "`+lexicon.Id+`", nil, input, &output); err != nil {`+"\n")
		fmt.Fprintf(s, "return nil, err\n")
		fmt.Fprintf(s, "}\n")
		fmt.Fprintf(s, "return &output, nil\n")
		fmt.Fprintf(s, "}\n\n")
	} else if def.Input == nil &&
		def.Output != nil && def.Output.Encoding == "application/json" {
		lexicon.generateStruct(s, defname+"_Output", "", def.Output.Schema.Properties, def.Output.Schema.Required, false)
		fmt.Fprintf(s, "// %s\n", def.Description)
		fmt.Fprintf(s, "func "+defname+"(ctx context.Context, c common.Client) (*"+defname+"_Output"+", error) {\n")
		fmt.Fprintf(s, "var output "+defname+"_Output"+"\n")
		fmt.Fprintf(s, `if err := c.Do(ctx, common.Procedure, "", "`+lexicon.Id+`", nil, nil, &output); err != nil {`+"\n")
		fmt.Fprintf(s, "return nil, err\n")
		fmt.Fprintf(s, "}\n")
		fmt.Fprintf(s, "return &output, nil\n")
		fmt.Fprintf(s, "}\n\n")
	} else if def.Input != nil && def.Input.Encoding == "application/json" &&
		def.Output == nil {
		lexicon.generateStruct(s, defname+"_Input", "", def.Input.Schema.Properties, def.Input.Schema.Required, false)
		fmt.Fprintf(s, "// %s\n", def.Description)
		fmt.Fprintf(s, "func "+defname+"(ctx context.Context, c common.Client, input *"+defname+"_Input) error {\n")
		fmt.Fprintf(s, `return c.Do(ctx, common.Procedure, "`+def.Input.Encoding+`", "`+lexicon.Id+`", nil, input, nil)`+"\n")
		fmt.Fprintf(s, "}\n\n")
	} else if def.Input == nil && def.Output == nil {
		fmt.Fprintf(s, "// %s\n", def.Description)
		fmt.Fprintf(s, "func "+defname+"(ctx context.Context, c common.Client) error {\n")
		fmt.Fprintf(s, `return c.Do(ctx, common.Procedure, "", "`+lexicon.Id+`", nil, nil, nil)`+"\n")
		fmt.Fprintf(s, "}\n\n")
	} else if def.Input != nil &&
		def.Output != nil && def.Output.Encoding == "application/json" {
		if def.Output.Schema.Type == "ref" {
			if def.Output.Schema.Ref[0] == '#' {
				fmt.Fprintf(s, "type %s = %s\n", defname+"_Output", defname+"_"+capitalize(def.Output.Schema.Ref[1:]))
			} else {
				parts := strings.Split(def.Output.Schema.Ref, "#")
				fmt.Fprintf(s, "type %s = %s\n", defname+"_Output", symbolForID(parts[0])+"_"+capitalize(parts[1]))
			}
		} else {
			lexicon.generateStruct(s, defname+"_Output", "", def.Output.Schema.Properties, def.Output.Schema.Required, false)
		}
		fmt.Fprintf(s, "// %s\n", def.Description)
		fmt.Fprintf(s, "func "+defname+"(ctx context.Context, c common.Client, input io.Reader) (*"+defname+"_Output"+", error) {\n")
		fmt.Fprintf(s, "var output "+defname+"_Output"+"\n")
		fmt.Fprintf(s, `if err := c.Do(ctx, common.Procedure, "`+def.Input.Encoding+`", "`+lexicon.Id+`", nil, input, &output); err != nil {`+"\n")
		fmt.Fprintf(s, "return nil, err\n")
		fmt.Fprintf(s, "}\n")
		fmt.Fprintf(s, "return &output, nil\n")
		fmt.Fprintf(s, "}\n\n")
	} else if def.Input != nil && def.Output == nil {
		fmt.Fprintf(s, "// %s\n", def.Description)
		fmt.Fprintf(s, "func "+defname+"(ctx context.Context, c common.Client, input io.Reader) error {\n")
		fmt.Fprintf(s, `if err := c.Do(ctx, common.Procedure, "`+def.Input.Encoding+`", "`+lexicon.Id+`", nil, input, nil); err != nil {`+"\n")
		fmt.Fprintf(s, "return err\n")
		fmt.Fprintf(s, "}\n")
		fmt.Fprintf(s, "return nil\n")
		fmt.Fprintf(s, "}\n\n")
	} else {
		fmt.Fprintf(s, "// FIXME skipping procedure with unhandled types\n")
	}
}
