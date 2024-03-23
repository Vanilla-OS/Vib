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
	Name          string
	Id            string
	Stages        []Stage
	Path          string
	ParentPath    string
	DownloadsPath string
	SourcesPath   string
	PluginPath    string
	Containerfile string
}

type Stage struct {
	Id          string            `json:"id"`
	Base        string            `json:"base"`
	SingleLayer bool              `json:"singlelayer"`
	Copy        []Copy            `json:"copy"`
	Labels      map[string]string `json:"labels"`
	Env         map[string]string `json:"env"`
	Adds        map[string]string `json:"adds"`
	Args        map[string]string `json:"args"`
	Runs        []string          `json:"runs"`
	Expose      map[string]string `json:"expose"`
	Cmd         string            `json:"cmd"`
	Modules     []interface{}     `json:"modules"`
	Entrypoint  []string
}

type Copy struct {
	From  string
	Paths []Path
}

type Path struct {
	Src string
	Dst string
}
