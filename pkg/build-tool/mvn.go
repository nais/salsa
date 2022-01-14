package build_tool

import (
	"fmt"
	"github.com/nais/salsa/pkg/vcs"
	"os/exec"

	"github.com/nais/salsa/pkg/scan/jvm"
	"github.com/nais/salsa/pkg/utils"
	log "github.com/sirupsen/logrus"
)

const mavenBuildFileName = "pom.xml"

type Maven struct {
	BuildFilePatterns []string
}

func NewMaven() BuildTool {
	return &Maven{
		BuildFilePatterns: []string{mavenBuildFileName},
	}
}

func (m Maven) Build(workDir, project string, context *vcs.AnyContext) error {

	cmd := exec.Command(
		"mvn",
		"dependency:list",
	)
	cmd.Dir = workDir

	output, err := utils.Exec(cmd)
	if err != nil {
		return fmt.Errorf("exec: %v\n", err)
	}

	deps, err := jvm.MavenCompileAndRuntimeTimeDeps(output)
	if err != nil {
		return fmt.Errorf("scan: %v\n", err)
	}

	log.Println(deps)

	err = GenerateProvenance(workDir, project, deps, context)
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
