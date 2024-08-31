package main

import (
	"C"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/vanilla-os/vib/api"
)

type Sysext struct {
	Name               string `json:"name"`
	Type               string `json:"type"`
	OSReleaseID        string `json:"osreleaseid"`
	OSReleaseVersionID string `json:"osreleaseversionid"`
}

//export PlugInfo
func PlugInfo() *C.char {
	plugininfo := &api.PluginInfo{Name: "sysext", Type: api.FinalizePlugin}
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
	var module *Sysext
	var data *api.ScopeData

	err := json.Unmarshal([]byte(C.GoString(moduleInterface)), &module)
	if err != nil {
		return C.CString(fmt.Sprintf("ERROR: %s", err.Error()))
	}

	err = json.Unmarshal([]byte(C.GoString(extraData)), &data)
	if err != nil {
		return C.CString(fmt.Sprintf("ERROR: %s", err.Error()))
	}

	var extensionRelease strings.Builder
	fmt.Fprintf(&extensionRelease, "ID=%s\n", module.OSReleaseID)
	fmt.Fprintf(&extensionRelease, "VERSION_ID=%s\n", module.OSReleaseVersionID)

	err = os.MkdirAll(filepath.Join(data.FS, "usr/lib/extension-release.d"), 0o777)
	if err != nil {
		return C.CString(fmt.Sprintf("ERROR: %s", err.Error()))
	}
	err = os.WriteFile(filepath.Join(data.FS, fmt.Sprintf("usr/lib/extension-release.d/extension-release.%s", data.Recipe.Id)), []byte(extensionRelease.String()), 0o777)
	if err != nil {
		return C.CString(fmt.Sprintf("ERROR: %s", err.Error()))
	}

	mksquashfs, err := exec.LookPath("mksquashfs")
	if err != nil {
		return C.CString(fmt.Sprintf("ERROR: %s", err.Error()))
	}
	cmd := exec.Command(
		mksquashfs, data.FS,
		filepath.Join(data.Recipe.ParentPath, fmt.Sprintf("%s.raw", data.Recipe.Id)),
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
