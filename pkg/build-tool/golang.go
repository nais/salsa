package build_tool

import (
	"fmt"
	"github.com/nais/salsa/pkg/scan/golang"
	"github.com/nais/salsa/pkg/vcs"
	log "github.com/sirupsen/logrus"
	"os"
)

const golangBuildFileName = "go.sum"

type Golang struct {
	BuildFilePatterns []string
}

func NewGolang() BuildTool {
	return &Golang{
		BuildFilePatterns: []string{golangBuildFileName},
	}
}

func (g Golang) Build(workDir, project string, context *vcs.AnyContext) error {
	fileContent, err := os.ReadFile(fmt.Sprintf("%s/%s", workDir, golangBuildFileName))
	goMetadata := golang.GoDeps(string(fileContent))
	log.Println(goMetadata.Deps)

	err = GenerateProvenance(workDir, project, goMetadata, context)
	if err != nil {
		return fmt.Errorf("generating provencance %v", err)
	}
	return nil
}

func (g Golang) BuildTool(pattern string) bool {
	return Contains(g.BuildFilePatterns, pattern)
}

func (g Golang) BuildFiles() []string {
	return g.BuildFilePatterns
}
