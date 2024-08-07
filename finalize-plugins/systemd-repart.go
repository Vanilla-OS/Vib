package main

import (
	"C"
	"encoding/json"
	"fmt"
	"github.com/vanilla-os/vib/api"
	"os"
	"os/exec"
)

type SystemdRepart struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Output string `json:"output"`
	Size   string `json:"size"`
	Seed   string `json:"seed"`
}

//export PlugInfo
func PlugInfo() *C.char {
	plugininfo := &api.PluginInfo{Name: "systemd-repart", Type: api.FinalizePlugin}
	pluginjson, err := json.Marshal(plugininfo)
	if err != nil {
		return C.CString(fmt.Sprintf("ERROR: %s", err.Error()))
	}
	return C.CString(string(pluginjson))
}

//export PluginScope
func PluginScope() int32 { // int32 is defined as GoInt32 in cgo which is the same as a C int
	return api.IMAGENAME | api.FS | api.RECIPE
}

//export FinalizeBuild
func FinalizeBuild(moduleInterface *C.char, extraData *C.char) *C.char {
	var module *SystemdRepart
	var data *api.ScopeData

	err := json.Unmarshal([]byte(C.GoString(moduleInterface)), &module)
	if err != nil {
		return C.CString(fmt.Sprintf("ERROR: %s", err.Error()))
	}

	err = json.Unmarshal([]byte(C.GoString(extraData)), &data)
	if err != nil {
		return C.CString(fmt.Sprintf("ERROR: %s", err.Error()))
	}

	repart, err := exec.LookPath("systemd-repart")
	if err != nil {
		return C.CString(fmt.Sprintf("ERROR: %s", err.Error()))
	}
	cmd := exec.Command(
		repart,
		"--definitions=definitions",
		"--empty=create",
		fmt.Sprintf("--size=%s", module.Size),
		"--dry-run=no",
		"--discard=no",
		"--offline=true",
		"--no-pager",
		fmt.Sprintf("--seed=%s", module.Seed),
		fmt.Sprintf("--root=%s", data.FS),
		module.Output,
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
