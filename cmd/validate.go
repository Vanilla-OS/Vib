package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vanilla-os/vib/core"
)

func NewValidateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate the given recipe",
		RunE:  validateCommand,
	}
	cmd.Flags().SetInterspersed(false)

	return cmd
}

func validateCommand(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no recipe path specified")
	}

	recipePath := args[0]
	recipe, err := core.LoadRecipe(recipePath)
	if err != nil {
		return err
	}

	fmt.Printf(
		"Recipe %s looks valid.\nNote that this is not a full validation (yet).\n)",
		recipe.Name,
	)
	return nil
}
