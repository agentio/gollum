package slink

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

func RespondWithJSON(w http.ResponseWriter, response any) {
	b, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "failed to serialize response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	fmt.Fprintf(w, "%s\n", string(b))
}
