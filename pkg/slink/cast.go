package slink

import (
	"encoding/json"
	"log"
)

func CastStringToPointer(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func CastInt64ToPointer(v int64) *int64 {
	if v == 0 {
		return nil
	}
	return &v
}

func CastBoolToPointer(v bool) *bool {
	if v == false {
		return nil
	}
	return &v
}

func CastAnyToStruct[T any](v *any) *T {
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

func CastAnyToArray[T any](v *any) []*T {
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
