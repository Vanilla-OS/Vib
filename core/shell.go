package core

import (
	"errors"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/vanilla-os/vib/api"
)

// Configuration for shell modules
type ShellModule struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Sources  []api.Source
	Commands []string
}

// Build shell module commands and return them as a single string
//
// Returns: Concatenated shell commands or an error if any step fails
func BuildShellModule(moduleInterface interface{}, recipe *api.Recipe) (string, error) {
	var module ShellModule
	err := mapstructure.Decode(moduleInterface, &module)
	if err != nil {
		return "", err
	}

	for _, source := range module.Sources {
		if strings.TrimSpace(source.Type) != "" {
			err := api.DownloadSource(recipe.DownloadsPath, source, module.Name)
			if err != nil {
				return "", err
			}
			err = api.MoveSource(recipe.DownloadsPath, recipe.SourcesPath, source, module.Name)
			if err != nil {
				return "", err
			}
		}
	}

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
