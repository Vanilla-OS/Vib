package core

// BuildMakeModule builds a module that builds a Make project
func BuildMakeModule(module Module) (string, error) {
	return "cd /sources/" + module.Name + " && make && make install", nil
}
