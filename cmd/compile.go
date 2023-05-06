package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vanilla-os/vib/core"
)

func NewCompileCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "compile",
		Short: "Compile a recipe",
		RunE:  compileCommand,
	}
	cmd.Flags().SetInterspersed(false)

	return cmd
}

func compileCommand(cmd *cobra.Command, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("no recipe path or runtime specified")
	}

	recipePath := args[0]
	runtime := args[1]
	err := core.CompileRecipe(recipePath, runtime)
	if err != nil {
		return err
	}

	return nil
}
