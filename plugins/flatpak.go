package main

import (
	"C"
	"encoding/json"
	"fmt"

	"github.com/vanilla-os/vib/api"
)
import (
	"os"
	"path"
	"strings"
)

type innerFlatpakModule struct {
	Repourl  string   `json:"repo-url"`
	Reponame string   `json:"repo-name"`
	Install  []string `json:"install"`
	Remove   []string `json:"remove"`
}

type FlatpakModule struct {
	Name   string             `json:"name"`
	Type   string             `json:"type"`
	System innerFlatpakModule `json:"system"`
	User   innerFlatpakModule `json:"user"`
}

var SystemService string = `
[Unit]
Description=Manage system flatpaks
Wants=network-online.target
After=network-online.target

[Service]
Type=oneshot
ExecStart=/usr/bin/system-flatpak-setup
Restart=on-failure
RestartSec=30

[Install]
WantedBy=default.target
`

var UserService string = `
[Unit]
Description=Configure Flatpaks for current user
Wants=network-online.target
After=system-flatpak-setup.service

[Service]
Type=simple
ExecStart=/usr/bin/user-flatpak-setup
Restart=on-failure
RestartSec=30

[Install]
WantedBy=default.target
`

//export PlugInfo
func PlugInfo() *C.char {
	plugininfo := &api.PluginInfo{Name: "flatpak", Type: api.BuildPlugin}
	pluginjson, err := json.Marshal(plugininfo)
	if err != nil {
		return C.CString(fmt.Sprintf("ERROR: %s", err.Error()))
	}
	return C.CString(string(pluginjson))
}

func createRepo(module innerFlatpakModule, isSystem bool) string {
	fmt.Println("Adding remote ", isSystem, " ", module)
	command := "flatpak remote-add --if-not-exists"
	if isSystem {
		command = fmt.Sprintf("%s --system", command)
	} else {
		command = fmt.Sprintf("%s --user", command)
	}
	return fmt.Sprintf("%s %s %s", command, module.Reponame, module.Repourl)
}

//export BuildModule
func BuildModule(moduleInterface *C.char, recipeInterface *C.char) *C.char {
	var module *FlatpakModule
	var recipe *api.Recipe

	err := json.Unmarshal([]byte(C.GoString(moduleInterface)), &module)
	if err != nil {
		return C.CString(fmt.Sprintf("ERROR: %s", err.Error()))
	}

	err = json.Unmarshal([]byte(C.GoString(recipeInterface)), &recipe)
	if err != nil {
		return C.CString(fmt.Sprintf("ERROR: %s", err.Error()))
	}

	os.MkdirAll(path.Join(recipe.ParentPath, "includes.container/usr/bin/"), 0o775)
	if module.System.Reponame != "" {
		syscommands := "#!/usr/bin/env sh"
		if module.System.Repourl != "" {
			syscommands = fmt.Sprintf("%s\n%s", syscommands, createRepo(module.System, true))
			fmt.Println(syscommands)
		}
		if len(module.System.Install) != 0 {
			syscommands = fmt.Sprintf("%s\nflatpak install --system --noninteractive %s %s", syscommands, module.System.Reponame, strings.Join(module.System.Install, " "))
		}
		if len(module.System.Remove) != 0 {
			syscommands = fmt.Sprintf("%s\nflatpak uninstall --system --noninteractive %s %s", syscommands, module.User.Reponame, strings.Join(module.System.Remove, " "))
		}

		syscommands = fmt.Sprintf("%s\nsystemctl disable flatpak-system-setup.service", syscommands)
		err := os.WriteFile(path.Join(recipe.ParentPath, "includes.container/usr/bin/system-flatpak-setup"), []byte(syscommands), 0o777)
		if err != nil {
			return C.CString(fmt.Sprintf("ERROR: %s", err.Error()))
		}

	}
	if module.User.Reponame != "" {
		usercommands := "#!/usr/bin/env sh"
		if module.User.Repourl != "" {
			usercommands = fmt.Sprintf("%s\n%s", usercommands, createRepo(module.User, false))
			fmt.Println(usercommands)
		}
		if len(module.User.Install) != 0 {
			usercommands = fmt.Sprintf("%s\nflatpak install --user --noninteractive %s", usercommands, strings.Join(module.User.Install, " "))
		}
		if len(module.User.Remove) != 0 {
			usercommands = fmt.Sprintf("%s\nflatpak uninstall --user --noninteractive %s", usercommands, strings.Join(module.User.Remove, " "))
		}

		err := os.WriteFile(path.Join(recipe.ParentPath, "includes.container/usr/bin/user-flatpak-setup"), []byte(usercommands), 0o777)
		if err != nil {
			return C.CString(fmt.Sprintf("ERROR: %s", err.Error()))
		}
	}
	os.MkdirAll(path.Join(recipe.ParentPath, "includes.container/etc/systemd/user"), 0o775)
	os.MkdirAll(path.Join(recipe.ParentPath, "includes.container/etc/systemd/system"), 0o775)
	os.WriteFile(path.Join(recipe.ParentPath, "includes.container/etc/systemd/user/flatpak-user-setup.service"), []byte(UserService), 0o666)
	os.WriteFile(path.Join(recipe.ParentPath, "includes.container/etc/systemd/system/flatpak-system-setup.service"), []byte(SystemService), 0o666)

	return C.CString("systemctl enable --global flatpak-user-setup.service && systemctl enable --system flatpak-system-setup.service")
}

func main() { fmt.Println("This plugin is not meant to run standalone!") }
