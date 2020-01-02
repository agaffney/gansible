package cmd

import (
	"fmt"
)

type adhocFlags struct {
	core flagsCore
}

func init() {
	addCommand(command{
		name:       `gansible`,
		entrypoint: adhocMain,
	})
}

func adhocMain() {
	fmt.Printf("running gansible\n")
	flags := adhocFlags{}
	addFlagsCore(&flags.core)
	parseFlags()
}
