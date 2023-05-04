package core

import "fmt"

// BuildCMakeModule builds a module that builds a CMake project
func BuildCMakeModule(module Module) (string, error) {
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
