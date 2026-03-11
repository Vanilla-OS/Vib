package core

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/vanilla-os/vib/api"
)

var modulesCount int
var includeDepth int
var maxIncludeDepth = 1
var errorCount = 0

// Add a WORKDIR instruction to the containerfile
func ChangeWorkingDirectory(workdir string, containerfile *os.File) error {
	if workdir != "" {
		_, err := containerfile.WriteString(
			fmt.Sprintf("WORKDIR %s\n", workdir),
		)
		if err != nil {
			return err
		}
	}
	return nil
}

// Add a WORKDIR instruction to reset to the root directory
func RestoreWorkingDirectory(workdir string, containerfile *os.File) error {
	if workdir != "" {
		_, err := containerfile.WriteString(
			fmt.Sprintf("WORKDIR %s\n", "/"),
		)
		if err != nil {
			return err
		}
	}

	return nil
}

// Load and build a Containerfile from the specified recipe
func BuildRecipe(recipePath string, arch string, containerfilePath string) (api.Recipe, error) {
	// load the recipe
	recipe, err := LoadRecipe(recipePath)
	if err != nil {
		return api.Recipe{}, err
	}

	fmt.Printf("Building recipe `%s`\n", recipe.Name)

	// assuming the Containerfile location is relative
	if len(containerfilePath) == 0 {
		recipe.Containerfile = filepath.Join(filepath.Dir(recipePath), "Containerfile")
	} else {
		recipe.Containerfile = filepath.Join(filepath.Dir(recipePath), containerfilePath)
		fmt.Printf("Containerfile path: %s\n", recipe.Containerfile)
	}

	// build the Containerfile
	err = BuildContainerfile(recipe, arch)
	if err != nil {
		return api.Recipe{}, err
	}

	modules := 0
	for _, stage := range recipe.Stages {
		modules += len(stage.Modules)
	}

	fmt.Printf("Recipe %s built successfully\n", recipe.Name)
	fmt.Printf("Processed %d stages\n", len(recipe.Stages))
	fmt.Printf("Processed %d modules\n", modules)

	return *recipe, nil
}

