package core

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/vanilla-os/vib/api"
)

type DpkgModule struct {
	Name   string `json:"name"`
	Type   string `json:"string"`
	Source api.Source
}

// BuildDpkgModule builds a module that installs a .deb package
func BuildDpkgModule(moduleInterface interface{}, recipe *api.Recipe) (string, error) {
	var module CMakeModule
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
	cmd := ""
	for _, path := range module.Source.Paths {
		cmd += fmt.Sprintf(" dpkg -i /sources/%s && apt install -f && ", path)
	}

	cmd += " && apt clean"
	return cmd, nil
}
