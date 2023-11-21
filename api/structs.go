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
	Module   string
}
