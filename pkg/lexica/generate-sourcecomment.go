package lexica

import (
	"encoding/json"
	"strings"
)

func (lexicon *Lexicon) generateSourceComment(s *strings.Builder) {
	filter := func(s string) string {
		return strings.ReplaceAll(s, "*/*", "[ANY]")
	}
	b, _ := json.MarshalIndent(lexicon, "", "  ")
	s.WriteString("/*\n")
	s.WriteString(filter(string(b)) + "\n")
	s.WriteString("*/\n")
}
