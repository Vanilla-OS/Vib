package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"C"

	"github.com/vanilla-os/vib/api"
)
import (
	"strings"
)

// Configuration for an APT module
type AptModule struct {
	Name    string       `json:"name"`
	Type    string       `json:"type"`
	Options AptOptions   `json:"options"`
	Sources []api.Source `json:"sources"`
}

// Options for APT package management
type AptOptions struct {
	NoRecommends    bool `json:"no_recommends"`
	InstallSuggests bool `json:"install_suggests"`
	FixMissing      bool `json:"fix_missing"`
	FixBroken       bool `json:"fix_broken"`
}

// Provide plugin information as a JSON string
//
//export PlugInfo
func PlugInfo() *C.char {
	plugininfo := &api.PluginInfo{Name: "apt", Type: api.BuildPlugin, UseContainerCmds: false}
	pluginjson, err := json.Marshal(plugininfo)
	if err != nil {
		return C.CString(fmt.Sprintf("ERROR: %s", err.Error()))
	}
	return C.CString(string(pluginjson))
}

// Generate an apt-get install command from the provided module and recipe.
// Handle package installation and apply appropriate options.
//
//export BuildModule
func BuildModule(moduleInterface *C.char, recipeInterface *C.char) *C.char {
	var module *AptModule
	var recipe *api.Recipe

	err := json.Unmarshal([]byte(C.GoString(moduleInterface)), &module)
	if err != nil {
		return C.CString(fmt.Sprintf("ERROR: %s", err.Error()))
	}

	err = json.Unmarshal([]byte(C.GoString(recipeInterface)), &recipe)
	if err != nil {
		return C.CString(fmt.Sprintf("ERROR: %s", err.Error()))
	}

	args := ""
	if module.Options.NoRecommends {
		args += "--no-install-recommends "
	}
	if module.Options.InstallSuggests {
		args += "--install-suggests "
	}
	if module.Options.FixMissing {
		args += "--fix-missing "
	}
	if module.Options.FixBroken {
		args += "--fix-broken "
	}

	packages := ""
	for _, source := range module.Sources {
		if len(source.Packages) > 0 {
			for _, pkg := range source.Packages {
				packages += pkg + " "
			}
		}

		if len(strings.TrimSpace(source.Path)) > 0 {
			fileInfo, err := os.Stat(source.Path)
			if err != nil {
				return C.CString(fmt.Sprintf("ERROR: %s", err.Error()))
			}
			if !fileInfo.Mode().IsRegular() {
				continue
			}
			file, err := os.Open(source.Path)
			if err != nil {
				return C.CString(fmt.Sprintf("ERROR: %s", err.Error()))
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				packages += scanner.Text() + " "
			}

			if err := scanner.Err(); err != nil {
				return C.CString(fmt.Sprintf("ERROR: %s", err.Error()))
			}
		}
	}

	if len(packages) >= 1 {
		cmd := fmt.Sprintf("apt-get install -y %s %s && apt-get clean", args, packages)

		return C.CString(cmd)
	}

	return C.CString("ERROR: no packages or paths specified")
}

func main() {}