// Generate a Containerfile from the recipe
func BuildContainerfile(recipe *api.Recipe, arch string) error {
	err := os.RemoveAll(recipe.Containerfile)
	if err != nil {
		return err
	}

	containerfile, err := os.Create(recipe.Containerfile)
	if err != nil {
		return err
	}

	defer containerfile.Close()

	for _, stage := range recipe.Stages {
		// build the modules*
		// * actually just build the commands that will be used
		//   in the Containerfile to build the modules
		cmds, err := BuildModules(recipe, stage.Modules, arch, stage.Id)
		if err != nil {
			return err
		}

		// FROM
		if stage.Id != "" {
			_, err = containerfile.WriteString(
				fmt.Sprintf("# Stage: %s\n", stage.Id),
			)
			if err != nil {
				return err
			}
			_, err = containerfile.WriteString(
				fmt.Sprintf("FROM %s AS %s\n", stage.Base, stage.Id),
			)
			if err != nil {
				return err
			}
		} else {
			_, err = containerfile.WriteString(
				fmt.Sprintf("FROM %s\n", stage.Base),
			)
			if err != nil {
				return err
			}
		}

		// COPY
		if len(stage.Copy) > 0 {
			for _, copy := range stage.Copy {
				if len(copy.SrcDst) > 0 {
					err = ChangeWorkingDirectory(copy.Workdir, containerfile)
					if err != nil {
						return err
					}

					for src, dst := range copy.SrcDst {
						if copy.From != "" {
							_, err = containerfile.WriteString(
								fmt.Sprintf("COPY --from=%s %s %s\n", copy.From, src, dst),
							)
							if err != nil {
								return err
							}
						} else {
							_, err = containerfile.WriteString(
								fmt.Sprintf("COPY %s %s\n", src, dst),
							)
							if err != nil {
								return err
							}
						}
					}

					err = RestoreWorkingDirectory(copy.Workdir, containerfile)
					if err != nil {
						return err
					}
				}
			}
		}

		// LABELS
		for key, value := range stage.Labels {
			_, err = containerfile.WriteString(
				fmt.Sprintf("LABEL %s='%s'\n", key, value),
			)
			if err != nil {
				return err
			}
		}

		// ENV
		for key, value := range stage.Env {
			_, err = containerfile.WriteString(
				fmt.Sprintf("ENV %s=%s\n", key, value),
			)
			if err != nil {
				return err
			}
		}

		// ARGS
		for key, value := range stage.Args {
			_, err = containerfile.WriteString(
				fmt.Sprintf("ARG %s=%s\n", key, value),
			)
			if err != nil {
				return err
			}
		}

		// RUN(S)
		if len(stage.Runs.Commands) > 0 {
			err = ChangeWorkingDirectory(stage.Runs.Workdir, containerfile)
			if err != nil {
				return err
			}

			for _, cmd := range stage.Runs.Commands {
				_, err = containerfile.WriteString(
					fmt.Sprintf("RUN %s\n", cmd),
				)
				if err != nil {
					return err
				}
			}

			err = RestoreWorkingDirectory(stage.Runs.Workdir, containerfile)
			if err != nil {
				return err
			}
		}

		// EXPOSE
		for key, value := range stage.Expose {
			_, err = containerfile.WriteString(
				fmt.Sprintf("EXPOSE %s/%s\n", key, value),
			)
			if err != nil {
				return err
			}
		}

		// ADDS
		if len(stage.Adds) > 0 {
			for _, add := range stage.Adds {
				if len(add.SrcDst) > 0 {
					err = ChangeWorkingDirectory(add.Workdir, containerfile)
					if err != nil {
						return err
					}

					for src, dst := range add.SrcDst {
						_, err = containerfile.WriteString(
							fmt.Sprintf("ADD %s %s\n", src, dst),
						)
						if err != nil {
							return err
						}
					}
				}

				err = RestoreWorkingDirectory(add.Workdir, containerfile)
				if err != nil {
					return err
				}
			}
		}

		// INCLUDES.CONTAINER
		if stage.Addincludes {
			_, err = containerfile.WriteString(fmt.Sprintf("ADD %s /\n", recipe.IncludesPath))
			if err != nil {
				return err
			}
		}

		// SOURCES
		sourcePath := filepath.Join("sources", stage.Id)
		err = os.MkdirAll(sourcePath, 0o755)
		if err != nil {
			return fmt.Errorf("could not create source path: %w", err)
		}
		_, err = containerfile.WriteString(fmt.Sprintf("ADD %s /sources\n", sourcePath))
		if err != nil {
			return err
		}

		for _, cmd := range cmds {
			err = ChangeWorkingDirectory(cmd.Workdir, containerfile)
			if err != nil {
				return err
			}

			_, err = containerfile.WriteString(strings.Join(cmd.Command, "\n"))
			if err != nil {
				return err
			}

			err = RestoreWorkingDirectory(cmd.Workdir, containerfile)
			if err != nil {
				return err
			}
		}

		// CMD
		err = ChangeWorkingDirectory(stage.Cmd.Workdir, containerfile)
		if err != nil {
			return err
		}

		if len(stage.Cmd.Exec) > 0 {
			_, err = containerfile.WriteString(
				fmt.Sprintf("CMD [\"%s\"]\n", strings.Join(stage.Cmd.Exec, "\",\"")),
			)
			if err != nil {
				return err
			}

			err = RestoreWorkingDirectory(stage.Cmd.Workdir, containerfile)
			if err != nil {
				return err
			}
		}

		// DELETE SOURCES
		_, err = containerfile.WriteString("RUN rm -r /sources\n")
		if err != nil {
			return err
		}

		// ENTRYPOINT
		err = ChangeWorkingDirectory(stage.Entrypoint.Workdir, containerfile)
		if err != nil {
			return err
		}

		if len(stage.Entrypoint.Exec) > 0 {
			_, err = containerfile.WriteString(
				fmt.Sprintf("ENTRYPOINT [\"%s\"]\n", strings.Join(stage.Entrypoint.Exec, "\",\"")),
			)
			if err != nil {
				return err
			}

			err = RestoreWorkingDirectory(stage.Entrypoint.Workdir, containerfile)
			if err != nil {
				return err
			}
		}

		containerfile.WriteString("\n")
	}

	return nil
}

