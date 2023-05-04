package core

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

// ResolveSources resolves the sources of a recipe and downloads them
// to the downloads directory. Note that modules in this function are
// returned in the order they should be built.
func ResolveSources(recipe *Recipe) ([]Module, []Source, error) {
	modules := GetAllModules(recipe.Modules)
	var sources []Source

	for _, module := range modules {
		if module.Source.URL == "" {
			continue
		}

		module.Source.Module = module.Name
		err := DownloadSource(recipe, module.Source)
		if err != nil {
			return nil, nil, err
		}

		sources = append(sources, module.Source)
	}

	return modules, sources, nil
}

// GetAllModules returns a list of all modules in a ordered list
func GetAllModules(modules []Module) []Module {
	var orderedList []Module

	for _, module := range modules {
		orderedList = append(orderedList, GetAllModules(module.Modules)...)
		orderedList = append(orderedList, module)
	}

	return orderedList
}

// DownloadSource downloads a source to the downloads directory
// according to its type (git, tar, ...)
func DownloadSource(recipe *Recipe, source Source) error {
	if source.Type == "git" {
		return DownloadGitSource(recipe, source)
	} else if source.Type == "tar" {
		return DownloadTarSource(recipe, source)
	} else {
		return fmt.Errorf("unsupported source type %s", source.Type)
	}
}

// DownloadGitSource downloads a git source to the downloads directory
// and checks out the commit or tag
func DownloadGitSource(recipe *Recipe, source Source) error {
	dest := filepath.Join(recipe.DownloadsPath, source.Module)

	if source.Commit == "" && source.Tag == "" {
		return fmt.Errorf("missing source commit or tag")
	}

	if source.Tag != "" {
		cmd := exec.Command(
			"git",
			"clone", source.URL,
			"--depth", "1",
			"--branch", source.Tag,
			dest,
		)
		err := cmd.Run()
		if err != nil {
			return err
		}
	} else {
		cmd := exec.Command(
			"git",
			"clone", source.URL,
			"--depth", "1",
			dest,
		)
		err := cmd.Run()
		if err != nil {
			return err
		}

		cmd = exec.Command(
			"git",
			"checkout", source.Commit,
		)
		cmd.Dir = dest
		err = cmd.Run()
		if err != nil {
			return err
		}
	}

	return nil
}

// DownloadTarSource downloads a tar archive to the downloads directory
func DownloadTarSource(recipe *Recipe, source Source) error {
	dest := filepath.Join(recipe.DownloadsPath, source.Module)

	res, err := http.Get(source.URL)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	file, err := os.Create(dest)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = io.Copy(file, res.Body)
	if err != nil {
		return err
	}

	return nil
}

// MoveSources moves all sources from the downloads directory to the
// sources directory
func MoveSources(recipe *Recipe, sources []Source) error {
	for _, source := range sources {
		err := MoveSource(recipe, source)
		if err != nil {
			return err
		}
	}

	return nil
}

// MoveSource moves a source from the downloads directory to the
// sources directory, by extracting if a tar archive or moving if a
// git repository
func MoveSource(recipe *Recipe, source Source) error {
	if source.Type == "git" {
		return os.Rename(
			filepath.Join(recipe.DownloadsPath, source.Module),
			filepath.Join(recipe.SourcesPath, source.Module),
		)
	} else if source.Type == "tar" {
		cmd := exec.Command(
			"tar",
			"-xf", filepath.Join(recipe.DownloadsPath, source.Module),
			"-C", recipe.SourcesPath,
		)
		err := cmd.Run()
		if err != nil {
			return err
		}

		return os.Remove(filepath.Join(recipe.DownloadsPath, source.Module))
	} else {
		return fmt.Errorf("unsupported source type %s", source.Type)
	}
}
