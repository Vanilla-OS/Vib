package core

import "errors"

type ShellModule struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Commands []string
}

func BuildShellModule(module ShellModule) (string, error) {
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
