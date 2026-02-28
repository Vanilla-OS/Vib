package core

import (
	"fmt"
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
	Cleanup  []string
}

// Build shell module commands and return them as a single string
//
// Returns: Concatenated shell commands or an error if any step fails
func BuildShellModule(module interface{}, recipe *api.Recipe, cleanup []string, arch string) (string, error) {
	var shellModule ShellModule

	if err := mapstructure.Decode(module, &shellModule); err != nil {
		return "", err
	}

	for _, source := range shellModule.Sources {
		if api.TestArch(source.OnlyArches, arch) {
			if strings.TrimSpace(source.Type) != "" {
				err := api.DownloadSource(recipe, source, shellModule.Name)
				if err != nil {
					return "", err
				}
				err = api.MoveSource(recipe.DownloadsPath, recipe.SourcesPath, source, shellModule.Name)
				if err != nil {
					return "", err
				}
			}
		}
	}

	if len(shellModule.Commands) == 0 {
		return "", fmt.Errorf("no commands specified")
	}

	var cmd strings.Builder
	_, err := fmt.Fprintf(&cmd, "RUN --mount=source=sources/%s,target=/sources/%s,rw\nRUN ", shellModule.Name, shellModule.Name)
	if err != nil {
		panic(fmt.Sprintf("Fprintf failed during build of shell module `%s`", shellModule.Name))
	}
	for i, command := range shellModule.Commands {
		cmd.WriteString(command)
		if i < len(shellModule.Commands)-1 {
			cmd.WriteString(" && ")
		}
	}
	cmd.WriteString(api.GetCleanupSuffix(append(cleanup, shellModule.Cleanup...)))

	return cmd.String(), nil
}
