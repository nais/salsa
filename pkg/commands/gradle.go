package commands

import (
	"encoding/json"
	"os"
	"os/exec"

	"github.com/nais/salsa/pkg/intoto"
	"github.com/nais/salsa/pkg/scan"
	log "github.com/sirupsen/logrus"
)

func GradleScan(workDir string) error {
	c := exec.Command(
		"./gradlew",
		"-q", "dependencies", "--configuration", "runtimeClasspath",
	)
	c.Dir = workDir
	result, err := Exec(c)
	if err != nil {
		return err
	}

	gradleDeps, err := scan.GradleDeps(result)
	if err != nil {
		log.Printf("failed: %v\n", err)
		os.Exit(1)
	}

	log.Print(gradleDeps)

	// Should send in Runner/Github info
	app := createApp("tokendings", gradleDeps)
	s := intoto.GenerateSlsaPredicate(app)

	statement, err := json.Marshal(s)
	if err != nil {
		return err
	}
	log.Print(string(statement))
	os.WriteFile("tokendings.provenance", statement, 0644)
	return nil
}

func createApp(name string, deps map[string]string) intoto.App {
	return intoto.App{
		Name:         name,
		BuilderId:    "todoId",
		BuildType:    "todoType",
		Dependencies: deps,
	}
}
