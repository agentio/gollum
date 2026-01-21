package common

import (
	"encoding/json"
	"fmt"
	"os"
)

func ReadJSONFile(name string) (any, error) {
	if name == "" {
		return nil, nil
	}
	b, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}
	var mapValue map[string]any
	err = json.Unmarshal(b, &mapValue)
	if err == nil {
		return mapValue, nil
	}
	var arrayValue []any
	err = json.Unmarshal(b, &arrayValue)
	if err == nil {
		return arrayValue, nil
	}
	return nil, fmt.Errorf("unable to unmarshal contents of %s, is it json?", name)
}
