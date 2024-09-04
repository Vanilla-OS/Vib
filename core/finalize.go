package core

import (
	"fmt"
	cstorage "github.com/containers/storage"
	"os/exec"
	"strings"
)

// Configuration for storage drivers
type StorageConf struct {
	Driver    string
	Runroot   string
	Graphroot string
}

// Retrieve the container storage configuration based on the runtime
func GetContainerStorage(runtime string) (cstorage.Store, error) {
	storageconfig := &StorageConf{}
	if runtime == "podman" {
		podmanPath, err := exec.LookPath("podman")
		out, err := exec.Command(
			podmanPath, "info", "-f json").Output()
		if err != nil {
			fmt.Println("Failed to get podman info")
		} else {
			driver := strings.Split(strings.Split(string(out), "\"graphDriverName\": \"")[1], "\",")[0]
			storageconfig.Driver = driver

			graphRoot := strings.Split(strings.Split(string(out), "\"graphRoot\": \"")[1], "\",")[0]
			storageconfig.Graphroot = graphRoot

			runRoot := strings.Split(strings.Split(string(out), "\"runRoot\": \"")[1], "\",")[0]
			storageconfig.Runroot = runRoot
		}

	}
	if storageconfig.Runroot == "" {
		storageconfig.Runroot = "/var/lib/vib/runroot"
		storageconfig.Graphroot = "/var/lib/vib/graphroot"
		storageconfig.Driver = "overlay"
	}
	store, err := cstorage.GetStore(cstorage.StoreOptions{
		RunRoot:         storageconfig.Runroot,
		GraphRoot:       storageconfig.Graphroot,
		GraphDriverName: storageconfig.Driver,
	})
	if err != nil {
		return store, err
	}

	return store, err
}

// Retrieve the image ID for a given image name from the storage
func GetImageID(name string, store cstorage.Store) (string, error) {
	images, err := store.Images()
	if err != nil {
		return "", err
	}
	for _, img := range images {
		for _, imgname := range img.Names {
			if imgname == name {
				return img.ID, nil
			}
		}
	}
	return "", fmt.Errorf("image not found")
}

// Retrieve the top layer ID for a given image ID from the storage
func GetTopLayerID(imageid string, store cstorage.Store) (string, error) {
	images, err := store.Images()
	if err != nil {
		return "", err
	}
	for _, img := range images {
		if img.ID == imageid {
			return img.TopLayer, nil
		}
	}
	return "", fmt.Errorf("no top layer for id %s found", imageid)
}

// Mount the image and return the mount directory
func MountImage(imagename string, imageid string, runtime string) (string, error) {
	store, err := GetContainerStorage(runtime)
	if err != nil {
		return "", err
	}
	topLayerID, err := GetTopLayerID(imageid, store)
	if err != nil {
		return "", err
	}
	mountDir, err := store.Mount(topLayerID, "")
	if err != nil {
		return "", err
	}
	return mountDir, err
}
