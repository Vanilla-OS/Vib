package api

type Source struct {
	URL      string   `json:"url"`
	Checksum string   `json:"checksum"`
	Type     string   `json:"type"`
	Commit   string   `json:"commit"`
	Tag      string   `json:"tag"`
	Branch   string   `json:"branch"`
	Packages []string `json:"packages"`
	Paths    []string `json:"paths"`
}

type Recipe struct {
	Base          string `json:"base"`
	Name          string
	Id            string
	SingleLayer   bool              `json:"singlelayer"`
	Labels        map[string]string `json:"labels"`
	Adds          map[string]string `json:"adds"`
	Args          map[string]string `json:"args"`
	Runs          []string          `json:"runs"`
	Expose        int               `json:"expose"`
	Cmd           string            `json:"cmd"`
	Modules       []interface{}     `json:"modules"`
	Path          string
	ParentPath    string
	DownloadsPath string
	SourcesPath   string
	PluginPath    string
	Containerfile string
	Entrypoint    []string
}