func ExhaustCollectedErrors(_errors *[]error) int {
	length := len(*_errors)
	if length == 0 {
		return 0
	}

	for _, err := range *_errors {
		fmt.Printf("%v\n", err)
	}

	*_errors = nil
	errorCount += length
	return length
}

func MapSlicesToInterfaceSlices(inter []map[string]interface{}) []interface{} {
	result := make([]interface{}, len(inter))
	for i, m := range inter {
		result[i] = m
	}
	return result
}

func MapToInterfaceSlices(m map[string]interface{}) []interface{} {
	panic("If you need to use this function, you're likely handling module interfaces wrong")
	// result := make([]interface{}, 0, len(m))
	// for _, v := range m {
	// 	result = append(result, v)
	// }
	// return result
}

func DecodeModuleToGenericModule(module interface{}, errors *[]error) (Module, error) {
	var decodedModule Module
	var customErr error = nil
	defaultErr := mapstructure.Decode(module, &decodedModule)

	if defaultErr != nil {
		customErr = fmt.Errorf("error: yaml decode error: failed to decode module to generic module with further error: %v", defaultErr)
		(*errors) = append((*errors), customErr)
	}

	return decodedModule, customErr
}

func CollectModulesRecursively(modules []interface{}, allModules *[]interface{}, occurances *map[string][]int, errors *[]error) {
	for _, module := range modules {
		modulesCount++

		decodedModule, err := DecodeModuleToGenericModule(module, errors)
		if err != nil {
			continue
		}

		c := 0
		if decodedModule.Name == "" {
			c |= 0b10
		}
		if decodedModule.Type == "" {
			c |= 0b01
		}

		switch c {
		case 0b11:
			fmt.Printf("error: module name and type cannot be == \"\"")
			continue
		case 0b10:
			fmt.Printf("error: module name cannot be == \"\"")
			continue
		case 0b01:
			fmt.Printf("error: module type cannot be == \"\"")
			continue
		case 0b00:
			// fallthrough
		}

		*allModules = append(*allModules, module)
		(*occurances)[decodedModule.Name] = append((*occurances)[decodedModule.Name], modulesCount)

		ExhaustCollectedErrors(errors)

		if len(decodedModule.Modules) > 0 {
			CollectModulesRecursively(MapSlicesToInterfaceSlices(decodedModule.Modules), allModules, occurances, errors)
		}
	}
}

