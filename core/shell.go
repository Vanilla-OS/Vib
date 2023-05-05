package core

import "errors"

func BuildShellModule(module Module) (string, error) {
	if len(module.Commands) == 0 {
		return "", errors.New("no commands specified")
	}

	cmd := ""
	for i, command := range module.Commands {
		cmd += command
		if i < len(module.Commands)-1 {
			cmd += " && "
		}
	}

	return cmd, nil
}
