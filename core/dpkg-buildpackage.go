package core

import (
	"fmt"
	"github.com/vanilla-os/vib/api"
)

type DpkgBuildModule struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Source api.Source
}

// BuildDpkgModule builds a module that builds a dpkg project
// and installs the resulting .deb package
func BuildDpkgBuildPkgModule(moduleInterface interface{}, _ *api.Recipe) (string, error) {
	module := moduleInterface.(DpkgBuildModule)
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
