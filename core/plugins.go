package core

import (
	"fmt"
	"plugin"
	"strings"
)

var openedPlugins []Plugin

func LoadPlugin(name string, module []byte) {
	pluginOpened := false
	var buildModule Plugin
	for _, plugin := range openedPlugins {
		// To avoid loading the same plugin multiple times, we check the openedPlugins variable to see if the plugin
		// was loaded and added to the variable before
		if strings.ToLower(plugin.Name) == strings.ToLower(name) {
			pluginOpened = true
			buildModule = plugin
		}
	}
	if !pluginOpened {
		fmt.Println("Loading new plugin")
		buildModule = Plugin{Name: name}
		var err error
		loadedPlugin, err := plugin.Open(fmt.Sprintf("./plugins/%s.so", name)) // TODO: Proper path resolving
		if err != nil {
			panic(err)
		}
		buildFunction, err := loadedPlugin.Lookup("BuildModule")
		if err != nil {
			panic(err)
		}
		buildModule.BuildFunc = buildFunction.(func([]byte) (string, error))
		buildModule.LoadedPlugin = loadedPlugin

		openedPlugins = append(openedPlugins, buildModule)
	}
	fmt.Printf("Using plugin: %s\n", buildModule.Name)
	fmt.Println(buildModule.BuildFunc([]byte(":3")))
}
