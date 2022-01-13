package scan

type BuildToolMetadata struct {
	Deps      map[string]string
	Checksums map[string]string
}

func CreateMetadata() *BuildToolMetadata {
	return &BuildToolMetadata{
		Deps:      make(map[string]string),
		Checksums: make(map[string]string),
	}
}
