package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vanilla-os/vib/core"
)

func NewLintCommand() *cobra.Command {
	var custom []string
	cmd := &cobra.Command{
		Use:   "lint",
		Short: "Lint a Vib recipe file",
		Long:  "Lint a Vib recipe file to check if it is valid",
		RunE:  lintCommand,
	}

	cmd.Flags().SetInterspersed(false)
	cmd.Flags().StringArrayVarP(&custom, "custom", "c", []string{}, "path to custom module CUE schema")

	return cmd
}

func lintCommand(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no recipe path specified")
	}

	recipePath := args[0]

	if isSet := cmd.Flags().Lookup("custom").Changed; isSet {
		custom, err := cmd.Flags().GetStringArray("custom")

		for _, mod := range custom {
			if _, err := os.Stat(mod); os.IsNotExist(err) {
				return err
			}
		}

		if err != nil {
			return err
		}

		err = core.LintCustomRecipe(recipePath, custom)
		if err != nil {
			return err
		}
	} else {
		err := core.LintRecipe(recipePath)
		if err != nil {
			return err
		}
	}

	return nil
}
