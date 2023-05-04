package core

import "fmt"

// BuildDpkgModule builds a module that builds a dpkg project
// and installs the resulting .deb package
func BuildDpkgBuildPkgModule(module Module) (string, error) {
	cmd := fmt.Sprintf(
		"cd /sources/%s && dpkg-buildpackage -us -uc -b &&",
		module.Name,
	)

	for i, path := range module.Paths {
		cmd += fmt.Sprintf(" dpkg -i ../%s*.deb && apt install -f", path)
		if i < len(module.Paths)-1 {
			cmd += " &&"
		}
	}

	return cmd, nil
}
