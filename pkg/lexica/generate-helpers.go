package lexica

import (
	"sort"
	"strings"
)

func capitalize(s string) string {
	return strings.ToUpper(s[0:1]) + s[1:]
}

func sortedPropertyNames(properties map[string]Property) []string {
	var propnames []string
	for propname := range properties {
		propnames = append(propnames, propname)
	}
	sort.Strings(propnames)
	return propnames
}
