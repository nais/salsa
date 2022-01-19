package build_tool

import (
	"fmt"
	"os/exec"

	"github.com/nais/salsa/pkg/scan/jvm"
	"github.com/nais/salsa/pkg/utils"
	"github.com/nais/salsa/pkg/vcs"
	log "github.com/sirupsen/logrus"
)

const gradleBuildFileName = "build.gradle.kts"

type Gradle struct {
	BuildFilePatterns []string
}

func NewGradle() BuildTool {
	return &Gradle{
		BuildFilePatterns: []string{gradleBuildFileName},
	}
}

func (g Gradle) Build(workDir string, project string, context *vcs.AnyContext) error {
	cmd := exec.Command(
		"./gradlew",
		"-q", "dependencies", "--configuration", "runtimeClasspath",
	)
	cmd.Dir = workDir

	depsOutput, err := utils.Exec(cmd)
	if err != nil {
		return fmt.Errorf("exec: %v\n", err)
	}
	log.Info(depsOutput)

	deps, err := jvm.GradleDeps(depsOutput)
	log.Info(workDir)

	err = GenerateProvenance(workDir, project, deps, context)
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
