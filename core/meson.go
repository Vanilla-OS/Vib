package core

import (
	"fmt"

	"github.com/google/uuid"
)

// BuildMesonModule builds a module that builds a Meson project
func BuildMesonModule(module Module) (string, error) {
	tmpDir := "/tmp/" + uuid.New().String()

	cmd := fmt.Sprintf(
		"cd /sources/%s && meson %s && ninja -C %s && ninja -C %s install",
		module.Name,
		tmpDir,
		tmpDir,
		tmpDir,
	)

	return cmd, nil
}
