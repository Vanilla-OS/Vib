package main

import (
	"C"
	"fmt"

	"github.com/vanilla-os/vib/api"
)
import "encoding/json"

type GoModule struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	Source     api.Source
	BuildVars  map[string]string
	BuildFlags string
}

//export PlugInfo
func PlugInfo() *C.char {
	plugininfo := &api.PluginInfo{Name: "go", Type: api.BuildPlugin}
	pluginjson, err := json.Marshal(plugininfo)
	if err != nil {
		return C.CString(fmt.Sprintf("ERROR: %s", err.Error()))
	}
	return C.CString(string(pluginjson))
}

// BuildGoModule builds a module that builds a Go project
// buildVars are used to customize the build command
// like setting the output binary name and location
//
//export BuildModule
func BuildModule(moduleInterface *C.char, recipeInterface *C.char) *C.char {
	var module *GoModule
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
		"cd /sources/%s && go build%s -o %s",
		api.GetSourcePath(module.Source, module.Name),
		buildFlags,
		buildVars["GO_OUTPUT_BIN"],
	)

	return C.CString(cmd)
}

func main() { fmt.Println("This plugin is not meant to run standalone!") }
