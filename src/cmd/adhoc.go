package cmd

import (
	"fmt"
)

func init() {
	addCommand(&command{
		name:       `gansible`,
		entrypoint: adhocMain,
	})
}

func adhocMain() {
	fmt.Printf("running gansible\n")
}
