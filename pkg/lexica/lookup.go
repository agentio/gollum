package lexica

func LookupLexicon(id string) *Lexicon {
	for _, lexicon := range _catalog.Lexicons {
		if lexicon.Id == id {
			return lexicon
		}
	}
	return nil
}

func (lexicon *Lexicon) Lookup(id string) *Def {
	d, ok := lexicon.Defs[id]
	if !ok {
		return nil
	}
	return d
}
