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

type BuildTool interface {
	BuildFiles() []string
	ResolveDeps(workDir string) (*ArtifactDependencies, error)
}

type SupportedBuildTools struct {
	Tools []BuildTool
}

func (t SupportedBuildTools) DetectDeps(workDir string) (*ArtifactDependencies, error) {
	for _, tool := range t.Tools {
		log.Infof("search for build files '%s'\n", tool.BuildFiles())
		// TODO handle error
		if match(tool, workDir) {
			log.Infof("found build type")
			deps, err := tool.ResolveDeps(workDir)
			if err != nil {
				return nil, fmt.Errorf("could not resolve deps, %v", err)
			}
			return deps, nil
		}
	}
	return nil, nil
}

func (d Dependency) ToUri() string {
	return fmt.Sprintf("%s:%s:%s", PkgArtifactType, d.Coordinates, d.Version)
}

func (d Dependency) ToDigestSet() slsa.DigestSet {
	return slsa.DigestSet{d.CheckSum.Algorithm: d.CheckSum.Digest}
}

func match(t BuildTool, workDir string) bool {
	for _, file := range t.BuildFiles() {
		buildFile := findBuildFile(workDir, file)
		if buildFile != "" {
			return true
		}
	}
	return false
}

func findBuildFile(root, pattern string) (result string) {
	files, err := ioutil.ReadDir(root)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.Name() == pattern {
			result = file.Name()
			break
		}
	}
	return result
}

func concat(slices ...[]string) (result []string) {
	for _, slice := range slices {
		for _, s := range slice {
			result = append(result, s)
		}
	}
	return result
}
