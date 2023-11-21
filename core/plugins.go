package core

import (
	"fmt"
	"plugin"
)

var openedPlugins map[string]Plugin

func LoadPlugin(name string, module interface{}) (string, error) {
	pluginOpened := false
	var buildModule Plugin
	buildModule, pluginOpened = openedPlugins[name]
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
		buildModule.BuildFunc = buildFunction.(func(interface{}) (string, error))
		buildModule.LoadedPlugin = loadedPlugin

		openedPlugins[name] = buildModule
	}
	fmt.Printf("Using plugin: %s\n", buildModule.Name)
	fmt.Println(buildModule.BuildFunc(module))
	return buildModule.BuildFunc(module)
}
