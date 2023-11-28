package core

import (
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/vanilla-os/vib/api"
)

type ShellModule struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Commands []string
}

func BuildShellModule(moduleInterface interface{}, _ *api.Recipe) (string, error) {
	var module ShellModule
	mapstructure.Decode(moduleInterface, &module)
	fmt.Println(moduleInterface)
	fmt.Println(module)
	fmt.Println(module.Commands)
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
