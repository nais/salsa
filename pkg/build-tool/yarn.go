package build_tool

import (
	"fmt"
	"github.com/nais/salsa/pkg/scan/nodejs"
	"github.com/nais/salsa/pkg/vcs"
	"os"
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

func (m Yarn) Build(workDir, project string, context *vcs.AnyContext) error {
	fileContent, err := os.ReadFile(fmt.Sprintf("%s/%s", workDir, yarnBuildFileName))
	if err != nil {
		return fmt.Errorf("read file: %w\n", err)
	}
	deps := nodejs.YarnDeps(string(fileContent))
	err = GenerateProvenance(workDir, project, deps, context)
	if err != nil {
		return fmt.Errorf("generating provencance %v", err)
	}
	return nil
}

func (m Yarn) BuildTool(pattern string) bool {
	return Contains(m.BuildFilePatterns, pattern)
}

func (m Yarn) BuildFiles() []string {
	return m.BuildFilePatterns
}
