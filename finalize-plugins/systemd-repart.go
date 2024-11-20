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

// Configuration for systemd repartitioning
type SystemdRepart struct {
	Name            string   `json:"name"`
	Type            string   `json:"type"`
	Output          string   `json:"output"`
	Json            string   `json:"json"`
	SpecOutput      string   `json:"spec_output"`
	Size            string   `json:"size"`
	Seed            string   `json:"seed"`
	Split           bool     `json:"split"`
	Empty           string   `json:"empty"`
	Root            string   `json:"root"`
	DeferPartitions []string `json:"defer_partitions"`
}

// Provide plugin information as a JSON string
//
//export PlugInfo
func PlugInfo() *C.char {
	plugininfo := &api.PluginInfo{Name: "systemd-repart", Type: api.FinalizePlugin}
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
func parsePath(path string, data *api.ScopeData) string {
	path = strings.ReplaceAll(path, "$PROJROOT", data.Recipe.ParentPath)
	path = strings.ReplaceAll(path, "$FSROOT", data.FS)
	return path
}

// Finalize the build by executing systemd-repart with the provided configuration
// to generate and apply partitioning specifications and output results
//
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

	if len(strings.TrimSpace(module.Json)) == 0 {
		module.Json = "off"
	}

	if len(strings.TrimSpace(module.Empty)) == 0 {
		module.Empty = "create"
	}

	if len(strings.TrimSpace(module.Root)) == 0 {
		module.Root = data.FS
	} else {
		module.Root = parsePath(module.Root, data)
	}

	args := []string{
		"--definitions=definitions",
		fmt.Sprintf("--empty=%s", module.Empty),
		fmt.Sprintf("--size=%s", module.Size),
		"--dry-run=no",
		"--discard=no",
		"--offline=true",
		"--no-pager",
		fmt.Sprintf("--split=%t", module.Split),
		fmt.Sprintf("--seed=%s", module.Seed),
		fmt.Sprintf("--root=%s", data.FS),
		module.Output,
		fmt.Sprintf("--json=%s", module.Json),
	}

	if len(module.DeferPartitions) > 0 {
		args = append(args, fmt.Sprintf("--defer-partitions=%s", strings.Join(module.DeferPartitions, ",")))
	}

	cmd := exec.Command(
		repart,
		args...,
	)
	jsonFile, err := os.Create(module.SpecOutput)
	if err != nil {
		return C.CString(fmt.Sprintf("ERROR: %s", err.Error()))
	}
	defer jsonFile.Close()
	cmd.Stdout = jsonFile
	cmd.Stderr = os.Stderr
	cmd.Dir = data.Recipe.ParentPath

	err = cmd.Run()
	if err != nil {
		return C.CString(fmt.Sprintf("ERROR: %s", err.Error()))
	}

	return C.CString("")
}

func main() {}
