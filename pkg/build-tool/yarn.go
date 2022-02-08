package build_tool

import (
	"fmt"
	"os"

	"github.com/nais/salsa/pkg/scan"
	"github.com/nais/salsa/pkg/scan/nodejs"
)

const yarnBuildFileName = "yarn.lock"

type Yarn struct {
	BuildFilePatterns []string
}

func NewYarn() BuildTool {
	return &Yarn{
		BuildFilePatterns: []string{yarnBuildFileName},
	}
}

func (y Yarn) ResolveDeps(workDir string) (*scan.ArtifactDependencies, error) {
	fileContent, err := os.ReadFile(fmt.Sprintf("%s/%s", workDir, yarnBuildFileName))
	if err != nil {
		return nil, fmt.Errorf("read file: %w\n", err)
	}
	deps := nodejs.YarnDeps(string(fileContent))

	return &scan.ArtifactDependencies{
		Cmd:         yarnBuildFileName,
		RuntimeDeps: deps,
	}, nil
}

func (y Yarn) Supported(pattern string) bool {
	return Contains(y.BuildFilePatterns, pattern)
}

func (y Yarn) BuildFiles() []string {
	return y.BuildFilePatterns
}
