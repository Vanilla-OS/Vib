package main

import (
	"os"
	"github.com/vanilla-os/vib/cmd"
)

var (
	Version = cmd.Version
)

func main() {
	err := cmd.Execute()
	if (err != nil) {
		os.Exit(1)
	}
}
