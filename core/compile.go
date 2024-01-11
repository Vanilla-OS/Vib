package core

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/vanilla-os/vib/api"
)

// CompileRecipe compiles a recipe into a runnable image.
func CompileRecipe(recipePath string, runtime string) error {
	recipe, err := BuildRecipe(recipePath)
	if err != nil {
		return err
	}

	storePath := fmt.Sprintf("%s/%s", os.Getenv("HOME"), ".vib/store")
	if _, err := os.Stat(storePath); os.IsNotExist(err) {
		err = os.MkdirAll(storePath, 0755)
		if err != nil {
			return err
		}
	}

	switch runtime {
	case "docker":
		err = compileDocker(recipe, storePath)
		if err != nil {
			return err
		}
	case "podman":
		err = compilePodman(recipe, storePath)
		if err != nil {
			return err
		}
	case "buildah":
		return fmt.Errorf("buildah not implemented yet")
	default:
		return fmt.Errorf("no runtime specified and the prometheus library is not implemented yet")
	}

	fmt.Printf("Image %s built successfully\n", recipe.Id)
	fmt.Printf("Remember to point %s to %s\n", runtime, storePath)

	return nil
}

func compileDocker(recipe api.Recipe, storePath string) error {
	docker, err := exec.LookPath("docker")
	if err != nil {
		return err
	}

	cmd := exec.Command(
		docker, "build",
		"-t", recipe.Id,
		"-f", recipe.Containerfile,
		"--root", storePath,
		".",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = recipe.ParentPath

	return cmd.Run()
}

func compilePodman(recipe api.Recipe, storePath string) error {
	podman, err := exec.LookPath("podman")
	if err != nil {
		return err
	}

	cmd := exec.Command(
		podman, "build",
		"-t", recipe.Id,
		"-f", recipe.Containerfile,
		"--root", storePath,
		".",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = recipe.ParentPath

	return cmd.Run()
}
