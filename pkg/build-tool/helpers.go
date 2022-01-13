package build_tool

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func sumSupported(buildFiles ...[]string) (supported []string) {
	for _, supportedBuildFiles := range buildFiles {
		for _, s := range supportedBuildFiles {
			supported = append(supported, s)
		}
	}
	return supported
}
