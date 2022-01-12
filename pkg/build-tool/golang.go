package build_tool

import (
	"fmt"
	"os/exec"

	"github.com/nais/salsa/pkg/scan/golang"
	"github.com/nais/salsa/pkg/utils"
	log "github.com/sirupsen/logrus"
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

func (g Golang) Build(workDir, project string) error {
	cmd := exec.Command(
		"cat",
		"go.sum",
	)
	cmd.Dir = workDir

	output, err := utils.Exec(cmd)

	if err != nil {
		return fmt.Errorf("exec: %v\n", err)
	}

	deps := golang.GoDeps(output)
	log.Println(deps)

	err = GenerateProvenance(workDir, project, deps)
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
