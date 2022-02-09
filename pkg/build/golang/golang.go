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

func (g Golang) ResolveDeps(workDir string) (*build.ArtifactDependencies, error) {
	fileContent, err := os.ReadFile(fmt.Sprintf("%s/%s", workDir, golangBuildFileName))
	deps := GoDeps(string(fileContent))
	if err != nil {
		return nil, fmt.Errorf("error parsing %s, %v", golangBuildFileName, err)
	}
	return &build.ArtifactDependencies{
		Cmd:         golangBuildFileName,
		RuntimeDeps: deps,
	}, nil
}

func NewGolang() build.BuildTool {
	return &Golang{
		BuildFilePatterns: []string{golangBuildFileName},
	}
}

func (g Golang) BuildFiles() []string {
	return g.BuildFilePatterns
}

func GoDeps(goSumContents string) []build.Dependency {
	deps := make([]build.Dependency, 0)
	lines := strings.Split(goSumContents, "\n")
	for _, line := range lines {
		if isNotInteresting(line) {
			continue
		}
		parts := strings.Split(line, " ")
		version := parts[1][1:]
		stringDecodedDigest := digest.Digest(strings.Split(parts[2], ":")[1])
		deps = append(deps, build.Dependency{
			Coordinates: parts[0],
			Version:     version,
			CheckSum:    build.CheckSum{Algorithm: digest.SHA256, Digest: string(stringDecodedDigest)},
		})
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
