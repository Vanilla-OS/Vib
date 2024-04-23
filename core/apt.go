package core

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mitchellh/mapstructure"
	"github.com/vanilla-os/vib/api"
)

type AptModule struct {
	Name    string     `json:"name"`
	Type    string     `json:"type"`
	Options AptOptions `json:"options"`
	Source  api.Source `json:"source"`
}

type AptOptions struct {
	NoRecommends    bool `json:"no_recommends"`
	InstallSuggests bool `json:"install_suggests"`
	FixMissing      bool `json:"fix_missing"`
	FixBroken       bool `json:"fix_broken"`
}

// BuildAptModule builds a module that installs packages
// using the apt package manager
func BuildAptModule(moduleInterface interface{}, recipe *api.Recipe) (string, error) {
	var module AptModule
	err := mapstructure.Decode(moduleInterface, &module)
	if err != nil {
		return "", err
	}

	args := ""
	if module.Options.NoRecommends {
		args += "--no-install-recommends "
	}
	if module.Options.InstallSuggests {
		args += "--install-suggests "
	}
	if module.Options.FixMissing {
		args += "--fix-missing "
	}
	if module.Options.FixBroken {
		args += "--fix-broken "
	}

	if len(module.Source.Packages) > 0 {
		packages := ""
		for _, pkg := range module.Source.Packages {
			packages += pkg + " "
		}

		return fmt.Sprintf("apt install -y %s %s && apt clean", args, packages), nil
	}

	if len(module.Source.Paths) > 0 {
		cmd := ""

		for i, path := range module.Source.Paths {
			instPath := filepath.Join(recipe.ParentPath, path+".inst")
			packages := ""
			file, err := os.Open(instPath)
			if err != nil {
				return "", err
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				packages += scanner.Text() + " "
			}

			if err := scanner.Err(); err != nil {
				return "", err
			}

			cmd += fmt.Sprintf("apt install -y %s %s", args, packages)

			if i != len(module.Source.Paths)-1 {
				cmd += "&& "
			} else {
				cmd += "&& apt clean"
			}
		}

		return cmd, nil
	}

	return "", errors.New("no packages or paths specified")
}

