package core

import (
	"fmt"

	"github.com/vanilla-os/vib/api"
)

type CMakeModule struct {
	Name       string            `json:"name"`
	Type       string            `json:"type"`
	BuildVars  map[string]string `json:"buildvars"`
	BuildFlags string            `json:"buildflags"`
	Source     api.Source
}

// BuildCMakeModule builds a module that builds a CMake project
func BuildCMakeModule(moduleInterface interface{}, recipe *api.Recipe) (string, error) {
	module := moduleInterface.(CMakeModule)
	err := api.DownloadSource(recipe.DownloadsPath, module.Source)
	if err != nil {
		return "", err
	}
	err = api.MoveSource(recipe.DownloadsPath, module.Source)
	if err != nil {
		return "", err
	}
	buildVars := map[string]string{}
	for k, v := range module.BuildVars {
		buildVars[k] = v
	}

	buildFlags := ""
	if module.BuildFlags != "" {
		buildFlags = " " + module.BuildFlags
	}

	cmd := fmt.Sprintf(
		"cd /sources/%s && mkdir -p build && cd build && cmake ..%s && make",
		module.Name,
		buildFlags,
	)

	return cmd, nil
}
