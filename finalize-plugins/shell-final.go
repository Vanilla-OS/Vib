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

type Shell struct {
	Name     string   `json:"name"`
	Type     string   `json:"type"`
	Commands []string `json:"commands"`
	Cwd      string   `json:"cwd"`
}

//export PlugInfo
func PlugInfo() *C.char {
	plugininfo := &api.PluginInfo{Name: "shell-final", Type: api.FinalizePlugin}
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

func parsePath(path string, data *api.ScopeData) string {
	path = strings.ReplaceAll(path, "$PROJROOT", data.Recipe.ParentPath)
	path = strings.ReplaceAll(path, "$FSROOT", data.FS)
	return path
}

func baseCommand(command string, data *api.ScopeData) string {
	commandParts := strings.Split(command, " ")
	if strings.Contains(commandParts[0], "/") {
		return parsePath(commandParts[0], data)
	} else {
		command, err := exec.LookPath(commandParts[0])
		if err != nil {
			return commandParts[0]
		}
		return command
	}
}

func getArgs(command string, data *api.ScopeData) []string {
	commandParts := strings.Split(parsePath(command, data), " ")
	return commandParts[1:]
}

func genCommand(command string, data *api.ScopeData) []string {
	baseCommand := baseCommand(command, data)
	args := getArgs(command, data)
	return append(append(append([]string{"-c", "'"}, strings.Join(args, " ")), baseCommand), "'")
}

//export FinalizeBuild
func FinalizeBuild(moduleInterface *C.char, extraData *C.char) *C.char {
	var module *Shell
	var data *api.ScopeData

	err := json.Unmarshal([]byte(C.GoString(moduleInterface)), &module)
	if err != nil {
		return C.CString(fmt.Sprintf("ERROR: %s", err.Error()))
	}

	err = json.Unmarshal([]byte(C.GoString(extraData)), &data)
	if err != nil {
		return C.CString(fmt.Sprintf("ERROR: %s", err.Error()))
	}

	for _, command := range module.Commands {
		fmt.Println("shell-final:: bash ", "-c ", command)

		cmd := exec.Command(
			"bash", "-c", parsePath(command, data),
		)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Env = os.Environ()
		if len(strings.TrimSpace(module.Cwd)) == 0 {
			cmd.Dir = data.Recipe.ParentPath
		} else {
			cmd.Dir = parsePath(module.Cwd, data)
		}

		err = cmd.Run()
		if err != nil {
			return C.CString(fmt.Sprintf("ERROR: %s", err.Error()))
		}
	}

	return C.CString("")
}

func main() {}
