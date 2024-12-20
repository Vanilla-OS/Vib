package api

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Generate the destination path for the source based on its type and module name
func GetSourcePath(source Source, moduleName string) string {
	if len(strings.TrimSpace(source.Path)) > 0 {
		return filepath.Join(moduleName, source.Path)
	}
	switch source.Type {
	case "git":
		repoName := strings.Split(source.URL, "/")
		return filepath.Join(moduleName, strings.ReplaceAll(repoName[len(repoName)-1], ".git", ""))
	case "tar":
		url := strings.Split(source.URL, "/")
		tarFile := strings.Split(url[len(url)-1], "?")[0]
		tarParts := strings.Split(tarFile, ".")
		if strings.TrimSpace(tarParts[len(tarParts)-2]) != "tar" {
			return filepath.Join(moduleName, strings.Join(tarParts[:len(tarParts)-1], "."))
		} else {
			return filepath.Join(moduleName, strings.Join(tarParts[:len(tarParts)-2], "."))
		}
	case "file":
		url := strings.Split(source.URL, "/")
		file := strings.Split(url[len(url)-1], "?")[0]
		fileParts := strings.Split(file, ".")
		return filepath.Join(moduleName, strings.Join(fileParts[:len(fileParts)-1], "."))
	case "local":
		toplevelDir := strings.Split(source.URL, "/")
		return filepath.Join(moduleName, toplevelDir[len(toplevelDir)-1])
	}

	return ""
}

// Download the source based on its type and validate its checksum
func DownloadSource(recipe *Recipe, source Source, moduleName string) error {
	fmt.Printf("Downloading source: %s\n", source.URL)

	switch source.Type {
	case "git":
		return DownloadGitSource(recipe.DownloadsPath, source, moduleName)
	case "tar":
		err := DownloadTarSource(recipe.DownloadsPath, source, moduleName)
		if err != nil {
			return err
		}
		return checksumValidation(source, filepath.Join(recipe.DownloadsPath, GetSourcePath(source, moduleName), moduleName+".tar"))
	case "file":
		err := DownloadFileSource(recipe.DownloadsPath, source, moduleName)
		if err != nil {
			return err
		}

		extension := filepath.Ext(source.URL)
		filename := fmt.Sprintf("%s%s", moduleName, extension)
		destinationPath := filepath.Join(recipe.DownloadsPath, GetSourcePath(source, moduleName), filename)

		return checksumValidation(source, destinationPath)
	case "local":
		return DownloadLocalSource(recipe.SourcesPath, source, moduleName)
	default:
		return fmt.Errorf("unsupported source type %s", source.Type)
	}
}

// Clone a specific tag from a Git repository to the destination directory
func gitCloneTag(url, tag, dest string) error {
	cmd := exec.Command(
		"git",
		"clone", url,
		"--depth", "1",
		"--branch", tag,
		dest,
	)
	return cmd.Run()
}

// Retrieve the latest Git repository commit hash for a given branch from the destination directory
func gitGetLatestCommit(branch, dest string) (string, error) {
	cmd := exec.Command("git", "--no-pager", "log", "-n", "1", "--pretty=format:\"%H\"", branch)
	cmd.Dir = dest
	latest_tag, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.Trim(string(latest_tag), "\""), nil
}

// Check out a specific Git repository branch or commit in the destination directory
func gitCheckout(value, dest string) error {
	cmd := exec.Command("git", "checkout", value)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = dest
	return cmd.Run()
}

// Download a Git source repository based on the specified tag, branch, or commit
func DownloadGitSource(downloadPath string, source Source, moduleName string) error {
	fmt.Printf("Downloading git source: %s\n", source.URL)

	if source.URL == "" {
		return fmt.Errorf("missing git remote URL")
	}
	if source.Commit == "" && source.Tag == "" && source.Branch == "" {
		return fmt.Errorf("missing source commit, tag or branch")
	}

	dest := filepath.Join(downloadPath, GetSourcePath(source, moduleName))
	os.MkdirAll(dest, 0o777)

	if source.Tag != "" {
		fmt.Printf("Using tag %s\n", source.Tag)
		return gitCloneTag(source.URL, source.Tag, dest)
	}

	fmt.Printf("Cloning repository: %s\n", source.URL)
	cmd := exec.Command("git", "clone", source.URL, dest)
	err := cmd.Run()
	if err != nil {
		return err
	}

	if source.Commit != "" {
		fmt.Printf("Checking out branch: %s\n", source.Branch)
		err := gitCheckout(source.Branch, dest)
		if err != nil {
			return err
		}
	}

	// Default to latest commit
	if len(strings.TrimSpace(source.Commit)) == 0 || strings.EqualFold(source.Commit, "latest") {
		source.Commit, err = gitGetLatestCommit(source.Branch, dest)
		if err != nil {
			return fmt.Errorf("could not get latest commit: %s", err.Error())
		}
	}
	fmt.Printf("Resetting to commit: %s\n", source.Commit)
	return gitCheckout(source.Commit, dest)
}

