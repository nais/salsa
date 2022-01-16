package scan

type BuildToolMetadata struct {
	Deps      map[string]string
	Checksums map[string]CheckSum
}

func CreateMetadata() *BuildToolMetadata {
	return &BuildToolMetadata{
		Deps:      make(map[string]string),
		Checksums: make(map[string]CheckSum),
	}
}

type CheckSum struct {
	Algorithm string
	Digest    string
}
