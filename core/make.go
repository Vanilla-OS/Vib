package core

import (
	"github.com/mitchellh/mapstructure"
	"github.com/vanilla-os/vib/api"
)

type MakeModule struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Source api.Source
}

// BuildMakeModule builds a module that builds a Make project
func BuildMakeModule(moduleInterface interface{}, recipe *api.Recipe) (string, error) {
	var module MakeModule
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

	return "cd /sources/" + api.GetSourcePath(module.Source, module.Name) + " && make && make install", nil
}
