package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vanilla-os/vib/core"
)

func NewLintCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lint",
		Short: "Lint a Vib recipe file",
		Long:  "Lint a Vib recipe file to check if it is valid",
		RunE:  lintCommand,
	}
	cmd.Flags().SetInterspersed(false)

	return cmd
}

func lintCommand(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no recipe path specified")
	}

	recipePath := args[0]
	err := core.LintRecipe(recipePath)
	if err != nil {
		return err
	}
	return nil
}
