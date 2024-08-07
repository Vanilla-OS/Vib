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

//export PluginScope
func PluginScope() int32 { // int32 is defined as GoInt32 in cgo which is the same as a C int
	return int32(api.IMAEGNAME)
}

//export FinalizeBuild
func FinalizeBuild(moduleInterface *C.char, extraData *C.char) *C.char {
	var module *DockerRemote

	err := json.Unmarshal([]byte(C.GoString(moduleInterface)), &module)
	if err != nil {
		return C.CString(fmt.Sprintf("ERROR: %s", err.Error()))
	}

	return C.CString("ERROR: failed to upload image")
}

func main() {}
