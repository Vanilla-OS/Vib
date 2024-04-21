package core

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/vanilla-os/vib/api"
)

type MesonModule struct {
	Name   string
	Type   string
	Source api.Source
}

// BuildMesonModule builds a module that builds a Meson project
func BuildMesonModule(moduleInterface interface{}, recipe *api.Recipe) (string, error) {
	var module MesonModule
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

	// Since the downloaded source goes through checksum verification already
	// it is safe to simply use the specified checksum from the module definition
	tmpDir := fmt.Sprintf("/tmp/%s-%s", module.Source.Checksum, module.Name)
	cmd := fmt.Sprintf(
		"cd /sources/%s && meson %s && ninja -C %s && ninja -C %s install",
		module.Name,
		tmpDir,
		tmpDir,
		tmpDir,
	)

	return cmd, nil
}
