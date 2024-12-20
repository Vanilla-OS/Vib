package main

import (
	"C"
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/vanilla-os/vib/api"
)

// Configuration for a CMake module
type CMakeModule struct {
	Name       string            `json:"name"`
	Type       string            `json:"type"`
	BuildVars  map[string]string `json:"buildvars"`
	BuildFlags string            `json:"buildflags"`
	Source     api.Source
}

// Provide plugin information as a JSON string
//
//export PlugInfo
func PlugInfo() *C.char {
	plugininfo := &api.PluginInfo{Name: "cmake", Type: api.BuildPlugin, UseContainerCmds: false}
	pluginjson, err := json.Marshal(plugininfo)
	if err != nil {
		return C.CString(fmt.Sprintf("ERROR: %s", err.Error()))
	}
	return C.CString(string(pluginjson))
}

// Generate a shell command to build a CMake project based on the provided module and recipe.
// Download and move the source, set up build variables and flags, and construct the CMake build command.
//
//export BuildModule
func BuildModule(moduleInterface *C.char, recipeInterface *C.char) *C.char {
	var module *CMakeModule
	var recipe *api.Recipe

	err := json.Unmarshal([]byte(C.GoString(moduleInterface)), &module)
	if err != nil {
		return C.CString(fmt.Sprintf("ERROR: %s", err.Error()))
	}

	err = json.Unmarshal([]byte(C.GoString(recipeInterface)), &recipe)
	if err != nil {
		return C.CString(fmt.Sprintf("ERROR: %s", err.Error()))
	}

	err = api.DownloadSource(recipe, module.Source, module.Name)
	if err != nil {
		return C.CString(fmt.Sprintf("ERROR: %s", err.Error()))
	}
	err = api.MoveSource(recipe.DownloadsPath, recipe.SourcesPath, module.Source, module.Name)
	if err != nil {
		return C.CString(fmt.Sprintf("ERROR: %s", err.Error()))
	}
	buildVars := map[string]string{}
	for k, v := range module.BuildVars {
		buildVars[k] = v
	}

	buildFlags := ""
	if module.BuildFlags != "" {
		buildFlags = " " + module.BuildFlags
	}

	cmd := fmt.Sprintf(
		"cd /sources/%s && mkdir -p build && cd build && cmake ../%s && make",
		filepath.Join(recipe.SourcesPath, api.GetSourcePath(module.Source, module.Name)),
		buildFlags,
	)

	return C.CString(cmd)
}

func main() { fmt.Println("This plugin is not meant to run standalone!") }
