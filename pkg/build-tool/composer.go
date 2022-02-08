package build_tool

import (
	"fmt"
	"os"

	"github.com/nais/salsa/pkg/scan"
	"github.com/nais/salsa/pkg/scan/php"
)

const composerLockFileName = "composer.lock"

type Composer struct {
	BuildFilePatterns []string
}

func (c Composer) ResolveDeps(workDir string) (*scan.ArtifactDependencies, error) {
	fileContent, err := os.ReadFile(fmt.Sprintf("%s/%s", workDir, composerLockFileName))
	if err != nil {
		return nil, fmt.Errorf("read file: %w\n", err)
	}
	deps, err := php.ComposerDeps(string(fileContent))
	if err != nil {
		return nil, fmt.Errorf("scan: %v\n", err)
	}
	return &scan.ArtifactDependencies{
		Cmd:         composerLockFileName,
		RuntimeDeps: deps,
	}, nil
}

func NewComposer() BuildTool {
	return &Composer{
		BuildFilePatterns: []string{composerLockFileName},
	}
}

func (c Composer) Supported(pattern string) bool {
	return Contains(c.BuildFilePatterns, pattern)
}

func (c Composer) BuildFiles() []string {
	return c.BuildFilePatterns
}
