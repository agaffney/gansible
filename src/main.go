package main

import (
	"fmt"
	"github.com/agaffney/gansible/cmd"
	"os"
	"path"
)

func main() {
	cmdName := path.Base(os.Args[0])
	entrypoint := cmd.GetEntrypoint(cmdName)
	if entrypoint == nil {
		fmt.Printf("Unknown command '%s'\n", cmdName)
		os.Exit(1)
	}
	entrypoint()
}
