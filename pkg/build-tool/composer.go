package build_tool

import (
	"fmt"
	"github.com/nais/salsa/pkg/scan/php"
	"github.com/nais/salsa/pkg/vcs"
	log "github.com/sirupsen/logrus"
	"os"
)

const composerLockFileName = "composer.lock"

type Composer struct {
	BuildFilePatterns []string
}

func NewComposer() BuildTool {
	return &Composer{
		BuildFilePatterns: []string{composerLockFileName},
	}
}

func (m Composer) Build(workDir, project string, context *vcs.AnyContext) error {
	fileContent, err := os.ReadFile(fmt.Sprintf("%s/%s", workDir, composerLockFileName))
	if err != nil {
		return fmt.Errorf("read file: %w\n", err)
	}
	deps, err := php.ComposerDeps(string(fileContent))
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

func (m Composer) BuildTool(pattern string) bool {
	return Contains(m.BuildFilePatterns, pattern)
}

func (m Composer) BuildFiles() []string {
	return m.BuildFilePatterns
}
