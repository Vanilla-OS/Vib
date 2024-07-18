package core

import (
	_ "embed"
	"fmt"
	"log"
	"os"

	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/errors"
	"cuelang.org/go/encoding/yaml"
)

//go:embed recipe.cue
var schemaFile string

func validateRecipe(recipeFile []byte) {
	ctx := cuecontext.New()
	schema := ctx.CompileString(schemaFile)
	if err := yaml.Validate(recipeFile, schema); err != nil {
		fmt.Println("❌ Recipe: NOT OK")
		log.Fatal(errors.Details(err, nil))
	}

	fmt.Println("✅ Recipe: OK")
}

// LintRecipe validates a recipe by loading it and checking for errors
func LintRecipe(path string) error {
	recipeFile, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Error reading recipe: %s\n", err)
		return err
	}

	validateRecipe(recipeFile)

	recipe, err := LoadRecipe(path)
	if err != nil {
		fmt.Printf("Error loading recipe: %s\n", err)
		return err
	}

	fmt.Printf("Recipe %s validated successfully\n", recipe.Id)

	modules := 0
	for _, stage := range recipe.Stages {
		modules += len(stage.Modules)
	}

	fmt.Printf("Found %d stages\n", len(recipe.Stages))
	fmt.Printf("Found %d modules\n", modules)
	return nil
}
