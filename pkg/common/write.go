package common

import (
	"encoding/json"
	"io"
)

func Write(w io.Writer, v any) error {
	switch v := v.(type) {
	case []byte:
		w.Write(v)
	default:
		b, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			return err
		}
		w.Write(b)
		w.Write([]byte("\n"))
	}
	return nil
}
