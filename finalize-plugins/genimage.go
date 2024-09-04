package main

import (
	"C"
	"encoding/json"
	"fmt"
	"github.com/vanilla-os/vib/api"
	"os"
	"os/exec"
	"strings"
)

// Configuration for generating an image
type Genimage struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	GenimagePath string `json:"genimagepath"`
	Config       string `json:"config"`
	Rootpath     string `json:"rootpath"`
	Inputpath    string `json:"inputpath"`
	Outputpath   string `json:"outputpath"`
}

// Provide plugin information as a JSON string
//
//export PlugInfo
func PlugInfo() *C.char {
	plugininfo := &api.PluginInfo{Name: "genimage", Type: api.FinalizePlugin}
	pluginjson, err := json.Marshal(plugininfo)
	if err != nil {
		return C.CString(fmt.Sprintf("ERROR: %s", err.Error()))
	}
	return C.CString(string(pluginjson))
}

// Provide the plugin scope
//
//export PluginScope
func PluginScope() int32 { // int32 is defined as GoInt32 in cgo which is the same as a C int
	return api.IMAGENAME | api.FS | api.RECIPE
}

// Replace placeholders in the path with actual values from ScopeData
// $PROJROOT -> Recipe.ParentPath
// $FSROOT -> FS
func ParsePath(path string, data *api.ScopeData) string {
	path = strings.Replace(path, "$PROJROOT", data.Recipe.ParentPath, 1)
	path = strings.Replace(path, "$FSROOT", data.FS, 1)
	return path
}

// Complete the build process for a generated image module.
// Find the binary if not specified, replace path placeholders 
// in the module paths, and run the command
// with the provided configuration
//
//export FinalizeBuild
func FinalizeBuild(moduleInterface *C.char, extraData *C.char) *C.char {
	var module *Genimage
	var data *api.ScopeData

	err := json.Unmarshal([]byte(C.GoString(moduleInterface)), &module)
	if err != nil {
		return C.CString(fmt.Sprintf("ERROR: %s", err.Error()))
	}

	err = json.Unmarshal([]byte(C.GoString(extraData)), &data)
	if err != nil {
		return C.CString(fmt.Sprintf("ERROR: %s", err.Error()))
	}

	genimage := module.GenimagePath
	if genimage == "" {
		genimage, err = exec.LookPath("genimage")
		if err != nil {
			return C.CString(fmt.Sprintf("ERROR: %s", err.Error()))
		}
	}

	cmd := exec.Command(
		genimage,
		"--config",
		ParsePath(module.Config, data),
		"--rootpath",
		ParsePath(module.Rootpath, data),
		"--outputpath",
		ParsePath(module.Outputpath, data),
		"--inputpath",
		ParsePath(module.Inputpath, data),
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = data.Recipe.ParentPath

	err = cmd.Run()
	if err != nil {
		return C.CString(fmt.Sprintf("ERROR: %s", err.Error()))
	}

	return C.CString("")
}

func main() {}
