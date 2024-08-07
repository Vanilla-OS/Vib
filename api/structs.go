package api

type Source struct {
	URL         string   `json:"url"`
	Checksum    string   `json:"checksum"`
	Type        string   `json:"type"`
	Destination string   `json:"destination"`
	Commit      string   `json:"commit"`
	Tag         string   `json:"tag"`
	Branch      string   `json:"branch"`
	Packages    []string `json:"packages"`
	Paths       []string `json:"paths"`
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
	Finalize      []interface{}
}

type Stage struct {
	Id          string            `json:"id"`
	Base        string            `json:"base"`
	SingleLayer bool              `json:"singlelayer"`
	Copy        []Copy            `json:"copy"`
	Labels      map[string]string `json:"labels"`
	Env         map[string]string `json:"env"`
	Adds        []Add             `json:"adds"`
	Args        map[string]string `json:"args"`
	Runs        Run               `json:"runs"`
	Expose      map[string]string `json:"expose"`
	Cmd         Cmd               `json:"cmd"`
	Modules     []interface{}     `json:"modules"`
	Entrypoint  Entrypoint
}

type PluginType int

const (
	BuildPlugin PluginType = iota
	FinalizePlugin
)

type PluginInfo struct {
	Name string
	Type PluginType
}

type Copy struct {
	From    string
	SrcDst  map[string]string
	Workdir string
}

type Add struct {
	SrcDst  map[string]string
	Workdir string
}

type Entrypoint struct {
	Exec    []string
	Workdir string
}

type Cmd struct {
	Exec    []string
	Workdir string
}

type Run struct {
	Commands []string
	Workdir  string
}
