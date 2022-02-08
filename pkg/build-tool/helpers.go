package build_tool

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func concat(slices ...[]string) (result []string) {
	for _, slice := range slices {
		for _, s := range slice {
			result = append(result, s)
		}
	}
	return result
}
