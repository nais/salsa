package scan

type BuildArtifactMetadata struct {
	Deps      map[string]string
	Checksums map[string]CheckSum
}

type ArtifactMetadata struct {
	Name         string
	Dependencies ArtifactDependencies
}

type ArtifactDependencies struct {
	Cmd         string
	RuntimeDeps []Dependency
}

type Dependency struct {
	Coordinates string
	Version     string
	CheckSum    CheckSum
}

func CreateMetadata() *BuildArtifactMetadata {
	return &BuildArtifactMetadata{
		Deps:      make(map[string]string),
		Checksums: make(map[string]CheckSum),
	}
}

type CheckSum struct {
	Algorithm string
	Digest    string
}
