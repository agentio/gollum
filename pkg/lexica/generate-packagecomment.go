package lexica

import (
	"fmt"
	"strings"
)

func packageComment(s *strings.Builder, packagename string) {
	fmt.Fprintf(s, "// Code generated ... DO NOT EDIT.\n\n")
	fmt.Fprintf(s, "// Package %s is generated from Lexicon source files by slink.\n", packagename)
	fmt.Fprintf(s, "// Get slink at https://github.com/agentio/slink.\n")
}
