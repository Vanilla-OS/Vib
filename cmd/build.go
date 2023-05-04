package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vanilla-os/vib/core"
)

func NewBuildCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "build",
		Short: "Build the given recipe",
		RunE:  buildCommand,
	}
	cmd.Flags().SetInterspersed(false)

	return cmd
}

func buildCommand(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no recipe path specified")
	}

	recipePath := args[0]
	err := core.BuildRecipe(recipePath)
	if err != nil {
		return err
	}

	return nil
}
