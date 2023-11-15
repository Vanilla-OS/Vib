package core

import "plugin"

type Recipe struct {
	Base          string `json:"base"`
	Name          string
	Id            string
	SingleLayer   bool              `json:"singlelayer"`
	Labels        map[string]string `json:"labels"`
	Adds          map[string]string `json:"adds"`
	Args          map[string]string `json:"args"`
	Runs          []string          `json:"runs"`
	Cmd           string            `json:"cmd"`
	Modules       []Module          `json:"modules"`
	Path          string
	ParentPath    string
	DownloadsPath string
	SourcesPath   string
	Containerfile string
}

type Module struct {
	Name       string            `json:"name"`
	Type       string            `json:"type"`
	Path       string            `json:"path"`
	Source     Source            `json:"source"`
	Modules    []Module          `json:"modules"`
	BuildFlags string            `json:"buildflags"`
	BuildVars  map[string]string `json:"buildvars"`
	Commands   []string          `json:"commands"`
	Includes   []string          `json:"includes"`
}

type Source struct {
	URL      string   `json:"url"`
	Checksum string   `json:"checksum"`
	Type     string   `json:"type"`
	Commit   string   `json:"commit"`
	Tag      string   `json:"tag"`
	Branch   string   `json:"branch"`
	Packages []string `json:"packages"`
	Paths    []string `json:"paths"`
	Module   string
}

type ModuleCommand struct {
	Name    string
	Command string
}

type Plugin struct {
	Name         string
	BuildFunc    func([]byte) (string, error)
	LoadedPlugin *plugin.Plugin
}
