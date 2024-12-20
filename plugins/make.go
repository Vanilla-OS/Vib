package main

import (
	"C"
	"encoding/json"
	"fmt"

	"github.com/vanilla-os/vib/api"
)
import "strings"

// Configuration for building a project using Make
type MakeModule struct {
	Name              string       `json:"name"`
	Type              string       `json:"type"`
	BuildCommand      string       `json:"buildcommand"`
	InstallCommand    string       `json:"installcommand"`
	IntermediateSteps []string     `json:"intermediatesteps"`
	Sources           []api.Source `json:"sources"`
}

// Provide plugin information as a JSON string
//
//export PlugInfo
func PlugInfo() *C.char {
	plugininfo := &api.PluginInfo{Name: "make", Type: api.BuildPlugin, UseContainerCmds: false}
	pluginjson, err := json.Marshal(plugininfo)
	if err != nil {
		return C.CString(fmt.Sprintf("ERROR: %s", err.Error()))
	}
	return C.CString(string(pluginjson))
}

// Generate a command to build a Make project. Change directory
// to the source path, run 'make' to build the project, and 'make install'
// to install the built project. Handle downloading and moving the source.
//
//export BuildModule
func BuildModule(moduleInterface *C.char, recipeInterface *C.char) *C.char {
	var module *MakeModule
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

	buildCommand := "make"
	installCommand := "make install"
	intermediateSteps := " && "

	if len(strings.TrimSpace(module.BuildCommand)) != 0 {
		buildCommand = module.BuildCommand
	}

	if len(strings.TrimSpace(module.InstallCommand)) != 0 {
		installCommand = module.InstallCommand
	}

	if len(module.IntermediateSteps) != 0 {
		intermediateSteps = " && " + strings.Join(module.IntermediateSteps, " && ") + " && "
	}

	cmd := "cd /sources/" + api.GetSourcePath(module.Sources[0], module.Name) + " && " + buildCommand + intermediateSteps + installCommand
	return C.CString(cmd)
}

func main() { fmt.Println("This plugin is not meant to run standalone!") }
