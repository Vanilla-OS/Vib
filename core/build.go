package core

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/vanilla-os/vib/api"
)

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
func BuildRecipe(recipePath string, arch string) (api.Recipe, error) {
	// load the recipe
	recipe, err := LoadRecipe(recipePath)
	if err != nil {
		return api.Recipe{}, err
	}

	fmt.Printf("Building recipe %s\n", recipe.Name)

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

// Build commands for each module in the recipe
func BuildModules(recipe *api.Recipe, modules []interface{}, arch string, stageName string) ([]ModuleCommand, error) {
	cmds := []ModuleCommand{}
	for _, moduleInterface := range modules {
		var module Module
		err := mapstructure.Decode(moduleInterface, &module)
		if err != nil {
			return nil, err
		}

		cmd, err := BuildModule(recipe, moduleInterface, arch, stageName)
		if err != nil {
			return nil, err
		}

		cmds = append(cmds, ModuleCommand{
			Name:    module.Name,
			Command: append(cmd, ""), // add empty entry to ensure proper newline in Containerfile
			Workdir: module.Workdir,
		})
	}

	return cmds, nil
}

func buildIncludesModule(moduleInterface interface{}, recipe *api.Recipe, arch string, stageName string) (string, error) {
	var include IncludesModule
	err := mapstructure.Decode(moduleInterface, &include)
	if err != nil {
		return "", err
	}

	if len(include.Includes) == 0 {
		return "", errors.New("includes module must have at least one module to include")
	}

	var commands []string
	for _, include := range include.Includes {
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

		includeModule, err := GenModule(modulePath)
		if err != nil {
			return "", err
		}

		buildModule, err := BuildModule(recipe, includeModule, arch, stageName)
		if err != nil {
			return "", err
		}
		commands = append(commands, buildModule...)
	}
	return strings.Join(commands, "\n"), nil
}

// Build a command string for the given module in the recipe
func BuildModule(recipe *api.Recipe, moduleInterface interface{}, arch string, stageName string) ([]string, error) {
	var module Module
	err := mapstructure.Decode(moduleInterface, &module)
	if err != nil {
		return []string{""}, err
	}

	fmt.Printf("Building module [%s] of type [%s]\n", module.Name, module.Type)

	commands := []string{fmt.Sprintf("\n# Begin Module %s - %s", module.Name, module.Type)}

	if len(module.Modules) > 0 {
		for _, nestedModule := range module.Modules {
			buildModule, err := BuildModule(recipe, nestedModule, arch, stageName)
			if err != nil {
				return []string{""}, err
			}
			commands = append(commands, buildModule...)
		}
	}

	switch module.Type {
	case "shell":
		command, err := BuildShellModule(moduleInterface, recipe, arch)
		if err != nil {
			return []string{""}, err
		}
		commands = append(commands, command)
	case "includes":
		command, err := buildIncludesModule(moduleInterface, recipe, arch, stageName)
		if err != nil {
			return []string{""}, err
		}
		commands = append(commands, command)
	default:
		command, err := LoadBuildPlugin(module.Type, moduleInterface, recipe, arch)
		if err != nil {
			return []string{""}, err
		}
		commands = append(commands, command...)
	}

	sourcePath := filepath.Join(recipe.SourcesPath, module.Name)
	stageSourcePath := filepath.Join(recipe.SourcesPath, stageName, module.Name)
	_ = os.MkdirAll(sourcePath, 0o777)
	_ = os.MkdirAll(filepath.Dir(stageSourcePath), 0o777)
	err = os.Rename(sourcePath, stageSourcePath)
	if err != nil {
		return []string{}, fmt.Errorf("could not move source: %w", err)
	}

	commands = append(commands, fmt.Sprintf("# End Module %s - %s\n", module.Name, module.Type))

	fmt.Printf("Module [%s] built successfully\n", module.Name)
	return commands, nil
}
