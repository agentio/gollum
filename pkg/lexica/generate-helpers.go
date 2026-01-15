package lexica

import (
	"strings"
)

func capitalize(s string) string {
	return strings.ToUpper(s[0:1]) + s[1:]
}
