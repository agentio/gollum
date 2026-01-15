package lexica

import (
	"strings"
)

func codePrefix(id string) string {
	id = strings.TrimPrefix(id, "com.atproto.") // put this in the top-level namespace
	var s strings.Builder
	for _, part := range strings.Split(id, ".") {
		s.WriteString(capitalize(part))
	}
	return s.String()
}
