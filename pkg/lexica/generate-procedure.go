package lexica

func (lexicon *Lexicon) generateProcedure(defname string, def *Def) string {
	var s string
	if def.Input != nil && def.Input.Encoding == "application/json" &&
		def.Output != nil && def.Output.Encoding == "application/json" {
		s += lexicon.generateStruct(defname+"_Input", "", def.Input.Schema.Properties, def.Input.Schema.Required)
		s += lexicon.generateStruct(defname+"_Output", "", def.Output.Schema.Properties, def.Output.Schema.Required)
		s += "// " + def.Description + "\n"
		s += "func " + defname + "(ctx context.Context, c xrpc.Client, input *" + defname + "_Input) (*" + defname + "_Output" + ", error) {\n"
		s += "var output " + defname + "_Output" + "\n"
		s += `if err := c.Do(ctx, xrpc.Procedure, "application/json", "` + lexicon.Id + `", nil, input, &output); err != nil {` + "\n"
		s += "return nil, err\n"
		s += "}\n"
		s += "return &output, nil\n"
		s += "}\n\n"
	} else if def.Input == nil &&
		def.Output != nil && def.Output.Encoding == "application/json" {
		s += lexicon.generateStruct(defname+"_Output", "", def.Output.Schema.Properties, def.Output.Schema.Required)
		s += "// " + def.Description + "\n"
		s += "func " + defname + "(ctx context.Context, c xrpc.Client) (*" + defname + "_Output" + ", error) {\n"
		s += "var output " + defname + "_Output" + "\n"
		s += `if err := c.Do(ctx, xrpc.Procedure, "application/json", "` + lexicon.Id + `", nil, nil, &output); err != nil {` + "\n"
		s += "return nil, err\n"
		s += "}\n"
		s += "return &output, nil\n"
		s += "}\n\n"
	} else if def.Input != nil && def.Input.Encoding == "application/json" &&
		def.Output == nil {
		s += lexicon.generateStruct(defname+"_Input", "", def.Input.Schema.Properties, def.Input.Schema.Required)
		s += "// " + def.Description + "\n"
		s += "func " + defname + "(ctx context.Context, c xrpc.Client, input *" + defname + "_Input) error {\n"
		s += `return c.Do(ctx, xrpc.Procedure, "", "` + lexicon.Id + `", nil, input, nil)` + "\n"
		s += "}\n\n"
	} else if def.Input == nil && def.Output == nil {
		s += "// " + def.Description + "\n"
		s += "func " + defname + "(ctx context.Context, c xrpc.Client) error {\n"
		s += `return c.Do(ctx, xrpc.Procedure, "", "` + lexicon.Id + `", nil, nil, nil)` + "\n"
		s += "}\n\n"
	} else {
		s += "// FIXME (procedure with unhandled types)\n"
	}
	return s
}
