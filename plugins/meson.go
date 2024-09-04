package main

import (
	"C"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/vanilla-os/vib/api"
)

// Configuration for building a Meson project
type MesonModule struct {
	Name   string
	Type   string
	BuildFlags []string            `json:"buildflags"`
	Source api.Source
}

// Provide plugin information as a JSON string
//
//export PlugInfo
func PlugInfo() *C.char {
	plugininfo := &api.PluginInfo{Name: "meson", Type: api.BuildPlugin}
	pluginjson, err := json.Marshal(plugininfo)
	if err != nil {
		return C.CString(fmt.Sprintf("ERROR: %s", err.Error()))
	}
	return C.CString(string(pluginjson))
}

// Generate a command to build a Meson project. Handle source downloading, moving,
// and use Meson and Ninja build tools with a temporary build directory based on the checksum.
//
//export BuildModule
func BuildModule(moduleInterface *C.char, recipeInterface *C.char) *C.char {
	var module *MesonModule
	var recipe *api.Recipe

	err := json.Unmarshal([]byte(C.GoString(moduleInterface)), &module)
	if err != nil {
		return C.CString(fmt.Sprintf("ERROR: %s", err.Error()))
	}

	err = json.Unmarshal([]byte(C.GoString(recipeInterface)), &recipe)
	if err != nil {
		return C.CString(fmt.Sprintf("ERROR: %s", err.Error()))
	}

	err = api.DownloadSource(recipe.DownloadsPath, module.Source, module.Name)
	if err != nil {
		return C.CString(fmt.Sprintf("ERROR: %s", err.Error()))
	}
	err = api.MoveSource(recipe.DownloadsPath, recipe.SourcesPath, module.Source, module.Name)
	if err != nil {
		return C.CString(fmt.Sprintf("ERROR: %s", err.Error()))
	}

	// Since the downloaded source goes through checksum verification already
	// it is safe to simply use the specified checksum from the module definition
	tmpDir := fmt.Sprintf("/tmp/%s-%s", module.Source.Checksum, module.Name)
	cmd := fmt.Sprintf(
		"cd /sources/%s && meson %s %s && ninja -C %s && ninja -C %s install",
		api.GetSourcePath(module.Source, module.Name),
		strings.Join(module.BuildFlags, " "),
		tmpDir,
		tmpDir,
		tmpDir,
	)

	return C.CString(cmd)
}

func main() { fmt.Println("This plugin is not meant to run standalone!") }
