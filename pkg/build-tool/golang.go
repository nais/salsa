package build_tool

import (
	"fmt"
	"os"

	"github.com/nais/salsa/pkg/scan"
	"github.com/nais/salsa/pkg/scan/golang"
)

const golangBuildFileName = "go.sum"

type Golang struct {
	BuildFilePatterns []string
}

func (g Golang) ResolveDeps(workDir string) (*scan.ArtifactDependencies, error) {
	fileContent, err := os.ReadFile(fmt.Sprintf("%s/%s", workDir, golangBuildFileName))
	deps := golang.GoDeps(string(fileContent))
	if err != nil {
		return nil, fmt.Errorf("error parsing %s, %v", golangBuildFileName, err)
	}
	return &scan.ArtifactDependencies{
		Cmd:         golangBuildFileName,
		RuntimeDeps: deps,
	}, nil
}

func NewGolang() BuildTool {
	return &Golang{
		BuildFilePatterns: []string{golangBuildFileName},
	}
}

func (g Golang) Supported(pattern string) bool {
	return Contains(g.BuildFilePatterns, pattern)
}

func (g Golang) BuildFiles() []string {
	return g.BuildFilePatterns
}
