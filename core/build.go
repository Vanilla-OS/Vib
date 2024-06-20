package core

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/mitchellh/mapstructure"

	"github.com/vanilla-os/vib/api"
)

// BuildRecipe builds a Containerfile from a recipe path
func BuildRecipe(recipePath string) (api.Recipe, error) {
	// load the recipe
	recipe, err := LoadRecipe(recipePath)
	if err != nil {
		return api.Recipe{}, err
	}

	fmt.Printf("Building recipe %s\n", recipe.Name)

	// build the Containerfile
	err = BuildContainerfile(recipe)
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

// BuildContainerfile builds a Containerfile from a recipe
// and a list of modules commands
func BuildContainerfile(recipe *api.Recipe) error {
	containerfile, err := os.Create(recipe.Containerfile)
	if err != nil {
		return err
	}

	defer containerfile.Close()

	for _, stage := range recipe.Stages {
		// build the modules*
		// * actually just build the commands that will be used
		//   in the Containerfile to build the modules
		cmds, err := BuildModules(recipe, stage.Modules)
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
				if len(copy.Paths) > 0 {
					if copy.Workdir != "" && copy.Workdir != recipe.Cwd {
						_, err = containerfile.WriteString(
							fmt.Sprintf("WORKDIR %s\n", copy.Workdir),
						)
						recipe.Cwd = copy.Workdir
						if err != nil {
							return err
						}
					}
					for _, path := range copy.Paths {
						if copy.From != "" {
							_, err = containerfile.WriteString(
								fmt.Sprintf("COPY --from=%s %s %s\n", copy.From, path.Src, path.Dst),
							)
							if err != nil {
								return err
							}
						} else {
							_, err = containerfile.WriteString(
								fmt.Sprintf("COPY %s %s\n", path.Src, path.Dst),
							)
							if err != nil {
								return err
							}
						}
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
		if !stage.SingleLayer {
			if len(stage.Runs.Shell) > 0 {
				if stage.Runs.Workdir != "" && stage.Runs.Workdir != recipe.Cwd {
					_, err = containerfile.WriteString(
						fmt.Sprintf("WORKDIR %s\n", stage.Runs.Workdir),
					)
					recipe.Cwd = stage.Runs.Workdir
					if err != nil {
						return err
					}
				}
				for _, cmd := range stage.Runs.Shell {
					_, err = containerfile.WriteString(
						fmt.Sprintf("RUN %s\n", cmd),
					)
					if err != nil {
						return err
					}
				}
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
		for key, value := range stage.Adds {
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
		if !stage.SingleLayer {
			for _, cmd := range cmds {
				if cmd.Command == "" {
					continue
				}

				if cmd.Workdir != "" && cmd.Workdir != recipe.Cwd {
					_, err = containerfile.WriteString(
						fmt.Sprintf("WORKDIR %s\n", cmd.Workdir),
					)
					recipe.Cwd = cmd.Workdir
					if err != nil {
						return err
					}
				}

				_, err = containerfile.WriteString(
					fmt.Sprintf("RUN %s\n", cmd.Command),
				)
				if err != nil {
					return err
				}
			}
		}

		// SINGLE LAYER
		if stage.SingleLayer {
			if len(stage.Runs.Shell) > 0 {
				if stage.Runs.Workdir != "" && stage.Runs.Workdir != recipe.Cwd {
					_, err = containerfile.WriteString(
						fmt.Sprintf("WORKDIR %s\n", stage.Runs.Workdir),
					)
					recipe.Cwd = stage.Runs.Workdir
					if err != nil {
						return err
					}
				}

				unifiedCmd := "RUN "

				for i, cmd := range stage.Runs.Shell {
					unifiedCmd += cmd
					if i != len(stage.Runs.Shell)-1 {
						unifiedCmd += " && "
					}
				}

				if len(cmds) > 0 {
					unifiedCmd += " && "
				}

				for i, cmd := range cmds {
					if cmd.Workdir != stage.Runs.Workdir {
						return errors.New("Workdir mismatch")
					}
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
		}

		// CMD
		if len(stage.Cmd) > 0 {
			_, err = containerfile.WriteString(
				fmt.Sprintf("CMD [\"%s\"]\n", strings.Join(stage.Cmd, "\",\"")),
			)
			if err != nil {
				return err
			}
		}

		// ENTRYPOINT
		if stage.Entrypoint.Workdir != "" && stage.Entrypoint.Workdir != recipe.Cwd {
			_, err = containerfile.WriteString(
				fmt.Sprintf("WORKDIR %s\n", stage.Entrypoint.Workdir),
			)
			recipe.Cwd = stage.Entrypoint.Workdir
			if err != nil {
				return err
			}
		}
		if len(stage.Entrypoint.Exec) > 0 {
			_, err = containerfile.WriteString(
				fmt.Sprintf("ENTRYPOINT [\"%s\"]\n", strings.Join(stage.Entrypoint.Exec, "\",\"")),
			)
			if err != nil {
				return err
			}
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

		cmd, err := BuildModule(recipe, moduleInterface)
		if err != nil {
			return nil, err
		}

		cmds = append(cmds, ModuleCommand{
			Name:    module.Name,
			Command: cmd,
			Workdir: module.Workdir,
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

	fmt.Printf("Building module [%s] of type [%s]\n", module.Name, module.Type)

	var commands string
	if len(module.Modules) > 0 {
		for _, nestedModule := range module.Modules {
			buildModule, err := BuildModule(recipe, nestedModule)
			if err != nil {
				return "", err
			}
			commands = commands + " && " + buildModule
		}
	}

	moduleBuilders := map[string]func(interface{}, *api.Recipe) (string, error){
		"shell":    BuildShellModule,
		"includes": func(interface{}, *api.Recipe) (string, error) { return "", nil },
	}

	if moduleBuilder, ok := moduleBuilders[module.Type]; ok {
		command, err := moduleBuilder(moduleInterface, recipe)
		if err != nil {
			return "", err
		}
		commands = commands + " && " + command
	} else {
		command, err := LoadPlugin(module.Type, moduleInterface, recipe)
		if err != nil {
			return "", err
		}
		commands = commands + " && " + command
	}

	fmt.Printf("Module [%s] built successfully\n", module.Name)
	result := strings.TrimPrefix(commands, " && ")

	if result == "&&" {
		return "", nil
	}

	return result, nil
}
