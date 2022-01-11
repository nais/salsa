package build_tool

import (
	"encoding/json"
	"fmt"
	"github.com/nais/salsa/pkg/exec"
	"github.com/nais/salsa/pkg/intoto"
	"github.com/nais/salsa/pkg/scan/golang"
	log "github.com/sirupsen/logrus"
	"os"
)

const golangBuildFileName = "go.sum"

func NewGolang() BuildTool {
	return &Golang{
		BuildFilePatterns: []string{golangBuildFileName},
	}
}

type Golang struct {
	BuildFilePatterns []string
}

func (g Golang) Build(workDir, project string) error {
	c := exec.CmdCfg{
		WorkDir: workDir,
		Cmd:     "cat",
		Args:    []string{"go.sum"},
	}
	command, err := c.Exec()

	if err != nil {
		return fmt.Errorf("exec: %v\n", err)
	}

	deps := golang.GoDeps(command.Output)
	log.Println(deps)

	app := CreateApp(project, deps)
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

func (g Golang) BuildTool(pattern string) bool {
	return Contains(g.BuildFilePatterns, pattern)
}

func (g Golang) BuildFiles() []string {
	return g.BuildFilePatterns
}
