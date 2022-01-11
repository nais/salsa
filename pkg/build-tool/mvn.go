package build_tool

import (
	"fmt"
	"github.com/nais/salsa/pkg/exec"
	"github.com/nais/salsa/pkg/scan/jvm"
	log "github.com/sirupsen/logrus"
)

const mavenBuildFileName = "pom.xml"

type Maven struct {
	BuildFilePatterns []string
	Cmd               exec.CmdCfg
}

func NewMaven(workDir string) BuildTool {
	return &Maven{
		BuildFilePatterns: []string{mavenBuildFileName},
		Cmd: exec.CmdCfg{
			WorkDir: workDir,
			Cmd:     "mvn",
			Args:    []string{"dependency:list"},
		},
	}
}

func (m Maven) Build(project string) error {
	command, err := m.Cmd.Exec()
	if err != nil {
		return fmt.Errorf("exec: %v\n", err)
	}

	deps, err := jvm.MavenCompileAndRuntimeTimeDeps(command.Output)
	if err != nil {
		return fmt.Errorf("scan: %v\n", err)
	}

	log.Println(deps)

	err = GenerateProvenance(project, deps)
	if err != nil {
		return fmt.Errorf("generating provencance %v", err)
	}
	return nil
}

func (m Maven) BuildTool(pattern string) bool {
	return Contains(m.BuildFilePatterns, pattern)
}

func (m Maven) BuildFiles() []string {
	return m.BuildFilePatterns
}
