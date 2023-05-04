package core

import "fmt"

// BuildDpkgModule builds a module that installs a .deb package
func BuildDpkgModule(module Module) (string, error) {
	cmd := ""
	for i, path := range module.Paths {
		cmd += fmt.Sprintf(" dpkg -i /sources/%s && apt install -f", path)
		if i < len(module.Paths)-1 {
			cmd += " &&"
		}
	}
	return cmd, nil
}
