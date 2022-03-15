package golang

import (
	"fmt"
	"github.com/nais/salsa/pkg/build"
	"github.com/nais/salsa/pkg/utils"
	"os"
	"strings"
)

const golangBuildFileName = "go.sum"

type Golang struct {
	BuildFilePatterns []string
}

func BuildGo() build.Tool {
	return &Golang{
		BuildFilePatterns: []string{golangBuildFileName},
	}
}

func (g Golang) BuildFiles() []string {
	return g.BuildFilePatterns
}

func (g Golang) ResolveDeps(workDir string) (*build.ArtifactDependencies, error) {
	path := fmt.Sprintf("%s/%s", workDir, golangBuildFileName)

	fileContent, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("readfile %v", err)
	}

	deps, err := GoDeps(string(fileContent))
	if err != nil {
		return nil, fmt.Errorf("error parsing %s, %v", golangBuildFileName, err)
	}

	return build.ArtifactDependency(deps, path, golangBuildFileName), nil
}

func GoDeps(goSumContents string) (map[string]build.Dependency, error) {
	deps := make(map[string]build.Dependency, 0)
	lines := strings.Split(goSumContents, "\n")
	for _, line := range lines {
		if isNotInteresting(line) {
			continue
		}
		parts := strings.Split(line, " ")
		version := parts[1][1:]
		coordinates := parts[0]
		base64EncodedDigest := strings.Split(parts[2], ":")[1]
		digest, err := utils.DecodeDigest(base64EncodedDigest)
		if err != nil {
			return nil, err
		}
		checksum := build.Verification(build.AlgorithmSHA256, digest)
		deps[coordinates] = build.Dependence(coordinates, version, checksum)
	}
	return deps, nil
}

func isNotInteresting(line string) bool {
	return isEmpty(line) || isMod(line)
}

func isEmpty(line string) bool {
	return strings.TrimSpace(line) == ""
}

func isMod(line string) bool {
	idx := strings.Index(line, "go.mod")
	return idx > -1
}
