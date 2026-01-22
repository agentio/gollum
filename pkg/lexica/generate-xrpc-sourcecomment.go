package lexica

import (
	"encoding/json"
	"fmt"
	"strings"
)

func (lexicon *Lexicon) generateSourceComment(s *strings.Builder) {
	filter := func(s string) string {
		return strings.ReplaceAll(s, "*/*", "[ANY]")
	}
	b, _ := json.MarshalIndent(lexicon, "", "  ")
	fmt.Fprintf(s, "/*\n%s\n*/\n", filter(string(b)))
}