// Build commands for each module in the recipe
func BuildModules(recipe *api.Recipe, modules []interface{}, arch string, stageName string) ([]ModuleCommand, error) {
	var _errors []error
	var allModules []interface{}
	modNameOccursInMod := make(map[string][]int)

	CollectModulesRecursively(modules, &allModules, &modNameOccursInMod, &_errors)

	cmds := []ModuleCommand{}

	for _, moduleInterface := range modules {
		decodedModule, cmd, err := BuildModule(recipe, moduleInterface, &allModules, &modNameOccursInMod, arch, stageName, &_errors)
		if err != nil {
			if !(decodedModule.Type == "includes" && (errors.Is(err, os.ErrNotExist) || errors.Is(err, fs.ErrNotExist))) {
				_errors = append(_errors, err)
			}
			ExhaustCollectedErrors(&_errors)
			fmt.Printf("Building [%s] module `%s`: failed\n", decodedModule.Type, decodedModule.Name)

			continue
		}

		ExhaustCollectedErrors(&_errors)

		cmds = append(cmds, ModuleCommand{
			Name:    decodedModule.Name,
			Command: append(cmd, ""), // add empty entry to ensure proper newline in Containerfile
			Workdir: decodedModule.Workdir,
		})
	}

	for _, occurancesInMods := range modNameOccursInMod {
		occurances := len(occurancesInMods)

		if occurances > 1 {
			decodedModule, err := DecodeModuleToGenericModule(allModules[occurancesInMods[0]], &_errors)
			if err != nil {
				panic("This module was previously decoded but now fails to. Needs fix in vib codebase.")
			}
			_errors = append(_errors, fmt.Errorf("error: found ambiguous module with name `%s` %d times:", decodedModule.Name, occurances))

			for j := range occurances {
				if j > 0 {
					decodedModule, err = DecodeModuleToGenericModule(allModules[occurancesInMods[j]], &_errors)
				} else if err != nil {
					continue
				}
				_errors = append(_errors, fmt.Errorf("note:                     found in file `%s`", decodedModule.Workdir))
				// TODO: This is not the correct variable to get the file path of the module, which we should display.
			}
			ExhaustCollectedErrors(&_errors)
			errorCount += occurances
		}
	}

	if errorCount > 0 {
		return nil, fmt.Errorf("Encoutered %d errors while building %d modules\n", errorCount, modulesCount)
	}

	return cmds, nil
}

func BuildIncludesModule(recipe *api.Recipe, module interface{}, allModules *[]interface{}, occurances *map[string][]int, arch string, stageName string, _errors *[]error) (string, error) {
	// Note: errors is called _errors here because this function needs the errors package.

	includeDepth++
	defer func() { includeDepth-- }()

	var includeModule IncludesModule
	if _err := mapstructure.Decode(module, &includeModule); _err != nil {
		return "", _err
	}

	if includeDepth > 1 {
		return "", fmt.Errorf("[includes] module nesting is currently limited to `%d` layers.\n       Found includes module in `%s`\n", maxIncludeDepth, includeModule.Name)
	}

	if len(includeModule.Includes) == 0 {
		return "", fmt.Errorf("[includes] module `%s` must have at least one module to include", includeModule.Name)
	}

	var commands []string
	var err error = nil
	for _, include := range includeModule.Includes {
		var modulePath string

		// in case of a remote include, we need to download the
		// recipe before including it
		if include[:4] == "http" {
			fmt.Printf("Downloading recipe from %s\n", include)
			modulePath, err = downloadRecipe(include)
			if err != nil {
				return "", err
			}
		} else if followsGhPattern(include) {
			// if the include follows the github pattern, we need to
			// download the recipe from the github repository
			fmt.Printf("Downloading recipe from %s\n", include)
			modulePath, err = downloadGhRecipe(include)
			if err != nil {
				return "", err
			}
		} else {
			modulePath = filepath.Join(recipe.ParentPath, include)
		}

		generatedModule, _err := GenModule(modulePath)

		if errors.Is(_err, os.ErrNotExist) || errors.Is(_err, fs.ErrNotExist) {
			customErr := fmt.Errorf("error: [%s] module `%s` includes\n       `%s`,\n       which doesn't exist", includeModule.Type, includeModule.Name, modulePath)

			(*_errors) = append(*_errors, customErr)
			(*_errors) = append(*_errors, _err)

			err = _err
			continue
		} else if _err != nil {
			(*_errors) = append(*_errors, _err)
			return "", _err
		}

		ExhaustCollectedErrors(_errors)

		var _errors []error // temporary

		decodedModule, cmd, _err := BuildModule(recipe, generatedModule, allModules, occurances, arch, stageName, &_errors)
		if _err != nil {
			// _errors = append(_errors, _err)
			ExhaustCollectedErrors(&_errors)
			err = _err
			continue
		}

		commands = append(commands, cmd...)
		*allModules = append(*allModules, decodedModule)
		(*occurances)[decodedModule.Name] = append((*occurances)[decodedModule.Name], modulesCount)

		fmt.Printf("Building all %d submodules of [%s] module `%s` included in `%s`\n", len(decodedModule.Modules), decodedModule.Type, decodedModule.Name, includeModule.Name)
		var failed bool = false
		includeModuleIdx := len(*allModules)
		CollectModulesRecursively(MapSlicesToInterfaceSlices(decodedModule.Modules), allModules, occurances, &_errors)

		modulesLeftToBuild := len(*allModules) - includeModuleIdx

		for i := modulesLeftToBuild; i > 0; i-- {
			_, buildModule, _err := BuildModule(recipe, (*allModules)[len(*allModules)-i], allModules, occurances, arch, stageName, &_errors)
			if _err != nil {
				ExhaustCollectedErrors(&_errors)
				fmt.Printf("%d/%d Building [%s] module of submodule `%s` included in `%s`: failed\n", i, len(decodedModule.Modules), decodedModule.Type, decodedModule.Name, includeModule.Name)
				failed = true
				err = _err
				continue
			}

			commands = append(commands, buildModule...)
		}

		ExhaustCollectedErrors(&_errors)
		if failed {
			fmt.Printf("Building all %d submodules of [%s] module `%s` included in `%s`: failed\n", len(decodedModule.Modules), decodedModule.Type, decodedModule.Name, includeModule.Name)
		} else {
			fmt.Printf("Buildung all %d submodules of [%s] module `%s` included in `%s`: success\n", len(decodedModule.Modules), decodedModule.Type, decodedModule.Name, includeModule.Name)
		}
	}
	return strings.Join(commands, "\n"), err
}

