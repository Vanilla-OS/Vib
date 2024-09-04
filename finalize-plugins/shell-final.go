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

// Configuration for a set of shell commands
type Shell struct {
	Name     string   `json:"name"`
	Type     string   `json:"type"`
	Commands []string `json:"commands"`
	Cwd      string   `json:"cwd"`
}

// Provide plugin information as a JSON string
//
//export PlugInfo
func PlugInfo() *C.char {
	plugininfo := &api.PluginInfo{Name: "shell-final", Type: api.FinalizePlugin}
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

// Check if the command is in $PATH or includes a directory path.
// Return the full path if found, otherwise return the command unchanged.
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

// Extract and return arguments from a command string
func getArgs(command string, data *api.ScopeData) []string {
	commandParts := strings.Split(parsePath(command, data), " ")
	return commandParts[1:]
}

// Generate an executable command by resolving the base command and arguments
// and wrapping them with appropriate syntax for execution.
func genCommand(command string, data *api.ScopeData) []string {
	baseCommand := baseCommand(command, data)
	args := getArgs(command, data)
	return append(append(append([]string{"-c", "'"}, strings.Join(args, " ")), baseCommand), "'")
}

// Execute shell commands from a Shell struct using the provided ScopeData.
// It parses and runs each command in the context of the provided working directory,
// or the recipe's parent path if no specific directory is given.
//
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
