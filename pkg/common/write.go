package common

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func Write(w io.Writer, output string, v any) error {
	switch v := v.(type) {
	case []byte:
		switch output {
		case "":
			fmt.Fprintf(w, "(%d bytes, use \"-o -\" to write to stdout)\n", len(v))
		case "-":
			w.Write(v)
		default:
			os.WriteFile(output, v, 0644)
		}
	default:
		b, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			return err
		}
		switch output {
		case "", "-":
			w.Write(b)
			w.Write([]byte("\n"))
		default:
			os.WriteFile(output, b, 0644)
		}
	}
	return nil
}
