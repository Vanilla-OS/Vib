package main

import (
	"C"
	"encoding/json"
	"fmt"

	"github.com/vanilla-os/vib/api"
)

// Configuration for a Cargo module
type CargoModule struct {
	Name       string            `json:"name"`
	Type       string            `json:"type"`
	BuildFlags string            `json:"buildflags"`
	Source     api.Source
}

// Provide plugin information as a JSON string
//
//export PlugInfo
func PlugInfo() *C.char {
	plugininfo := &api.PluginInfo{Name: "cargo-deb", Type: api.BuildPlugin, UseContainerCmds: false}
	pluginjson, err := json.Marshal(plugininfo)
	if err != nil {
		return C.CString(fmt.Sprintf("ERROR: %s", err.Error()))
	}
	return C.CString(string(pluginjson))
}

// Generate a shell command to build a Cargo project based on the provided module and recipe.
// Download and move the source, set up build variables and flags, and construct the Cargo build command.
//
//export BuildModule
func BuildModule(moduleInterface *C.char, recipeInterface *C.char, arch *C.char) *C.char {
	var module *CargoModule
	var recipe *api.Recipe

	err := json.Unmarshal([]byte(C.GoString(moduleInterface)), &module)
	if err != nil {
		return C.CString(fmt.Sprintf("ERROR: %s", err.Error()))
	}

	err = json.Unmarshal([]byte(C.GoString(recipeInterface)), &recipe)
	if err != nil {
		return C.CString(fmt.Sprintf("ERROR: %s", err.Error()))
	}

	if !api.TestArch(module.Source.OnlyArches, C.GoString(arch)) {
		return C.CString("")
	}

	err = api.DownloadSource(recipe, module.Source, module.Name)
	if err != nil {
		return C.CString(fmt.Sprintf("ERROR: %s", err.Error()))
	}
	err = api.MoveSource(recipe.DownloadsPath, recipe.SourcesPath, module.Source, module.Name)
	if err != nil {
		return C.CString(fmt.Sprintf("ERROR: %s", err.Error()))
	}

	buildFlags := ""
	if module.BuildFlags != "" {
		buildFlags = " " + module.BuildFlags
	}

	cmd := fmt.Sprintf(
		"cd /sources/%s && cargo install cargo-deb && cargo deb %s && apt install -y ./target/debian/*",
		api.GetSourcePath(module.Source, module.Name),
		buildFlags,
	)

	return C.CString(cmd)
}

func main() { fmt.Println("This plugin is not meant to run standalone!") }
