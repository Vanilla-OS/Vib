---
Title: Making a Build Plugin
Description: How to create a custom build plugin for Vib.
PublicationDate: 2024-02-14
Listed: true
Authors:
  - mirkobrombin
  - axtloss
Tags:
  - github
  - build
---

Vib supports custom, user-made plugins to expand the functionality of modules.

## Types of Plugin 

There are two types of plugins currently supported:

- Build Plugins
- Finalize Plugins

Build Plugins can be used in modules, they create commands that are part of the Containerfile. 

Finalize Plugins are run after the image has been built, they allow things such as generating isos from the filesystem or directly uploading the built image to a registry.

This article will focus on Build Plugins, Finalize Plugins are built differently, and will get their own documentation soon. 

## Plugin Requirements

Plugins are built as shared object files. In the `go` language, this can be achieved using the `-buildmode=c-shared` flag, while `gcc` requires the `-shared` flag. This way it is possible to write plugins in any compiled language capable of generating shared object files. Plugins can also be created using other languages, though additional steps may be necessary (details on that later).

The primary communication between `vib` and the plugin is handled through structs serialized to JSON.

Each build plugin must implement the following functions:

| Function Name | Arguments | Return Type | Description |
|---------------|-----------|-------------|-------------|
| `PlugInfo` |  | `char*` | Returns information about the plugin, typically as a JSON string. |
| `BuildModule` | `moduleInterface *char`, `recipeInterface *char` | `char*` | The main entry point for the plugin. Called by `vib` to retrieve the command to be executed. The command is returned as a JSON string. |

### char* PlugInfo()

This function returns information about the plugin, most notably the type of plugin.

Plugins that do not define this function are considered deprecated, while they still work, support may be dropped in future releases.

The function returns the `api.PluginInfo` struct serialised as a json:

```json
{
	"name": "<plugin name>",
	"type": 0
}
```

Vib gets the plugin type from the `type` field: `0` means `BuildPlugin`, and `1` means `FinalizePlugin`. For this article, it should be set to `0`, as it does not cover the requirements for a finalize plugin.

example function:

```C
char* PlugInfo() {
	return "{\"name\":\"example\",\"type\":0}";
}
```

### char* BuildModule(char* moduleInterface, char\* recipeInterface)

This is the entry point for plugins that vib calls. It returns a string prefixed with `ERROR:` if an error occurs, otherwise it returns the commands generated for the module.

The `moduleInterface` argument is a json serialised version of the module defined in the recipe.

The `recipeInterface` argument is a json serialised version of the entire recipe.

example function:

```C
char* BuildModule(char* moduleInterface, char* recipeInterface) {
	return "echo HAII";
}
```

## Making plugins without compiling to so files

One of the vib plugins is the `shim` plugin, it allows users to use plugins in any scripting languages, or regular executables.

The plugin writes the moduleInterface and recipeInterface into a temporary directory, the paths are given as arguments to the executable.

Shim then reads the generated commands from stdout.

example shim plugin:

```bash
#!/usr/bin/bash
username=$(cat $1 | jq .username)
echo "useradd -m ${username} && echo '${username}' | passwd ${username} --stdin"
```


## Plugin examples

We provide a plugin template for plugins written in go in the [vib-plugin repo](https://github.com/Vanilla-OS/vib-plugin).

Example plugins written in other languages than go can be found in axtlos' [vib-plugins repo](https://github.com/axtloss/vib-plugins/)

