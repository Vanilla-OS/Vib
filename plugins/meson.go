package main

import (
	"C"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/vanilla-os/vib/api"
)
import "crypto/sha1"

// Configuration for building a Meson project
type MesonModule struct {
	Name       string
	Type       string
	BuildFlags []string     `json:"buildflags"`
	Sources    []api.Source `json"sources"`
}

// Provide plugin information as a JSON string
//
//export PlugInfo
func PlugInfo() *C.char {
	plugininfo := &api.PluginInfo{Name: "meson", Type: api.BuildPlugin, UseContainerCmds: false}
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

	for _, source := range module.Sources {
		err = api.DownloadSource(recipe, source, module.Name)
		if err != nil {
			return C.CString(fmt.Sprintf("ERROR: %s", err.Error()))
		}
		err = api.MoveSource(recipe.DownloadsPath, recipe.SourcesPath, source, module.Name)
		if err != nil {
			return C.CString(fmt.Sprintf("ERROR: %s", err.Error()))
		}
	}

	var tmpDir string
	if strings.EqualFold(module.Sources[0].Type, "git") == true {
		tmpDir = fmt.Sprintf("/tmp/%s-%s", module.Sources[0].Commit, module.Name)
	} else if module.Sources[0].Type == "tar" || module.Sources[0].Type == "local" {
		tmpDir = fmt.Sprintf("/tmp/%s-%s", module.Sources[0].Checksum, module.Name)
	} else {
		tmpDir = fmt.Sprintf("/tmp/%s-%s", sha1.Sum([]byte(module.Sources[0].URL)), module.Name)
	}
	cmd := fmt.Sprintf(
		"cd /sources/%s && meson %s %s && ninja -C %s && ninja -C %s install",
		api.GetSourcePath(module.Sources[0], module.Name),
		strings.Join(module.BuildFlags, " "),
		tmpDir,
		tmpDir,
		tmpDir,
	)

	return C.CString(cmd)
}

func main() { fmt.Println("This plugin is not meant to run standalone!") }
