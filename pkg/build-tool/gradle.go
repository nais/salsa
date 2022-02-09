package build_tool

import (
	"fmt"
	"io/ioutil"
	"os/exec"

	"github.com/nais/salsa/pkg/scan"
	"github.com/nais/salsa/pkg/scan/jvm"
	"github.com/nais/salsa/pkg/utils"
)

const gradleBuildFileName = "build.gradle.kts"

type Gradle struct {
	BuildFilePatterns []string
}

func (g Gradle) ResolveDeps(workDir string) (*scan.ArtifactDependencies, error) {
	cmd := exec.Command(
		"gradle",
		"-q", "dependencies", "--configuration", "runtimeClasspath", "-M", "sha256",
	)
	cmd.Dir = workDir

	err := utils.RequireCommand("gradle")
	if err != nil {
		return nil, fmt.Errorf("exec: %v\n", err)
	}

	depsOutput, err := utils.Exec(cmd)
	if err != nil {
		return nil, fmt.Errorf("exec: %v\n", err)
	}

	xml, err := ioutil.ReadFile(workDir + "/gradle/verification-metadata.xml")
	if err != nil {
		return nil, fmt.Errorf("readfile: %v\n", err)
	}

	deps, err := jvm.GradleDeps(depsOutput, xml)
	if err != nil {
		return nil, fmt.Errorf("could not get gradle deps: %w", err)
	}

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
