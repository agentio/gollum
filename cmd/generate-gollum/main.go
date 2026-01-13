package main

import (
	"os"

	"github.com/agentio/gollum/cmd/generate-gollum/generate"
)

func main() {
	if err := generate.Cmd().Execute(); err != nil {
		os.Exit(1)
	}
}
