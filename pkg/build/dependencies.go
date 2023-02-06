package build

import (
	"fmt"
	slsa "github.com/in-toto/in-toto-golang/in_toto/slsa_provenance/common"
)

type artifactType string

const (
	pkgArtifactType artifactType = "pkg"
)

type ArtifactDependencies struct {
	Cmd         Cmd
	RuntimeDeps map[string]Dependency
}

func ArtifactDependency(deps map[string]Dependency, path, cmdFlags string) *ArtifactDependencies {
	return &ArtifactDependencies{
		Cmd: Cmd{
			Path:     path,
			CmdFlags: cmdFlags,
		},
		RuntimeDeps: deps,
	}
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

func Dependence(coordinates, version string, checksum CheckSum) Dependency {
	return Dependency{
		Coordinates: coordinates,
		Version:     version,
		CheckSum:    checksum,
		Type:        string(pkgArtifactType),
	}
}

type CheckSum struct {
	Algorithm string
	Digest    string
}

func Verification(algo, digest string) CheckSum {
	return CheckSum{
		Algorithm: algo,
		Digest:    digest,
	}
}
