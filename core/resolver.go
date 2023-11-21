package core

import (
	"crypto/sha256"
	"fmt"
	"github.com/vanilla-os/vib/api"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ResolveSources resolves the sources of a recipe and downloads them
// to the downloads directory. Note that modules in this function are
// returned in the order they should be built.
func ResolveSources(recipe *Recipe) ([]Module, []api.Source, error) {
	fmt.Println("Resolving sources")

	modules := GetAllModules(recipe.Modules)
	//var sources []Source

	for _, module := range modules {
		fmt.Printf("Resolving source for: %s\n", module.Name)

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

// GetAllModules returns a list of all modules in an ordered list
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
	fmt.Printf("Downloading source: %s\n", source.URL)

	if source.Type == "git" {
		return DownloadGitSource(recipe, source)
	} else if source.Type == "tar" {
		err := DownloadTarSource(recipe, source)
		if err != nil {
			return err
		}
		return checksumValidation(source, filepath.Join(recipe.DownloadsPath, source.Module))
	} else {
		return fmt.Errorf("unsupported source type %s", source.Type)
	}
}

// DownloadGitSource downloads a git source to the downloads directory
// and checks out the commit or tag
func DownloadGitSource(recipe *Recipe, source Source) error {
	fmt.Printf("Source is git: %s\n", source.URL)

	dest := filepath.Join(recipe.DownloadsPath, source.Module)

	if source.Commit == "" && source.Tag == "" && source.Branch == "" {
		return fmt.Errorf("missing source commit, tag or branch")
	}

	if source.Tag != "" {
		fmt.Printf("Using a tag: %s\n", source.Tag)

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
		fmt.Printf("Using a commit: %s\n", source.Commit)

		if source.Branch == "" {
			return fmt.Errorf("missing source branch, needed to checkout commit")
		}

		fmt.Printf("Cloning repository: %s\n", source.URL)
		cmd := exec.Command(
			"git",
			"clone", source.URL,
			dest,
		)
		err := cmd.Run()
		if err != nil {
			return err
		}

		if source.Commit == "latest" {
			cmd := exec.Command(
				"git", "--no-pager", "log", "-n", "1", "--pretty=format:\"%H\"", source.Branch,
			)
			cmd.Dir = dest
			latest_tag, err := cmd.Output()
			if err != nil {
				return err
			}
			source.Commit = strings.Trim(string(latest_tag), "\"")
		}

		fmt.Printf("Checking out branch: %s\n", source.Branch)
		cmd = exec.Command(
			"git",
			"checkout",
			"-B", source.Branch,
		)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Dir = dest
		err = cmd.Run()
		if err != nil {
			return err
		}

		fmt.Printf("Resetting to commit: %s\n", source.Commit)
		cmd = exec.Command(
			"git",
			"reset", "--hard", source.Commit,
		)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
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
	fmt.Printf("Source is tar: %s\n", source.URL)
	//Create the destination path
	dest := filepath.Join(recipe.DownloadsPath, source.Module)
	//Download the resource
	res, err := http.Get(source.URL)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	//Create the destination tar file
	file, err := os.Create(dest)
	if err != nil {
		return err
	}
	//Close the file when the function ends
	defer file.Close()
	//Copy the response body to the destination file
	_, err = io.Copy(file, res.Body)
	if err != nil {
		return err
	}

	return nil
}

// MoveSources moves all sources from the downloads directory to the
// sources directory
func MoveSources(recipe *Recipe, sources []Source) error {
	fmt.Println("Moving sources")

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
	fmt.Printf("Moving source: %s\n", source.Module)

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

// checksumValidation validates the checksum of a file
func checksumValidation(source Source, path string) error {
	//No checksum provided
	if len(strings.TrimSpace(source.Checksum)) == 0 {
		return nil
	}
	//Open the file
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	//Close the file when the function ends
	defer file.Close()
	//Calculate the checksum
	checksum := sha256.New()
	_, err = io.Copy(checksum, file)
	if err != nil {
		return fmt.Errorf("could not calculate tar file checksum")
	}

	//Validate the checksum
	if fmt.Sprintf("%x", checksum.Sum(nil)) != source.Checksum {

		return fmt.Errorf("tar file checksum doesn't match")
	}

	return nil
}
