package core

import (
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/vanilla-os/vib/api"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// LoadRecipe loads a recipe from a file and returns a Recipe
// Does not validate the recipe but it will catch some errors
// a proper validation will be done in the future
func LoadRecipe(path string) (*api.Recipe, error) {
	recipe := &api.Recipe{}

	// we use the absolute path to the recipe file as the
	// root path for the recipe and all its files
	recipePath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	// here we open the recipe file and unmarshal it into
	// the Recipe struct, this is not a full validation
	// but it will catch some errors
	recipeFile, err := os.Open(recipePath)
	if err != nil {
		return nil, err
	}
	defer recipeFile.Close()

	recipeYAML, err := io.ReadAll(recipeFile)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(recipeYAML, recipe)
	if err != nil {
		return nil, err
	}

	// the recipe path is stored in the recipe itself
	// for convenience
	recipe.Path = recipePath
	recipe.ParentPath = filepath.Dir(recipePath)

	// assuming the Containerfile location is relative
	recipe.Containerfile = filepath.Join(filepath.Dir(recipePath), "Containerfile")
	err = os.RemoveAll(recipe.Containerfile)
	if err != nil {
		return nil, err
	}

	// we create the sources directory which is the place where
	// all the sources will be stored and be available to all
	// the modules
	recipe.SourcesPath = filepath.Join(filepath.Dir(recipePath), "sources")
	err = os.RemoveAll(recipe.SourcesPath)
	if err != nil {
		return nil, err
	}
	err = os.MkdirAll(recipe.SourcesPath, 0755)
	if err != nil {
		return nil, err
	}

	// the downloads directory is a transient directory, here all
	// the downloaded sources will be stored before being moved
	// to the sources directory. This is useful since some sources
	// types need to be extracted, this way we can extract them
	// directly to the sources directory after downloading them
	recipe.DownloadsPath = filepath.Join(filepath.Dir(recipePath), "downloads")
	err = os.RemoveAll(recipe.DownloadsPath)
	if err != nil {
		return nil, err
	}
	err = os.MkdirAll(recipe.DownloadsPath, 0755)
	if err != nil {
		return nil, err
	}

	// the includes directory is the place where we store all the
	// files to be included in the container, this is useful for
	// example to include configuration files. Each file must follow
	// the File Hierachy Standard (FHS) and be placed in the correct
	// directory. For example, if you want to include a file in
	// /etc/nginx/nginx.conf you must place it in includes/etc/nginx/nginx.conf
	// so it will be copied to the correct location in the container
	includesContainerPath := filepath.Join(filepath.Dir(recipePath), "includes.container")
	_, err = os.Stat(includesContainerPath)
	if os.IsNotExist(err) {
		err := os.MkdirAll(includesContainerPath, 0755)
		if err != nil {
			return nil, err
		}
	}

	// here we check if the extra Adds path exists
	for src := range recipe.Adds {
		fullPath := filepath.Join(filepath.Dir(recipePath), src)
		_, err = os.Stat(fullPath)
		if os.IsNotExist(err) {
			return nil, err
		}
	}

	// here we expand modules of type "includes"
	var newRecipeModules []interface{}

	for _, moduleInterface := range recipe.Modules {

		var module Module
		err := mapstructure.Decode(moduleInterface, &module)
		if err != nil {
			return nil, err
		}

		if module.Type == "includes" {
			var include IncludesModule
			err := mapstructure.Decode(moduleInterface, &include)
			if err != nil {
				return nil, err
			}

			if len(include.Includes) == 0 {
				return nil, errors.New("includes module must have at least one module to include")
			}

			for _, include := range include.Includes {
				includeModule, err := GenModule(filepath.Join(recipe.ParentPath, include+".yml"))
				fmt.Printf("!!!!adding new module %s\n", includeModule)
				if err != nil {
					return nil, err
				}

				newRecipeModules = append(newRecipeModules, includeModule)
			}

			continue
		}

		newRecipeModules = append(newRecipeModules, moduleInterface)
	}

	recipe.Modules = newRecipeModules

	return recipe, nil
}

// GenModule generate a Module struct from a module path
func GenModule(modulePath string) (map[string]interface{}, error) {
	var module map[string]interface{}

	moduleFile, err := os.Open(modulePath)
	if err != nil {
		return module, err
	}
	defer moduleFile.Close()

	moduleYAML, err := io.ReadAll(moduleFile)
	if err != nil {
		return module, err
	}

	err = yaml.Unmarshal(moduleYAML, &module)
	if err != nil {
		return module, err
	}

	return module, nil
}
