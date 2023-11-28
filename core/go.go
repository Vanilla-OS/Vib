package core

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/vanilla-os/vib/api"
)

type GoModule struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	Source     api.Source
	BuildVars  map[string]string
	BuildFlags string
}

// BuildGoModule builds a module that builds a Go project
// buildVars are used to customize the build command
// like setting the output binary name and location
func BuildGoModule(moduleInterface interface{}, recipe *api.Recipe) (string, error) {
	var module GoModule
	err := mapstructure.Decode(moduleInterface, &module)
	if err != nil {
		return "", err
	}
	err = api.DownloadSource(recipe.DownloadsPath, module.Source, module.Name)
	if err != nil {
		return "", err
	}
	err = api.MoveSource(recipe.DownloadsPath, recipe.SourcesPath, module.Source, module.Name)
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

	buildVars["GO_OUTPUT_BIN"] = module.Name
	if module.BuildVars["GO_OUTPUT_BIN"] != "" {
		buildVars["GO_OUTPUT_BIN"] = module.BuildVars["GO_OUTPUT_BIN"]
	}

	cmd := fmt.Sprintf(
		"cd /sources/%s && go build%s -o %s",
		module.Name,
		buildFlags,
		buildVars["GO_OUTPUT_BIN"],
	)

	return cmd, nil
}
