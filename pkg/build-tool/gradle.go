package build_tool

import (
	"fmt"
	"os/exec"

	"github.com/nais/salsa/pkg/scan"
	"github.com/nais/salsa/pkg/scan/jvm"
	"github.com/nais/salsa/pkg/utils"
	log "github.com/sirupsen/logrus"
)

const gradleBuildFileName = "build.gradle.kts"

type Gradle struct {
	BuildFilePatterns []string
}

func (g Gradle) ResolveDeps(workDir string) (*scan.ArtifactDependencies, error) {
	cmd := exec.Command(
		"./gradlew",
		"-q", "dependencies", "--configuration", "runtimeClasspath",
	)
	cmd.Dir = workDir

	depsOutput, err := utils.Exec(cmd)
	if err != nil {
		return nil, fmt.Errorf("exec: %v\n", err)
	}
	log.Info(depsOutput)

	deps, err := jvm.GradleDeps(depsOutput)
	if err != nil {
		return nil, fmt.Errorf("could not get gradle deps: %w", err)
	}
	log.Info(workDir)

	return &scan.ArtifactDependencies{
		Cmd:         fmt.Sprintf("%s %v", cmd.Path, cmd.Args),
		RuntimeDeps: deps,
	}, nil
}

func NewGradle() BuildTool {
	return &Gradle{
		BuildFilePatterns: []string{gradleBuildFileName},
	}
}

func (g Gradle) Supported(pattern string) bool {
	return Contains(g.BuildFilePatterns, pattern)
}

func (g Gradle) BuildFiles() []string {
	return g.BuildFilePatterns
}
