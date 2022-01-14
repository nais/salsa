package build_tool

import (
	"fmt"
	"github.com/nais/salsa/pkg/scan/nodejs"
	"github.com/nais/salsa/pkg/vcs"
	log "github.com/sirupsen/logrus"
	"os"
)

const npmBuildFileName = "package-lock.json"

type Npm struct {
	BuildFilePatterns []string
}

func NewNpm() BuildTool {
	return &Npm{
		BuildFilePatterns: []string{npmBuildFileName},
	}
}

func (m Npm) Build(workDir, project string, context *vcs.AnyContext) error {
	fileContent, err := os.ReadFile(fmt.Sprintf("%s/%s", workDir, npmBuildFileName))
	if err != nil {
		return fmt.Errorf("read file: %w\n", err)
	}
	deps, err := nodejs.NpmDeps(string(fileContent))
	if err != nil {
		return fmt.Errorf("scan: %v\n", err)
	}

	log.Println(deps)

	err = GenerateProvenance(workDir, project, deps, context)
	if err != nil {
		return fmt.Errorf("generating provencance %v", err)
	}
	return nil
}

func (m Npm) BuildTool(pattern string) bool {
	return Contains(m.BuildFilePatterns, pattern)
}

func (m Npm) BuildFiles() []string {
	return m.BuildFilePatterns
}
