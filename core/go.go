package core

import "fmt"

// BuildGoModule builds a module that builds a Go project
// buildVars are used to customize the build command
// like setting the output binary name and location
func BuildGoModule(module Module) (string, error) {
	buildVars := map[string]string{}
	for k, v := range module.BuildVars {
		buildVars[k] = v
	}

	buildFlags := ""
	if module.BuildFlags != "" {
		buildFlags = " " + module.BuildFlags
	}

	buildVars["GO_OUTPUT_BIN"] = module.Name
	if module.BuildVars["GO_OUTPUT_BIN"] != "" {
		buildVars["GO_OUTPUT_BIN"] = module.BuildVars["GO_OUTPUT_BIN"]
	}

	cmd := fmt.Sprintf(
		"cd /sources/%s && go build%s -o %s%s",
		module.Name,
		buildFlags,
		buildVars["GO_OUTPUT_BIN"],
		module.Source.Module,
	)

	return cmd, nil
}
