package build_tool

import (
	"fmt"
	"os/exec"

	"github.com/nais/salsa/pkg/scan"
	"github.com/nais/salsa/pkg/scan/jvm"
	"github.com/nais/salsa/pkg/utils"
)

const mavenBuildFileName = "pom.xml"

type Maven struct {
	BuildFilePatterns []string
}

func (m Maven) ResolveDeps(workDir string) (*scan.ArtifactDependencies, error) {
	cmd := exec.Command(
		"mvn",
		"dependency:list",
	)
	cmd.Dir = workDir

	output, err := utils.Exec(cmd)
	if err != nil {
		return nil, fmt.Errorf("exec: %v\n", err)
	}

	deps, err := jvm.MavenCompileAndRuntimeTimeDeps(output)
	if err != nil {
		return nil, fmt.Errorf("scan: %v\n", err)
	}
	return &scan.ArtifactDependencies{
		Cmd:         fmt.Sprintf("%s %v", cmd.Path, cmd.Args),
		RuntimeDeps: deps,
	}, nil
}

func NewMaven() BuildTool {
	return &Maven{
		BuildFilePatterns: []string{mavenBuildFileName},
	}
}

func (m Maven) Supported(pattern string) bool {
	return Contains(m.BuildFilePatterns, pattern)
}

func (m Maven) BuildFiles() []string {
	return m.BuildFilePatterns
}
