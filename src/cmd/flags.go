package cmd

import (
	"flag"
)

type flagsCore struct {
	inventory string
}

func addFlagsCore(flags *flagsCore) {
	flag.StringVar(&flags.inventory, "inventory", "/etc/ansible/hosts", "inventory file")
}

func parseFlags() {
	flag.Parse()
}
