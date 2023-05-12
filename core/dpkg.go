package core

import "fmt"

// BuildDpkgModule builds a module that installs a .deb package
func BuildDpkgModule(module Module) (string, error) {
	cmd := ""
	for _, path := range module.Source.Paths {
		cmd += fmt.Sprintf(" dpkg -i /sources/%s && apt install -f && ", path)
	}

	cmd += " && apt clean"
	return cmd, nil
}
