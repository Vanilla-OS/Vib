package main

import (
	"github.com/vanilla-os/vib/cmd"
)

var (
	Version = cmd.Version
)

func main() {
	cmd.Execute()
}
