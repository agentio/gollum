package lexica

func (lexicon *Lexicon) generateProcedure(defname string, def *Def) string {
	s := ""
	if def.Output != nil && def.Output.Encoding == "application/json" &&
		def.Input != nil && def.Input.Encoding == "application/json" {
		// input
		s += "type " + defname + "_Input struct {\n"
		s += renderProperties(lexicon, defname+"_Input", def.Input.Schema.Properties, def.Input.Schema.Required)
		s += "}\n\n"
		s += renderDependentTypes(lexicon, defname+"_Input", def.Input.Schema.Properties, def.Input.Schema.Required)
		// output
		s += "type " + defname + "_Output struct {\n"
		s += renderProperties(lexicon, defname+"_Output", def.Output.Schema.Properties, def.Output.Schema.Required)
		s += "}\n\n"
		s += renderDependentTypes(lexicon, defname+"_Output", def.Output.Schema.Properties, def.Output.Schema.Required)
		// func
		s += "// " + def.Description + "\n"
		s += "func " + defname + "(ctx context.Context, c xrpc.Client, input *" + defname + "_Input) (*" + defname + "_Output" + ", error) {\n"
		s += "  var output " + defname + "_Output" + "\n"
		s += `if err := c.Do(ctx, xrpc.Procedure, "application/json", "` + lexicon.Id + `", nil, input, &output); err != nil {` + "\n"
		s += "return nil, err\n"
		s += "}\n"
		s += "  return &output, nil\n"
		s += "}\n\n"
	} else if def.Input == nil && def.Output != nil && def.Output.Encoding == "application/json" {
		// output
		s += "type " + defname + "_Output struct {\n"
		s += renderProperties(lexicon, defname+"_Output", def.Output.Schema.Properties, def.Output.Schema.Required)
		s += "}\n\n"
		s += renderDependentTypes(lexicon, defname+"_Output", def.Output.Schema.Properties, def.Output.Schema.Required)
		// func
		s += "// " + def.Description + "\n"
		s += "func " + defname + "(ctx context.Context, c xrpc.Client) (*" + defname + "_Output" + ", error) {\n"
		s += "  var output " + defname + "_Output" + "\n"
		s += `if err := c.Do(ctx, xrpc.Procedure, "application/json", "` + lexicon.Id + `", nil, nil, &output); err != nil {` + "\n"
		s += "return nil, err\n"
		s += "}\n"
		s += "  return &output, nil\n"
		s += "}\n\n"
	} else if def.Output == nil && def.Input != nil && def.Input.Encoding == "application/json" {
		// input
		s += "type " + defname + "_Input struct {\n"
		s += renderProperties(lexicon, defname+"_Input", def.Input.Schema.Properties, def.Input.Schema.Required)
		s += "}\n\n"
		s += renderDependentTypes(lexicon, defname+"_Input", def.Input.Schema.Properties, def.Input.Schema.Required)
		// func
		s += "// " + def.Description + "\n"
		s += "func " + defname + "(ctx context.Context, c xrpc.Client, input *" + defname + "_Input) error {\n"
		s += `return c.Do(ctx, xrpc.Procedure, "", "` + lexicon.Id + `", nil, input, nil)` + "\n"
		s += "}\n\n"
	} else if def.Output == nil && def.Input == nil {
		// func
		s += "// " + def.Description + "\n"
		s += "func " + defname + "(ctx context.Context, c xrpc.Client) error {\n"
		s += `return c.Do(ctx, xrpc.Procedure, "", "` + lexicon.Id + `", nil, nil, nil)` + "\n"
		s += "}\n\n"
	} else {
		s += "// FIXME (procedure with unhandled types)\n"
	}
	return s
}
