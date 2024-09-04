package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vanilla-os/vib/core"
)

// Create and return a new test command for the Cobra CLI
func NewTestCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "test",
		Short: "Test the given recipe",
		Long:  "Test the given Vib recipe to check if it's valid",
		RunE:  testCommand,
	}
	cmd.Flags().SetInterspersed(false)

	return cmd
}

// Validate the provided recipe by testing it
func testCommand(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no recipe path specified")
	}

	recipePath := args[0]
	_, err := core.TestRecipe(recipePath)
	if err != nil {
		return err
	}
	return nil
}
