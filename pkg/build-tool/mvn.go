package build_tool

import (
	"encoding/json"
	"fmt"
	"github.com/nais/salsa/pkg/exec"
	"github.com/nais/salsa/pkg/intoto"
	"github.com/nais/salsa/pkg/scan/jvm"
	log "github.com/sirupsen/logrus"
	"os"
)

const mavenBuildFileName = "pom.xml"

func NewMaven() BuildTool {
	return &Maven{
		BuildFilePatterns: []string{mavenBuildFileName},
	}
}

type Maven struct {
	BuildFilePatterns []string
}

func (m Maven) Build(workDir, project string) error {
	c := exec.CmdCfg{
		WorkDir: workDir,
		Cmd:     "mvn",
		Args:    []string{"dependency:list"},
	}
	command, err := c.Exec()

	if err != nil {
		return fmt.Errorf("exec: %v\n", err)
	}

	deps, err := jvm.MavenCompileAndRuntimeTimeDeps(command.Output)
	if err != nil {
		return fmt.Errorf("scan: %v\n", err)
	}

	log.Println(deps)

	app := createApp(project, deps)
	s := intoto.GenerateSlsaPredicate(app)

	statement, err := json.Marshal(s)
	if err != nil {
		return fmt.Errorf("marshal: %v\n", err)
	}
	log.Println(string(statement))
	provenanceName := fmt.Sprintf("%s.provenance", project)
	err = os.WriteFile(provenanceName, statement, 0644)
	if err != nil {
		return fmt.Errorf("write to file: %v\n", err)
	}
	return nil
}

func (m Maven) BuildTool(pattern string) bool {
	return Contains(m.BuildFilePatterns, pattern)
}

func (m Maven) BuildFiles() []string {
	return m.BuildFilePatterns
}
