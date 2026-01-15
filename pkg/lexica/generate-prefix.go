package lexica

import (
	"strings"
)

func idPrefix(id string) string {
	id = strings.TrimPrefix(id, "com.atproto.") // put these symbols in the top-level namespace
	var s strings.Builder
	for _, part := range strings.Split(id, ".") {
		s.WriteString(capitalize(part))
	}
	return s.String()
}
