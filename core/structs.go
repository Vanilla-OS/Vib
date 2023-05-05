package core

type Recipe struct {
	Base          string `json:"base"`
	Name          string
	Labels        map[string]string `json:"labels"`
	Args          map[string]string `json:"args"`
	Runs          []string          `json:"runs"`
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
	BuildFlags string            `json:"buildFlags"`
	BuildVars  map[string]string `json:"buildVars"`
}

type Source struct {
	URL      string   `json:"url"`
	Type     string   `json:"type"`
	Commit   string   `json:"commit"`
	Tag      string   `json:"tag"`
	Packages []string `json:"packages"`
	Paths    []string `json:"paths"`
	Module   string
}

type ModuleCommand struct {
	Name    string
	Command string
}
