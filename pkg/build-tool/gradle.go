package build_tool

import (
	"encoding/json"
	"fmt"
	"github.com/nais/salsa/pkg/exec"
	"github.com/nais/salsa/pkg/intoto"
	"github.com/nais/salsa/pkg/scan"
	log "github.com/sirupsen/logrus"
	"os"
)

func NewGradle() BuildTool {
	return &Gradle{
		BuildFilePatterns: []string{"build.gradle.kts"},
	}
}

type Gradle struct {
	BuildFilePatterns []string
}

func (g Gradle) Build(workDir, project string) error {
	c := exec.CmdCfg{
		WorkDir: workDir,
		Cmd:     "./gradlew",
		Args:    []string{"-q", "dependencies", "--configuration", "runtimeClasspath"},
	}
	command, err := c.Exec()

	if err != nil {
		return fmt.Errorf("exec: %v\n", err)
	}

	gradleDeps, err := scan.GradleDeps(command.Output)
	if err != nil {
		return fmt.Errorf("scan: %v\n", err)
	}

	log.Print(gradleDeps)

	// Should send in Runner/Github info
	app := createApp(project, gradleDeps)
	s := intoto.GenerateSlsaPredicate(app)

	statement, err := json.Marshal(s)
	if err != nil {
		return fmt.Errorf("marshal: %v\n", err)
	}
	log.Print(string(statement))
	provenanceName := fmt.Sprintf("%s.provenance", project)
	err = os.WriteFile(provenanceName, statement, 0644)
	if err != nil {
		return fmt.Errorf("write to file: %v\n", err)
	}
	return nil
}

func (g Gradle) BuildTool(pattern string) bool {
	return contains(g.BuildFilePatterns, pattern)
}
