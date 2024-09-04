package core_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/vanilla-os/vib/api"
)

// Test the DownloadSource function to ensure it downloads and verifies the source file
func TestDownloadSource(t *testing.T) {
	tmp := t.TempDir()

	source := api.Source{
		Type:     "tar",
		URL:      "https://github.com/Vanilla-OS/Vib/archive/refs/tags/v0.3.1.tar.gz",
		Checksum: "d28ab888c7b30fd1cc01e0a581169ea52dfb5bfcefaca721497f82734b6a5a98",
	}
	err := api.DownloadSource(tmp, source, "test")
	if err != nil {
		t.Errorf("DownloadSource returned an error: %v", err)
	}

	// Check if the file was downloaded
	dest := filepath.Join(tmp, "test")
	if _, err := os.Stat(dest); os.IsNotExist(err) {
		t.Errorf("Downloaded file does not exist: %v", err)
	}
	defer os.Remove("/tmp/example") // clean up
}

// Test the DownloadTarSource function to ensure it downloads and verifies the tar file
func TestDownloadTarSource(t *testing.T) {
	tmp := t.TempDir()

	source := api.Source{
		Type:     "tar",
		URL:      "https://github.com/Vanilla-OS/Vib/archive/refs/tags/v0.3.1.tar.gz",
		Checksum: "d28ab888c7b30fd1cc01e0a581169ea52dfb5bfcefaca721497f82734b6a5a98",
	}
	err := api.DownloadTarSource(tmp, source, "test2")
	if err != nil {
		t.Errorf("DownloadTarSource returned an error: %v", err)
	}

	// Check if the file was downloaded
	dest := filepath.Join(tmp, "test2")
	if _, err := os.Stat(dest); os.IsNotExist(err) {
		t.Errorf("Downloaded file does not exist: %v", err)
	}
	defer os.Remove("/tmp/example") // clean up

}
