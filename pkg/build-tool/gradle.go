package build_tool

import (
	"fmt"
	"github.com/nais/salsa/pkg/exec"
	"github.com/nais/salsa/pkg/scan/jvm"
	log "github.com/sirupsen/logrus"
)

const gradleBuildFileName = "build.gradle.kts"

type Gradle struct {
	BuildFilePatterns []string
	Cmd               exec.CmdCfg
}

func NewGradle(workDir string) BuildTool {
	return &Gradle{
		BuildFilePatterns: []string{gradleBuildFileName},
		Cmd: exec.CmdCfg{
			WorkDir: workDir,
			Cmd:     "./gradlew",
			Args:    []string{"-q", "dependencies", "--configuration", "runtimeClasspath"},
		},
	}
}

func (g Gradle) Build(project string) error {
	command, err := g.Cmd.Exec()
	if err != nil {
		return fmt.Errorf("exec: %v\n", err)
	}

	deps, err := jvm.GradleDeps(command.Output)
	if err != nil {
		return fmt.Errorf("scan: %v\n", err)
	}

	log.Print(deps)

	err = GenerateProvenance(project, deps)
	if err != nil {
		return fmt.Errorf("generating provencance %v", err)
	}
	return nil
}

func (g Gradle) BuildTool(pattern string) bool {
	return Contains(g.BuildFilePatterns, pattern)
}

func (g Gradle) BuildFiles() []string {
	return g.BuildFilePatterns
}
