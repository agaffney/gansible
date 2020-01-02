package cmd

type entrypointFunc func()

type command struct {
	name       string
	entrypoint entrypointFunc
}

var commands = []command{}

func addCommand(cmd command) {
	commands = append(commands, cmd)
}

func GetEntrypoint(name string) entrypointFunc {
	for _, cmd := range commands {
		if cmd.name == name {
			return cmd.entrypoint
		}
	}
	return nil
}
