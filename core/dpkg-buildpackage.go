package core

import "fmt"

// BuildDpkgModule builds a module that builds a dpkg project
// and installs the resulting .deb package
func BuildDpkgBuildPkgModule(module Module) (string, error) {
	cmd := fmt.Sprintf(
		"cd /sources/%s && dpkg-buildpackage -d -us -uc -b",
		module.Name,
	)

	for _, path := range module.Source.Paths {
		cmd += fmt.Sprintf(" && dpkg -i ../%s*.deb && apt install -f", path)
	}

	cmd += " && apt clean"
	return cmd, nil
}