// Download a tarball from the specified URL and save it to the destination path
func DownloadTarSource(downloadPath string, source Source, moduleName string) error {
	fmt.Printf("Source is tar: %s\n", source.URL)
	// Create the destination path
	dest := filepath.Join(downloadPath, GetSourcePath(source, moduleName))
	os.MkdirAll(dest, 0o777)
	// Download the resource
	res, err := http.Get(source.URL)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	// Create the destination tar file
	file, err := os.Create(filepath.Join(dest, moduleName+".tar"))
	if err != nil {
		return err
	}
	// Close the file when the function ends
	defer file.Close()
	// Copy the response body to the destination file
	_, err = io.Copy(file, res.Body)
	if err != nil {
		return err
	}

	return nil
}

// Copies a local source for use during the build, skips the Download directory and copies directly into the source path
func DownloadLocalSource(sourcesPath string, source Source, moduleName string) error {
	fmt.Printf("Source is local: %s\n", source.URL)
	dest := filepath.Join(sourcesPath, GetSourcePath(source, moduleName))
	os.MkdirAll(dest, 0o777)
	fileInfo, err := os.Stat(source.URL)
	if err != nil {
		return err
	}
	if fileInfo.IsDir() {
		fmt.Println("ROOTDIR:: ", source.URL)
		root := os.DirFS(source.URL)
		return os.CopyFS(dest, root)
	} else {
		fileName := strings.Split(source.URL, "/")
		out, err := os.Create(filepath.Join(dest, fileName[len(fileName)-1]))
		if err != nil {
			return err
		}
		defer out.Close()

		in, err := os.Open(source.URL)
		if err != nil {
			return err
		}
		defer in.Close()

		_, err = io.Copy(out, in)
		if err != nil {
			return err
		}
		return nil
	}

}

// Move downloaded sources from the download path to the sources path
func MoveSources(downloadPath string, sourcesPath string, sources []Source, moduleName string) error {
	fmt.Println("Moving sources for " + moduleName)

	err := os.MkdirAll(filepath.Join(sourcesPath, moduleName), 0777)
	if err != nil {
		return err
	}
	for _, source := range sources {
		err = MoveSource(downloadPath, sourcesPath, source, moduleName)
		if err != nil {
			return err
		}
	}

	return nil
}

// Move or extract a source from the download path to the sources path depending on its type
// tarballs: extract
// git repositories: move
func MoveSource(downloadPath string, sourcesPath string, source Source, moduleName string) error {
	fmt.Printf("Moving source: %s\n", moduleName)

	err := os.MkdirAll(filepath.Join(sourcesPath, moduleName), 0777)
	if err != nil {
		return err
	}

	switch source.Type {
	case "git", "file":
		dest := GetSourcePath(source, moduleName)
		return os.Rename(
			filepath.Join(downloadPath, dest),
			filepath.Join(sourcesPath, dest),
		)
	case "tar":
		os.MkdirAll(filepath.Join(sourcesPath, GetSourcePath(source, moduleName)), 0o777)
		cmd := exec.Command(
			"tar",
			"-xf", filepath.Join(downloadPath, GetSourcePath(source, moduleName), moduleName+".tar*"),
			"-C", filepath.Join(sourcesPath, GetSourcePath(source, moduleName)),
		)
		err := cmd.Run()
		if err != nil {
			return err
		}

		return os.Remove(filepath.Join(downloadPath, GetSourcePath(source, moduleName), moduleName+".tar"))
	case "local":
		return nil
	default:
		return fmt.Errorf("unsupported source type %s", source.Type)
	}
}

// Validate the checksum of the downloaded file
func checksumValidation(source Source, path string) error {
	// No checksum provided
	if len(strings.TrimSpace(source.Checksum)) == 0 {
		return nil
	}

	// Open the file
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("could not open file: %v", err)
	}

	// Close the file when the function ends
	defer file.Close()

	// Calculate the checksum
	checksum := sha256.New()
	_, err = io.Copy(checksum, file)
	if err != nil {
		return fmt.Errorf("could not calculate checksum: %v", err)
	}

	// Validate the checksum based on source type
	calculatedChecksum := fmt.Sprintf("%x", checksum.Sum(nil))
	if (source.Type == "tar" || source.Type == "file") && calculatedChecksum != source.Checksum {
		return fmt.Errorf("%s source module checksum doesn't match: expected %s, got %s", source.Type, source.Checksum, calculatedChecksum)
	}

	return nil
}

// Download a file source from a URL and save it to the specified download path.
// Create necessary directories and handle file naming based on the URL extension.
func DownloadFileSource(downloadPath string, source Source, moduleName string) error {
	fmt.Printf("Source is file: %s\n", source.URL)

	destDir := filepath.Join(downloadPath, GetSourcePath(source, moduleName))
	os.MkdirAll(destDir, 0o777)
	// Download the resource
	res, err := http.Get(source.URL)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	// Create the destination file
	extension := filepath.Ext(source.URL)
	filename := fmt.Sprintf("%s%s", moduleName, extension)
	dest := filepath.Join(destDir, filename)

	file, err := os.Create(dest)
	if err != nil {
		return err
	}
	// Close the file when the function ends
	defer file.Close()
	// Copy the response body to the destination file
	_, err = io.Copy(file, res.Body)
	if err != nil {
		return err
	}

	return nil
}
