package cmd

import (
	"fmt"
	"github.com/agaffney/gansible/python"
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
	python.Init()
}
