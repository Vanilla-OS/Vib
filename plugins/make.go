package main

import (
	"C"
	"encoding/json"
	"fmt"

	"github.com/vanilla-os/vib/api"
)

type MakeModule struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Source api.Source
}

// BuildMakeModule builds a module that builds a Make project
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

	err = api.DownloadSource(recipe.DownloadsPath, module.Source, module.Name)
	if err != nil {
		return C.CString(fmt.Sprintf("ERROR: %s", err.Error()))
	}
	err = api.MoveSource(recipe.DownloadsPath, recipe.SourcesPath, module.Source, module.Name)
	if err != nil {
		return C.CString(fmt.Sprintf("ERROR: %s", err.Error()))
	}

	cmd := "cd /sources/" + api.GetSourcePath(module.Source, module.Name) + " && make && make install"
	return C.CString(cmd)
}

func main() { fmt.Println("This plugin is not meant to run standalone!") }
