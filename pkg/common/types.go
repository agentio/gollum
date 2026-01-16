package common

func StringPointerOrNil(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
