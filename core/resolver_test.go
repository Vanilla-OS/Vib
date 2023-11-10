package core_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/vanilla-os/vib/core"
)

func TestDownloadSource(t *testing.T) {
	recipe := &core.Recipe{
		DownloadsPath: "/tmp/",
	}
	source := core.Source{
		Type:     "tar",
		URL:      "https://github.com/Vanilla-OS/Vib/archive/refs/tags/v0.3.1.tar.gz",
		Module:   "example",
		Checksum: "d28ab888c7b30fd1cc01e0a581169ea52dfb5bfcefaca721497f82734b6a5a98",
	}
	err := core.DownloadSource(recipe, source)
	if err != nil {
		t.Errorf("DownloadSource returned an error: %v", err)
	}

	// Check if the file was downloaded
	dest := filepath.Join(recipe.DownloadsPath, source.Module)
	if _, err := os.Stat(dest); os.IsNotExist(err) {
		t.Errorf("Downloaded file does not exist: %v", err)
	}
	defer os.Remove("/tmp/example") // clean up
}

func TestDownloadTarSource(t *testing.T) {
	recipe := &core.Recipe{
		DownloadsPath: "/tmp/",
	}
	source := core.Source{
		Type:     "tar",
		URL:      "https://github.com/Vanilla-OS/Vib/archive/refs/tags/v0.3.1.tar.gz",
		Module:   "example",
		Checksum: "d28ab888c7b30fd1cc01e0a581169ea52dfb5bfcefaca721497f82734b6a5a98",
	}
	err := core.DownloadTarSource(recipe, source)
	if err != nil {
		t.Errorf("DownloadTarSource returned an error: %v", err)
	}

	// Check if the file was downloaded
	dest := filepath.Join(recipe.DownloadsPath, source.Module)
	if _, err := os.Stat(dest); os.IsNotExist(err) {
		t.Errorf("Downloaded file does not exist: %v", err)
	}
	defer os.Remove("/tmp/example") // clean up

}
