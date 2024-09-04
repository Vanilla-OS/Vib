package core

import "C"

// Configuration for a module
type Module struct {
	Name    string `json:"name"`
	Workdir string
	Type    string `json:"type"`
	Modules []map[string]interface{}
	Content []byte // The entire module unparsed as a []byte, used by plugins
}

// Configuration for finalization steps
type Finalize struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Content []byte // The entire module unparsed as a []byte, used by plugins
}

// Configuration for including other modules or recipes
type IncludesModule struct {
	Name     string   `json:"name"`
	Type     string   `json:"type"`
	Includes []string `json:"includes"`
}

// Information for building a module
type ModuleCommand struct {
	Name    string
	Command string
	Workdir string
}

// Configuration for a plugin
type Plugin struct {
	Name         string
	BuildFunc    func(*C.char, *C.char) string
	LoadedPlugin uintptr
}
