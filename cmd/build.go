package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vanilla-os/vib/core"
)

func NewBuildCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "build",
		Short: "Build the given recipe",
		Long:  "Build the given Vib recipe into a Containerfile",
		Example: `  Using the recipe.yml/yaml or vib.yml/yaml file in the current directory:
    vib build

  To specify a recipe file, use:
    vib build /path/to/recipe.yml`,
		RunE: buildCommand,
	}
	cmd.Flags().SetInterspersed(false)

	return cmd
}

func buildCommand(cmd *cobra.Command, args []string) error {
	commonNames := []string{
		"recipe.yml",
		"recipe.yaml",
		"vib.yml",
		"vib.yaml",
	}
	var recipePath string

	if len(args) == 0 {
		for _, name := range commonNames {
			if _, err := os.Stat(name); err == nil {
				recipePath = name
				break
			}
		}
	} else {
		recipePath = args[0]
	}

	if recipePath == "" {
		return fmt.Errorf("missing recipe path")
	}

	recipePath = args[0]
	_, err := core.BuildRecipe(recipePath)
	if err != nil {
		return err
	}

	return nil
}
