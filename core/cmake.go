package core

import "fmt"

type CMakeModule struct {
	Name       string            `json:"name"`
	Type       string            `json:"name"`
	BuildVars  map[string]string `json:"name"`
	BuildFlags string            `json:"name"`
}

// BuildCMakeModule builds a module that builds a CMake project
func BuildCMakeModule(module CMakeModule) (string, error) {
	buildVars := map[string]string{}
	for k, v := range module.BuildVars {
		buildVars[k] = v
	}

	buildFlags := ""
	if module.BuildFlags != "" {
		buildFlags = " " + module.BuildFlags
	}

	cmd := fmt.Sprintf(
		"cd /sources/%s && mkdir -p build && cd build && cmake ..%s && make",
		module.Name,
		buildFlags,
	)

	return cmd, nil
}
