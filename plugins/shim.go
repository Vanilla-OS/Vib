package main

import (
	"encoding/json"
	"fmt"

	"C"

	"github.com/vanilla-os/vib/api"
)
import (
	"os"
	"os/exec"
	"path/filepath"
)

type ShimModule struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	ShimType string `json:"shimtype"`
}

//export PlugInfo
func PlugInfo() *C.char {
	plugininfo := &api.PluginInfo{Name: "shim", Type: api.BuildPlugin}
	pluginjson, err := json.Marshal(plugininfo)
	if err != nil {
		return C.CString(fmt.Sprintf("ERROR: %s", err.Error()))
	}
	return C.CString(string(pluginjson))
}

//export BuildModule
func BuildModule(moduleInterface *C.char, recipeInterface *C.char) *C.char {
	var module *ShimModule
	var recipe *api.Recipe

	err := json.Unmarshal([]byte(C.GoString(moduleInterface)), &module)
	if err != nil {
		return C.CString(fmt.Sprintf("ERROR: %s", err.Error()))
	}

	err = json.Unmarshal([]byte(C.GoString(recipeInterface)), &recipe)
	if err != nil {
		return C.CString(fmt.Sprintf("ERROR: %s", err.Error()))
	}

	fmt.Printf("[SHIM] Starting plugin: %s\n", module.ShimType)

	dataDir, err := os.MkdirTemp("", fmt.Sprintf("*-vibshim-%s", module.ShimType))
	if err != nil {
		return C.CString(fmt.Sprintf("ERROR: %s", err.Error()))
	}
	defer os.RemoveAll(dataDir)

	pluginCommand := fmt.Sprintf("%s/%s", recipe.PluginPath, module.ShimType)
	modulePath := filepath.Join(dataDir, "moduleInterface")
	recipePath := filepath.Join(dataDir, "recipeInterface")

	err = os.WriteFile(modulePath, []byte(C.GoString(moduleInterface)), 0o777)
	if err != nil {
		return C.CString(fmt.Sprintf("ERROR: %s", err.Error()))
	}
	err = os.WriteFile(recipePath, []byte(C.GoString(recipeInterface)), 0o777)
	if err != nil {
		return C.CString(fmt.Sprintf("ERROR: %s", err.Error()))
	}

	out, err := exec.Command(pluginCommand, modulePath, recipePath).Output()
	if err != nil {
		return C.CString(fmt.Sprintf("ERROR: %s", err.Error()))
	}
	return C.CString(string(out))
}

func main() {}
