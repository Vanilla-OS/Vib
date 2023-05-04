package core

// BuildAptModule builds a module that installs packages
// using the apt package manager
func BuildAptModule(module Module) (string, error) {
	packages := ""
	for _, pkg := range module.Source.Packages {
		packages += pkg + " "
	}

	return "apt install -y " + packages, nil
}
