package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/vanilla-os/vib/core"
)

// Create and return a new compile command for the Cobra CLI
func NewCompileCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "compile",
		Short: "Compile the given recipe",
		Long:  "Compile the given Vib recipe into a working container image, using the specified runtime (docker/podman)",
		Example: `  vib compile // using the recipe in the current directory and the system's default runtime
  vib compile --runtime podman // using the recipe in the current directory and Podman as the runtime
  vib compile /path/to/recipe.yml --runtime podman // using the recipe at the specified path and Podman as the runtime
  Both docker and podman are supported as runtimes. If none is specified, the detected runtime will be used, giving priority to Docker.`,
		RunE: compileCommand,
	}
	cmd.Flags().StringP("runtime", "r", "", "The runtime to use (docker/podman)")
	cmd.Flags().SetInterspersed(false)

	return cmd
}

// Execute the compile command: compile the given recipe into a container image
func compileCommand(cmd *cobra.Command, args []string) error {
	commonNames := []string{
		"recipe.yml",
		"recipe.yaml",
		"vib.yml",
		"vib.yaml",
	}
	var recipePath string
	var arch string
	var containerRuntime string

	arch = runtime.GOARCH
	containerRuntime, _ = cmd.Flags().GetString("runtime")

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

	detectedRuntime := detectRuntime()
	if containerRuntime == "" && detectedRuntime == "" {
		return fmt.Errorf("missing runtime, and no one was detected")
	} else if containerRuntime == "" {
		containerRuntime = detectedRuntime
	}

	err := core.CompileRecipe(recipePath, arch, containerRuntime, IsRoot, OrigGID, OrigUID)
	if err != nil {
		return err
	}

	return nil
}

// Detect the container runtime by checking the system path
//
// Returns: runtime name or an empty string if no runtime is found
func detectRuntime() string {
	path, _ := exec.LookPath("docker")
	if path != "" {
		return "docker"
	}

	path, _ = exec.LookPath("podman")
	if path != "" {
		return "podman"
	}

	return ""
}
