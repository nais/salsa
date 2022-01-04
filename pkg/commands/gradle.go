package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/nais/salsa/pkg/intoto"
	"github.com/nais/salsa/pkg/scan"
	log "github.com/sirupsen/logrus"
)

func GradleScan(workDir string) {
	c := CmdConfig{
		workDir: workDir,
		cmd:     "./gradlew",
		args:    []string{"-q", "dependencies", "--configuration", "runtimeClasspath"},
	}
	command, err := c.ExecuteCommand()

	if err != nil {
		log.Printf("failed: %v\n", err)
		os.Exit(1)
	}

	gradleDeps, err := scan.GradleDeps(command.Output)
	if err != nil {
		log.Printf("failed: %v\n", err)
		os.Exit(1)
	}

	log.Print(gradleDeps)

	// Should send in Runner/Github info
	app := createApp("tokendings", gradleDeps)
	s := intoto.GenerateStatement(app)

	statement, err := json.Marshal(s)
	if err != nil {
		fmt.Printf("failed: %v\n", err)
		os.Exit(1)
	}
	log.Print(string(statement))
	os.WriteFile("tokendings.provenance", statement, 0644)
}

func createApp(name string, deps map[string]string) intoto.App {
	return intoto.App{
		Name:         name,
		BuilderId:    "todoId",
		BuildType:    "todoType",
		Dependencies: deps,
	}
}
