package core

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/vanilla-os/vib/api"
	"os"
)

// BuildRecipe builds a Containerfile from a recipe path
func BuildRecipe(recipePath string) (api.Recipe, error) {
	// load the recipe
	recipe, err := LoadRecipe(recipePath)
	if err != nil {
		return api.Recipe{}, err
	}

	fmt.Printf("Building recipe %s\n", recipe.Name)

	// resolve (and download) the sources
	//modules, sources, err := ResolveSources(recipe)
	//if err != nil {
	//	return Recipe{}, err
	//}

	// move them to the sources directory so they can be
	// used by the modules during the build
	//err = MoveSources(recipe, sources)
	//if err != nil {
	//	return Recipe{}, err
	//}

	// build the modules*
	// * actually just build the commands that will be used
	//   in the Containerfile to build the modules
	cmds, err := BuildModules(recipe, recipe.Modules)
	if err != nil {
		return api.Recipe{}, err
	}

	// build the Containerfile
	err = BuildContainerfile(recipe, cmds)
	if err != nil {
		return api.Recipe{}, err
	}

	return *recipe, nil
}

// BuildContainerfile builds a Containerfile from a recipe
// and a list of modules commands
func BuildContainerfile(recipe *api.Recipe, cmds []ModuleCommand) error {
	containerfile, err := os.Create(recipe.Containerfile)
	if err != nil {
		return err
	}

	defer containerfile.Close()

	// FROM
	_, err = containerfile.WriteString(
		fmt.Sprintf("FROM %s\n", recipe.Base),
	)
	if err != nil {
		return err
	}

	// LABELS
	for key, value := range recipe.Labels {
		_, err = containerfile.WriteString(
			fmt.Sprintf("LABEL %s='%s'\n", key, value),
		)
		if err != nil {
			return err
		}
	}

	// ARGS
	for key, value := range recipe.Args {
		_, err = containerfile.WriteString(
			fmt.Sprintf("ARG %s=%s\n", key, value),
		)
		if err != nil {
			return err
		}
	}

	// RUN(S)
	if !recipe.SingleLayer {
		for _, cmd := range recipe.Runs {
			_, err = containerfile.WriteString(
				fmt.Sprintf("RUN %s\n", cmd),
			)
			if err != nil {
				return err
			}
		}
	}
	// ADDS
	for key, value := range recipe.Adds {
		_, err = containerfile.WriteString(
			fmt.Sprintf("ADD %s %s\n", key, value),
		)
		if err != nil {
			return err
		}
	}

	// INCLUDES.CONTAINER
	_, err = containerfile.WriteString("ADD includes.container /\n")
	if err != nil {
		return err
	}

	// SOURCES
	_, err = containerfile.WriteString("ADD sources /sources\n")
	if err != nil {
		return err
	}

	// MODULES RUN(S)
	if !recipe.SingleLayer {
		for _, cmd := range cmds {
			_, err = containerfile.WriteString(
				fmt.Sprintf("RUN %s\n", cmd.Command),
			)
			if err != nil {
				return err
			}
		}
	}

	// SINGLE LAYER
	if recipe.SingleLayer {
		unifiedCmd := "RUN "

		for i, cmd := range recipe.Runs {
			unifiedCmd += cmd
			if i != len(recipe.Runs)-1 {
				unifiedCmd += " && "
			}
		}

		if len(cmds) > 0 {
			unifiedCmd += " && "
		}

		for i, cmd := range cmds {
			unifiedCmd += cmd.Command
			if i != len(cmds)-1 {
				unifiedCmd += " && "
			}
		}

		if len(unifiedCmd) > 4 {
			_, err = containerfile.WriteString(fmt.Sprintf("%s\n", unifiedCmd))
			if err != nil {
				return err
			}
		}
	}

	// CMD
	if recipe.Cmd != "" {
		_, err = containerfile.WriteString(
			fmt.Sprintf("CMD %s\n", recipe.Cmd),
		)
		if err != nil {
			return err
		}
	}

	return nil
}

// BuildModules builds a list of modules commands from a list of modules
func BuildModules(recipe *api.Recipe, modules []interface{}) ([]ModuleCommand, error) {
	cmds := []ModuleCommand{}
	for _, moduleInterface := range modules {
		var module Module
		err := mapstructure.Decode(moduleInterface, &module)
		if err != nil {
			return nil, err
		}
		fmt.Printf("Creating build command for %s\n", module)

		cmd, err := BuildModule(recipe, moduleInterface)
		if err != nil {
			return nil, err
		}

		cmds = append(cmds, ModuleCommand{
			Name:    module.Name,
			Command: cmd,
		})
	}

	return cmds, nil
}

// BuildModule builds a module command from a module
// this is done by calling the appropriate module builder
// function based on the module type
func BuildModule(recipe *api.Recipe, moduleInterface interface{}) (string, error) {
	var module Module
	err := mapstructure.Decode(moduleInterface, &module)
	if err != nil {
		return "", err
	}
	fmt.Printf("Processing module: %s\n", module.Type)
	fmt.Println(moduleInterface)
	switch module.Type {
	case "apt":
		return BuildAptModule(moduleInterface, recipe)
	case "cmake":
		return BuildCMakeModule(moduleInterface, recipe)
	case "dpkg":
		return BuildDpkgModule(moduleInterface, recipe)
	case "dpkg-buildpackage":
		return BuildDpkgBuildPkgModule(moduleInterface, recipe)
	case "go":
		return BuildGoModule(moduleInterface, recipe)
	case "make":
		return BuildMakeModule(moduleInterface, recipe)
	case "meson":
		return BuildMesonModule(moduleInterface, recipe)
	case "shell":
		return BuildShellModule(moduleInterface, recipe)
	case "includes":
		return "", nil
	default:
		return LoadPlugin(module.Type, moduleInterface, recipe)
	}
}
