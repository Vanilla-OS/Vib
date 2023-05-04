package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "vib",
	Short: "A tool for building and validating",
}

func init() {
	rootCmd.AddCommand(NewBuildCommand())
	rootCmd.AddCommand(NewValidateCommand())
}

func Execute() error {
	return rootCmd.Execute()
}
