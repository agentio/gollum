package common

func StringPointerOrNil(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func Int64PointerOrNil(v int64) *int64 {
	if v == 0 {
		return nil
	}
	return &v
}

func Truncate(s string) string {
	const maxlen = 77
	if len(s) < maxlen {
		return s
	}
	return s[0:maxlen] + "..."
}
