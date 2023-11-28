package core

import "github.com/vanilla-os/vib/api"

type MakeModule struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Source api.Source
}

// BuildMakeModule builds a module that builds a Make project
func BuildMakeModule(moduleInterface interface{}, recipe *api.Recipe) (string, error) {
	module := moduleInterface.(MakeModule)

	err := api.DownloadSource(recipe.DownloadsPath, module.Source)
	if err != nil {
		return "", err
	}
	err = api.MoveSource(recipe.DownloadsPath, module.Source)
	if err != nil {
		return "", err
	}

	return "cd /sources/" + module.Name + " && make && make install", nil
}
