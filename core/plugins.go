package core

import (
	"fmt"
	"github.com/vanilla-os/vib/api"
	"plugin"
)

var openedPlugins map[string]Plugin

func LoadPlugin(name string, module interface{}, recipe *api.Recipe) (string, error) {
    if openedPlugins == nil {
        openedPlugins = make(map[string]Plugin)
    }
	pluginOpened := false
	var buildModule Plugin
	buildModule, pluginOpened = openedPlugins[name]
	if !pluginOpened {
		fmt.Println("Loading new plugin")
		buildModule = Plugin{Name: name}
		var err error
		loadedPlugin, err := plugin.Open(fmt.Sprintf("%s/%s.so", recipe.PluginPath, name))
		if err != nil {
			panic(err)
		}
		buildFunction, err := loadedPlugin.Lookup("BuildModule")
		if err != nil {
			panic(err)
		}
		buildModule.BuildFunc = buildFunction.(func(interface{}, *api.Recipe) (string, error))
		buildModule.LoadedPlugin = loadedPlugin

		openedPlugins[name] = buildModule
	}
	fmt.Printf("Using plugin: %s\n", buildModule.Name)
	fmt.Println(buildModule.BuildFunc(module, recipe))
	return buildModule.BuildFunc(module, recipe)
}
