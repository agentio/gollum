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
