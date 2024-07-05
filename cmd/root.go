package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "vib",
	Short:   "Vib is a tool to build container images from recipes using modules",
	Long:    "Vib is a tool to build container images from YAML recipes using modules to define the steps to build the image.",
	Version: "0.7.3",
}

func init() {
	rootCmd.AddCommand(NewBuildCommand())
	rootCmd.AddCommand(NewTestCommand())
	rootCmd.AddCommand(NewCompileCommand())
}

func Execute() error {
	return rootCmd.Execute()
}
