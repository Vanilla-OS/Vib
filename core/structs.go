package core

import "C"

type Module struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Modules []map[string]interface{}
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
	BuildFunc    func(*C.char, *C.char) string
	LoadedPlugin uintptr
}
