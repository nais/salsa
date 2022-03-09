package golang

import (
	"fmt"
	"os"
	"strings"

	"github.com/nais/salsa/pkg/build"
	"github.com/nais/salsa/pkg/digest"
)

const golangBuildFileName = "go.sum"

type Golang struct {
	BuildFilePatterns []string
}

func NewGolang() build.Tool {
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
	deps := GoDeps(string(fileContent))
	if err != nil {
		return nil, fmt.Errorf("error parsing %s, %v", golangBuildFileName, err)
	}
	return build.ArtifactDependency(deps, path, golangBuildFileName), nil
}

func GoDeps(goSumContents string) map[string]build.Dependency {
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
		checksum := build.Verification(digest.AlgorithmSHA256, base64EncodedDigest)
		deps[coordinates] = build.Dependence(coordinates, version, checksum)
	}
	return deps
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
