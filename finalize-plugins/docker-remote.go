package main

import (
	"encoding/json"
	"fmt"

	"C"

	"github.com/vanilla-os/vib/api"
)

type DockerRemote struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	OCIRegistry string `json:"registry"`
	ImageName   string `json:"imagename"`
	ImageTag    string `json:"imagetag"`
}

//export PlugInfo
func PlugInfo() *C.char {
	plugininfo := &api.PluginInfo{Name: "docker-remote", Type: api.FinalizePlugin}
	pluginjson, err := json.Marshal(plugininfo)
	if err != nil {
		return C.CString(fmt.Sprintf("ERROR: %s", err.Error()))
	}
	return C.CString(string(pluginjson))
}

//export PluginScope
func PluginScope() int32 { // int32 is defined as GoInt32 in cgo which is the same as a C int
	return api.IMAGENAME | api.RUNTIME | api.IMAGEID
}

//export FinalizeBuild
func FinalizeBuild(moduleInterface *C.char, extraData *C.char) *C.char {
	var module *DockerRemote
	var data *api.ScopeData

	err := json.Unmarshal([]byte(C.GoString(moduleInterface)), &module)
	if err != nil {
		return C.CString(fmt.Sprintf("ERROR: %s", err.Error()))
	}

	err = json.Unmarshal([]byte(C.GoString(extraData)), &data)
	if err != nil {
		return C.CString(fmt.Sprintf("ERROR: %s", err.Error()))
	}

	fmt.Printf("Pushing %s with id %s\n", data.ImageName, data.ImageID)
	fmt.Printf("%s push %s %s\n", data.Runtime, data.ImageName, module.OCIRegistry)

	return C.CString("")
	//	return C.CString("ERROR: failed to upload image")
}

func main() {}
