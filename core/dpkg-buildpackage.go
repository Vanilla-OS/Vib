package core

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/vanilla-os/vib/api"
)

type DpkgBuildModule struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Source api.Source
}

// BuildDpkgModule builds a module that builds a dpkg project
// and installs the resulting .deb package
func BuildDpkgBuildPkgModule(moduleInterface interface{}, recipe *api.Recipe) (string, error) {
	var module DpkgBuildModule
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

	cmd := fmt.Sprintf(
		"cd /sources/%s && dpkg-buildpackage -d -us -uc -b",
		module.Name,
	)

	for _, path := range module.Source.Paths {
		cmd += fmt.Sprintf(" && apt install -y --allow-downgrades ../%s*.deb", path)
	}

	cmd += " && apt clean"
	return cmd, nil
}
