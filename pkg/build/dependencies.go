package build

import (
	"fmt"
	"io/ioutil"

	slsa "github.com/in-toto/in-toto-golang/in_toto/slsa_provenance/v0.2"
	log "github.com/sirupsen/logrus"
)

type ArtifactType string

const (
	PkgArtifactType ArtifactType = "pkg"
)

type ArtifactDependencies struct {
	Cmd         Cmd
	RuntimeDeps []Dependency
}

type Cmd struct {
	Path     string
	CmdFlags string
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

type BuildTool interface {
	BuildFiles() []string
	ResolveDeps(workDir string) (*ArtifactDependencies, error)
}

type SupportedBuildTools struct {
	Tools []BuildTool
}

func (t SupportedBuildTools) DetectDeps(workDir string) (*ArtifactDependencies, error) {
	log.Info("search for build files\n")
	for _, tool := range t.Tools {
		foundMatch, err := match(tool, workDir)
		if err != nil {
			return nil, fmt.Errorf("could not find match, %v", err)
		}

		if foundMatch {
			log.Infof("found build type '%s'\n", tool.BuildFiles())
			deps, err := tool.ResolveDeps(workDir)
			if err != nil {
				return nil, fmt.Errorf("could not resolve deps, %v", err)
			}

			return deps, nil
		}
	}
	return nil, nil
}

func (in ArtifactDependencies) CmdPath() string {
	return in.Cmd.Path
}

func (in ArtifactDependencies) CmdFlags() string {
	return in.Cmd.CmdFlags
}

func (d Dependency) ToUri() string {
	return fmt.Sprintf("%s:%s:%s", PkgArtifactType, d.Coordinates, d.Version)
}

func (d Dependency) ToDigestSet() slsa.DigestSet {
	return slsa.DigestSet{d.CheckSum.Algorithm: d.CheckSum.Digest}
}

func match(t BuildTool, workDir string) (bool, error) {
	for _, file := range t.BuildFiles() {
		buildFile, err := findBuildFile(workDir, file)

		if err != nil {
			return false, err
		}

		if buildFile != "" {
			return true, nil
		}
	}
	return false, nil
}

func findBuildFile(root, pattern string) (string, error) {
	var result = ""
	files, err := ioutil.ReadDir(root)
	if err != nil {
		return "", fmt.Errorf("reading dir %v", err)
	}

	for _, file := range files {
		if file.Name() == pattern {
			result = file.Name()
			break
		}
	}
	return result, nil
}
