package build_tool

import (
	"fmt"
	"github.com/nais/salsa/pkg/exec"
	"github.com/nais/salsa/pkg/scan/golang"
	log "github.com/sirupsen/logrus"
)

const golangBuildFileName = "go.sum"

type Golang struct {
	BuildFilePatterns []string
	Cmd               exec.CmdCfg
}

func NewGolang(workDir string) BuildTool {
	return &Golang{
		BuildFilePatterns: []string{golangBuildFileName},
		Cmd: exec.CmdCfg{
			WorkDir: workDir,
			Cmd:     "cat",
			Args:    []string{"go.sum"},
		},
	}
}

func (g Golang) Build(project string) error {
	command, err := g.Cmd.Exec()

	if err != nil {
		return fmt.Errorf("exec: %v\n", err)
	}

	deps := golang.GoDeps(command.Output)
	log.Println(deps)

	err = GenerateProvenance(project, deps)
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
