package build_tool

import (
	"fmt"
	"os"

	"github.com/nais/salsa/pkg/scan"
	"github.com/nais/salsa/pkg/scan/nodejs"
)

const npmBuildFileName = "package-lock.json"

type Npm struct {
	BuildFilePatterns []string
}

func (n Npm) ResolveDeps(workDir string) (*scan.ArtifactDependencies, error) {
	fileContent, err := os.ReadFile(fmt.Sprintf("%s/%s", workDir, npmBuildFileName))
	if err != nil {
		return nil, fmt.Errorf("read file: %w\n", err)
	}
	deps, err := nodejs.NpmDeps(string(fileContent))
	if err != nil {
		return nil, fmt.Errorf("error parsing deps: %v\n", err)
	}
	return &scan.ArtifactDependencies{
		Cmd:         npmBuildFileName,
		RuntimeDeps: deps,
	}, nil
}

func NewNpm() BuildTool {
	return &Npm{
		BuildFilePatterns: []string{npmBuildFileName},
	}
}

func (n Npm) Supported(pattern string) bool {
	return Contains(n.BuildFilePatterns, pattern)
}

func (n Npm) BuildFiles() []string {
	return n.BuildFilePatterns
}
