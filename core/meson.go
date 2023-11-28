package core

import (
	"fmt"
	"github.com/vanilla-os/vib/api"

	"github.com/google/uuid"
)

type MesonModule struct {
	Name   string
	Type   string
	Source api.Source
}

// BuildMesonModule builds a module that builds a Meson project
func BuildMesonModule(moduleInterface interface{}, recipe *api.Recipe) (string, error) {
	module := moduleInterface.(MesonModule)
	tmpDir := "/tmp/" + uuid.New().String()

	err := api.DownloadSource(recipe.DownloadsPath, module.Source)
	if err != nil {
		return "", err
	}
	err = api.MoveSource(recipe.DownloadsPath, module.Source)
	if err != nil {
		return "", err
	}

	cmd := fmt.Sprintf(
		"cd /sources/%s && meson %s && ninja -C %s && ninja -C %s install",
		module.Name,
		tmpDir,
		tmpDir,
		tmpDir,
	)

	return cmd, nil
}
