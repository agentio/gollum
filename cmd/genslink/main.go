package main

import (
	"os"

	"github.com/agentio/slink/cli"
)

func main() {
	if err := cli.Cmd().Execute(); err != nil {
		os.Exit(1)
	}
}
