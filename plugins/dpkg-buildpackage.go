package main

import (
	"C"
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/vanilla-os/vib/api"
)

// Configuration for building a Debian package using dpkg
type DpkgBuildModule struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Source api.Source
}

// Provide plugin information as a JSON string
//
//export PlugInfo
func PlugInfo() *C.char {
	plugininfo := &api.PluginInfo{Name: "dpkg-buildpackage", Type: api.BuildPlugin}
	pluginjson, err := json.Marshal(plugininfo)
	if err != nil {
		return C.CString(fmt.Sprintf("ERROR: %s", err.Error()))
	}
	return C.CString(string(pluginjson))
}

// Generate a command to build a Debian package using dpkg and install
// the resulting .deb package. Handle downloading, moving the source,
// and running dpkg-buildpackage with appropriate options.
//
//export BuildModule
func BuildModule(moduleInterface *C.char, recipeInterface *C.char) *C.char {
	var module *DpkgBuildModule
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

	cmd := fmt.Sprintf(
		"cd /sources/%s && dpkg-buildpackage -d -us -uc -b",
		filepath.Join(api.GetSourcePath(module.Source, module.Name)),
	)

	for _, path := range module.Source.Paths {
		cmd += fmt.Sprintf(" && apt install -y --allow-downgrades ../%s*.deb", path)
	}

	cmd += " && apt clean"
	return C.CString(cmd)
}

func main() { fmt.Println("This plugin is not meant to run standalone!") }
