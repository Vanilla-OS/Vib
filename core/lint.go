package core

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/errors"
	"cuelang.org/go/encoding/yaml"
	"github.com/vanilla-os/vib/api"
)

//go:embed recipe.cue
var schemaFile string

func getModuleType(val cue.Value) (string, error) {
	definitions := []string{}
	it, err := val.Fields(cue.Definitions(true))
	if err != nil {
		log.Fatal(err)
	}
	for it.Next() {
		if !it.Selector().IsDefinition() {
			continue
		}
		definitions = append(definitions, it.Selector().String())
	}

	if len(definitions) == 1 {
		return definitions[0], nil
	}

	return "", errors.New("Multiple definitions in custom module schema")
}

func validateRecipe(recipeFile []byte) {
	ctx := cuecontext.New()
	schema := ctx.CompileString(schemaFile)
	if err := yaml.Validate(recipeFile, schema); err != nil {
		fmt.Println("❌ Recipe: NOT OK")
		log.Fatal(errors.Details(err, nil))
	}

	fmt.Println("✅ Recipe: OK")
}

func validateCustomRecipe(recipeFile []byte, custom []string) {
	modules := []string{}
	types := []string{}
	pairs := make(map[string]string)

	ctx := cuecontext.New()
	for _, mod := range custom {
		data, err := os.ReadFile(mod)
		if err != nil {
			fmt.Println("Custom module schema could not be read")
			log.Fatal(err, nil)
		}
		customSchema := ctx.CompileBytes(data)
		moduleType, err := getModuleType(customSchema)
		if err != nil {
			log.Fatal(err, nil)
		}
		moduleTypeName := strings.TrimSuffix(filepath.Base(mod), filepath.Ext(mod))
		modules = append(modules, string(data))
		pairs[moduleTypeName] = moduleType
	}

	cueString := `
    #ModuleTypes: {
      "%s": %s
		}
	`

	for k, v := range pairs {
		types = append(types, fmt.Sprintf(cueString, k, v))
	}

	schema := ctx.CompileString(schemaFile + "\n" + strings.Join(modules, "\n") + "\n" + strings.Join(types, "\n"))

	if err := yaml.Validate(recipeFile, schema); err != nil {
		fmt.Println("❌ Recipe: NOT OK")
		log.Fatal(errors.Details(err, nil))
	}

	fmt.Println("✅ Recipe: OK")
}

func readRecipe(path string) ([]byte, error) {
	recipeFile, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return recipeFile, nil
}

func loadRecipe(path string) (*api.Recipe, error) {
	recipe, err := LoadRecipe(path)
	if err != nil {
		return nil, err
	}

	return recipe, nil
}

func displayRecipeInfo(recipe *api.Recipe) {
	modules := 0
	for _, stage := range recipe.Stages {
		modules += len(stage.Modules)
	}

	fmt.Printf("Found %d stages\n", len(recipe.Stages))
	fmt.Printf("Found %d modules\n", modules)
}

// LintRecipe validates a recipe by loading it and checking for errors
func LintRecipe(path string) error {
	recipeFile, err := readRecipe(path)

	if err != nil {
		fmt.Printf("Error reading YAML recipe: %s\n", err)
		return err
	}

	validateRecipe(recipeFile)

	recipe, err := loadRecipe(path)
	if err != nil {
		fmt.Printf("Error loading Vib recipe: %s\n", err)
		return err
	}

	fmt.Printf("Recipe %s validated successfully\n", recipe.Id)

	displayRecipeInfo(recipe)

	return nil
}

// LintCustomRecipe validates a recipe including custom modules by loading it and checking for errors
func LintCustomRecipe(path string, custom []string) error {
	recipeFile, err := readRecipe(path)

	if err != nil {
		fmt.Printf("Error reading YAML recipe: %s\n", err)
		return err
	}

	validateCustomRecipe(recipeFile, custom)

	recipe, err := loadRecipe(path)
	if err != nil {
		fmt.Printf("Error loading Vib recipe: %s\n", err)
		return err
	}

	fmt.Printf("Recipe %s validated successfully\n", recipe.Id)

	displayRecipeInfo(recipe)

	return nil

}
