package core

import (
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/vanilla-os/vib/api"
	"gopkg.in/yaml.v3"
)

var Min_Recipe_Version = []uint8{1, 0, 0}

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

	if len(strings.TrimSpace(recipe.Vibversion)) <= 0 {
		return nil, fmt.Errorf("version key not found in recipe file, assuming outdated recipe")
	}

	recipeVersionS := strings.Split(recipe.Vibversion, ".")
	if len(recipeVersionS) != 3 {
		return nil, fmt.Errorf("invalid version format, expected x.x.x, got %s", recipe.Vibversion)
	}

	recipeVersion := []uint8{0, 0, 0}
	for i := 0; i < len(recipeVersion); i++ {
		versionInt, err := strconv.ParseUint(recipeVersionS[i], 10, 8)
		if err != nil {
			return nil, err
		}
		if versionInt > math.MaxUint8 {
			recipeVersion[i] = math.MaxUint8
		} else {
			recipeVersion[i] = uint8(versionInt)
		}
	}

	if recipeVersion[0] < Min_Recipe_Version[0] || recipeVersion[1] < Min_Recipe_Version[1] || recipeVersion[2] < Min_Recipe_Version[2] {
		return nil, fmt.Errorf("outdated recipe, this version of vib supports recipes starting at version %s", strings.Join(strings.Fields(fmt.Sprint(Min_Recipe_Version)), "."))
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

	// the plugins directory contains all plugins that vib can load
	// and use for unknown modules in the recipe
	recipe.PluginPath = filepath.Join(filepath.Dir(recipePath), "plugins")

	// the includes directory is the place where we store all the
	// files to be included in the container, this is useful for
	// example to include configuration files. Each file must follow
	// the File Hierarchy Standard (FHS) and be placed in the correct
	// directory. For example, if you want to include a file in
	// /etc/nginx/nginx.conf you must place it in includes/etc/nginx/nginx.conf
	// so it will be copied to the correct location in the container
	if len(strings.TrimSpace(recipe.IncludesPath)) == 0 {
		recipe.IncludesPath = filepath.Join("includes.container")
	}
	_, err = os.Stat(recipe.IncludesPath)
	if os.IsNotExist(err) {
		err := os.MkdirAll(recipe.IncludesPath, 0755)
		if err != nil {
			return nil, err
		}
	}

	for i, stage := range recipe.Stages {
		// here we check if the extra Adds path exists
		for _, add := range stage.Adds {
			for src := range add.SrcDst {
				fullPath := filepath.Join(filepath.Dir(recipePath), src)
				_, err = os.Stat(fullPath)
				if os.IsNotExist(err) {
					return nil, err
				}
			}
		}

		recipe.Stages[i] = stage
	}

	return recipe, nil
}

// downloadRecipe downloads a recipe from a remote URL and stores it to
// a temporary file
func downloadRecipe(url string) (path string, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	tmpFile, err := os.CreateTemp("", "vib-recipe-")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		return "", err
	}

	return tmpFile.Name(), nil
}

// followsGhPattern checks if a given path follows the pattern:
// gh:org/repo:branch:path
func followsGhPattern(s string) bool {
	parts := strings.Split(s, ":")
	if len(parts) != 4 {
		return false
	}

	if parts[0] != "gh" {
		return false
	}

	return true
}

// downloadGhRecipe downloads a recipe from a github repository and stores it to
// a temporary file
func downloadGhRecipe(gh string) (path string, err error) {
	parts := strings.Split(gh, ":")
	repo := parts[1]
	branch := parts[2]
	file := parts[3]

	url := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s", repo, branch, file)
	return downloadRecipe(url)
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

// TestRecipe validates a recipe by loading it and checking for errors
func TestRecipe(path string) (*api.Recipe, error) {
	recipe, err := LoadRecipe(path)
	if err != nil {
		fmt.Printf("Error validating recipe: %s\n", err)
		return nil, err
	}

	modules := 0
	for _, stage := range recipe.Stages {
		modules += len(stage.Modules)
	}

	fmt.Printf("Recipe %s validated successfully\n", recipe.Id)
	fmt.Printf("Found %d stages\n", len(recipe.Stages))
	fmt.Printf("Found %d modules\n", modules)
	return recipe, nil
}
