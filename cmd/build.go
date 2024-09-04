package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vanilla-os/vib/core"
)

// Create a new build command for the Cobra CLI
//
// Returns: new Cobra command for building a recipe
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

// Handle the build command for the Cobra CLI
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

		/*
			Check whether the provided file has either yml or yaml extension,
			if not, then return an error

			Operations on recipePath:
			1. Get the recipePath extension, then
			2. Trim the left dot(.) and
			3. Convert the extension to lower case.

			Covers the following:
			1. filename.txt - Invalid extension
			2. filename. - No extension
			3. filename - No extension
			4. filename.YAML or filename.YML - uppercase extension
		*/
		extension := strings.ToLower(strings.TrimLeft(filepath.Ext(recipePath), "."))
		if len(extension) == 0 || (extension != "yml" && extension != "yaml") {
			return fmt.Errorf("%s is an invalid recipe file", recipePath)
		}

		// Check whether the provided file exists, if not, then return an error
		if _, err := os.Stat(recipePath); errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("%s does not exist", recipePath)
		}
	}

	if recipePath == "" {
		return fmt.Errorf("missing recipe path")
	}

	_, err := core.BuildRecipe(recipePath)
	if err != nil {
		return err
	}

	return nil
}
