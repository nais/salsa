package build

import (
	"fmt"
	slsa "github.com/in-toto/in-toto-golang/in_toto/slsa_provenance/v0.2"
)

type ArtifactType string

const (
	PkgArtifactType ArtifactType = "pkg"
)

type ArtifactDependencies struct {
	Cmd         Cmd
	RuntimeDeps map[string]Dependency
}

func (in ArtifactDependencies) CmdPath() string {
	return in.Cmd.Path
}

func (in ArtifactDependencies) CmdFlags() string {
	return in.Cmd.CmdFlags
}

type Cmd struct {
	Path     string
	CmdFlags string
}

type Dependency struct {
	Coordinates string
	Version     string
	CheckSum    CheckSum
	Type        string
}

func (d Dependency) ToUri() string {
	return fmt.Sprintf("%s:%s:%s", d.Type, d.Coordinates, d.Version)
}

func (d Dependency) ToDigestSet() slsa.DigestSet {
	return slsa.DigestSet{d.CheckSum.Algorithm: d.CheckSum.Digest}
}

func CreateDependency(coordinates, version string, checksum CheckSum) Dependency {
	return Dependency{
		Coordinates: coordinates,
		Version:     version,
		CheckSum:    checksum,
		Type:        string(PkgArtifactType),
	}
}

type CheckSum struct {
	Algorithm string
	Digest    string
}

func CreateChecksum(algo, digest string) CheckSum {
	return CheckSum{
		Algorithm: algo,
		Digest:    digest,
	}
}
