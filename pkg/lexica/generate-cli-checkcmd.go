package lexica

func (lexicon *Lexicon) generateCheckCommands(root string) {
	for defname, def := range lexicon.Defs {
		switch def.Type {
		case "record":
			lexicon.generateCheckCommandForDef(root, defname, def)
		}
	}
}

func (lexicon *Lexicon) generateCheckCommandForDef(root, defname string, def *Def) {

}