// Build a command string for the given module in the recipe
func BuildModule(recipe *api.Recipe, module interface{}, allModules *[]interface{}, occurances *map[string][]int, arch string, stageName string, _errors *[]error) (Module, []string, error) {
	decodedModule, err := DecodeModuleToGenericModule(module, _errors)
	if err != nil {
		return decodedModule, []string{""}, err
	}

	commands := []string{fmt.Sprintf("\n# Begin Module %s - %s", decodedModule.Name, decodedModule.Type)}
	defer func() {
		commands = append(commands, fmt.Sprintf("# End Module %s - %s\n", decodedModule.Name, decodedModule.Type))
	}()

	fmt.Printf("Building [%s] module `%s`\n", decodedModule.Type, decodedModule.Name)

	switch decodedModule.Type {
	case "shell":
		command, err := BuildShellModule(module, recipe, arch)
		if err != nil {
			return decodedModule, []string{""}, err
		}
		commands = append(commands, command)
	case "includes":
		command, err := BuildIncludesModule(recipe, module, allModules, occurances, arch, stageName, _errors)
		if err != nil {
			return decodedModule, []string{""}, err
		}
		commands = append(commands, command)
	case "":
		err := fmt.Errorf("error: module `%s` tried to use a plugin but specified no name", decodedModule.Name)
		return decodedModule, []string{""}, err
	default:
		command, err := LoadBuildPlugin(decodedModule.Type, module, recipe, arch)
		if err != nil {
			return decodedModule, []string{""}, err
		}
		commands = append(commands, command...)
	}

	sourcePath := filepath.Join(recipe.SourcesPath, decodedModule.Name)
	stageSourcePath := filepath.Join(recipe.SourcesPath, stageName, decodedModule.Name)

	_ = os.MkdirAll(sourcePath, 0o777)
	_ = os.MkdirAll(filepath.Dir(stageSourcePath), 0o777)

	err = os.Rename(sourcePath, stageSourcePath)
	if err != nil {
		if errors.Is(err, os.ErrExist) || errors.Is(err, fs.ErrExist) {
			fmt.Printf("Multiple module name error!\n")
			return decodedModule, []string{}, nil
		}
		return decodedModule, []string{}, fmt.Errorf("could not rename `%s` to `%s`: %w\n", sourcePath, stageSourcePath, err)
	}

	fmt.Printf("Building [%s] module `%s`: success\n", decodedModule.Type, decodedModule.Name)
	return decodedModule, commands, nil
}
