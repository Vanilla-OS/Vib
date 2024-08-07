package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/spf13/cobra"
)

var Version = "0.7.4"
var IsRoot = false
var OrigUID = 1000
var OrigGID = 1000
var OrigUser = "user"

var rootCmd = &cobra.Command{
	Use:     "vib",
	Short:   "Vib is a tool to build container images from recipes using modules",
	Long:    "Vib is a tool to build container images from YAML recipes using modules to define the steps to build the image.",
	Version: Version,
}

func init() {
	rootCmd.AddCommand(NewBuildCommand())
	rootCmd.AddCommand(NewTestCommand())
	rootCmd.AddCommand(NewCompileCommand())
}

func Execute() error {
	if os.Getuid() == 0 {
		IsRoot = true
		gid, err := strconv.Atoi(os.Getenv("SUDO_GID"))
		if err != nil {
			return fmt.Errorf("failed to get user uid through SUDO_UID: %s", err.Error())
		}
		OrigGID = gid // go moment??
		uid, err := strconv.Atoi(os.Getenv("SUDO_UID"))
		if err != nil {
			return fmt.Errorf("failed to get user uid through SUDO_GID: %s", err.Error())
		}
		OrigUID = uid
		user := os.Getenv("SUDO_USER")
		os.Setenv("HOME", filepath.Join("/home", user))
		err = syscall.Seteuid(OrigUID)
		if err != nil {
			fmt.Println("WARN: Failed to drop root privileges")
		}
		err = syscall.Setgid(OrigGID)
		if err != nil {
			fmt.Println("WARN: Failed to drop root privileges")
		}
	}
	return rootCmd.Execute()
}
