package common

import (
	"encoding/json"

	"github.com/charmbracelet/log"
)

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

func BoolPointerOrNil(v bool) *bool {
	if v == false {
		return nil
	}
	return &v
}

func Truncate(s string) string {
	return TruncateToLength(s, 80)
}

func TruncateToLength(s string, maxlen int) string {
	maxlen = maxlen - 3
	if len(s) < maxlen {
		return s
	}
	return s[0:maxlen] + "..."
}

func LexiconTypeFromJSONBytes(data []byte) string {
	type TypedRecord struct {
		LexiconTypeID string `json:"$type"`
	}
	var record TypedRecord
	err := json.Unmarshal(data, &record)
	if err != nil {
		return ""
	}
	return record.LexiconTypeID
}

type Blob struct {
	LexiconTypeID string `json:"$type,omitempty"`
	Ref           Link   `json:"ref,omitempty"`
	MimeType      string `json:"mimeType,omitempty"`
	Size          int64  `json:"size"`
}

type Link struct {
	LexiconLink string `json:"$link"`
}

func CastIntoStructType[T any](v *any) *T {
	if v == nil {
		return nil
	}
	var result T
	b, err := json.Marshal(v)
	if err != nil {
		log.Printf("%+v", err)
		return nil
	}
	err = json.Unmarshal(b, &result)
	if err != nil {
		log.Printf("%+v", err)
		return nil
	}
	return &result
}

func CastIntoArrayType[T any](v *any) []*T {
	var result []*T
	b, err := json.Marshal(v)
	if err != nil {
		log.Printf("%+v", err)
		return nil
	}
	err = json.Unmarshal(b, &result)
	if err != nil {
		log.Printf("%+v", err)
		return nil
	}
	return result
}
