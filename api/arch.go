package api

func TestArch(onlyArches []string, targetArch string) bool {
	if len(onlyArches) == 0 {
		return true
	}
	for _, arch := range onlyArches {
		if arch == targetArch {
			return true
		}
	}
	return false
}
