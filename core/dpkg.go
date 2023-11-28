package core

import (
	"fmt"
	"github.com/vanilla-os/vib/api"
)

type DpkgModule struct {
	Name   string `json:"name"`
	Type   string `json:"string"`
	Source api.Source
}

// BuildDpkgModule builds a module that installs a .deb package
func BuildDpkgModule(moduleInterface interface{}, _ *api.Recipe) (string, error) {
	module := moduleInterface.(DpkgModule)
	cmd := ""
	for _, path := range module.Source.Paths {
		cmd += fmt.Sprintf(" dpkg -i /sources/%s && apt install -f && ", path)
	}

	cmd += " && apt clean"
	return cmd, nil
}
