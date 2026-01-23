package slink

func TruncateShort(s string) string {
	return TruncateToLength(s, 80)
}

func TruncateToLength(s string, maxlen int) string {
	maxlen = maxlen - 3
	if len(s) < maxlen {
		return s
	}
	return s[0:maxlen] + "..."
}
