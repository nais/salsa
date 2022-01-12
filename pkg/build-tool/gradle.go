package build_tool

import (
	"fmt"
	"os/exec"

	"github.com/nais/salsa/pkg/scan/jvm"
	"github.com/nais/salsa/pkg/utils"
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

func (g Gradle) Build(workDir string, project string) error {
	cmd := exec.Command(
		"./gradlew",
		"-q", "dependencies", "--configuration", "runtimeClasspath",
	)
	cmd.Dir = workDir

	output, err := utils.Exec(cmd)
	if err != nil {
		return fmt.Errorf("exec: %v\n", err)
	}

	deps, err := jvm.GradleDeps(output)
	if err != nil {
		return fmt.Errorf("scan: %v\n", err)
	}

	log.Print(deps)

	err = GenerateProvenance(workDir, project, deps)
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
