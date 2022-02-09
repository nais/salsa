package common

import (
	"fmt"

	slsa "github.com/in-toto/in-toto-golang/in_toto/slsa_provenance/v0.2"
)

type ArtifactType string

const (
	PkgArtifactType ArtifactType = "pkg"
)

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

type CheckSum struct {
	Algorithm string
	Digest    string
}

func (d Dependency) ToUri() string {
	return fmt.Sprintf("%s:%s:%s", PkgArtifactType, d.Coordinates, d.Version)
}

func (d Dependency) ToDigestSet() slsa.DigestSet {
	return slsa.DigestSet{d.CheckSum.Algorithm: d.CheckSum.Digest}
}

// TODO: rename?
type BuildTool interface {
	BuildFiles() []string
	ResolveDeps(workDir string) (*ArtifactDependencies, error)
}

type BuildTools struct {
	Tools []BuildTool
}
