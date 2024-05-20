package main

import (
	"C"
	"fmt"
	"path/filepath"
	"encoding/json"

	"github.com/vanilla-os/vib/api"
)

type DpkgBuildModule struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Source api.Source
}

// BuildDpkgModule builds a module that builds a dpkg project
// and installs the resulting .deb package
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



func main() { fmt.Println("This plugin is not meant to run standalone!"); }
