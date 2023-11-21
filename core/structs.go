package core

import (
	"plugin"
)

type Recipe struct {
	Base          string `json:"base"`
	Name          string
	Id            string
	SingleLayer   bool                   `json:"singlelayer"`
	Labels        map[string]string      `json:"labels"`
	Adds          map[string]string      `json:"adds"`
	Args          map[string]string      `json:"args"`
	Runs          []string               `json:"runs"`
	Cmd           string                 `json:"cmd"`
	Modules       map[string]interface{} `json:"modules"`
	Path          string
	ParentPath    string
	DownloadsPath string
	SourcesPath   string
	Containerfile string
}

type Module struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Content []byte // The entire module unparsed as a []byte, used by plugins
}

type IncludesModule struct {
	Name     string   `json:"name"`
	Type     string   `json:"type"`
	Includes []string `json:"includes"`
}

type ModuleCommand struct {
	Name    string
	Command string
}

type Plugin struct {
	Name         string
	BuildFunc    func(interface{}) (string, error)
	LoadedPlugin *plugin.Plugin
}
