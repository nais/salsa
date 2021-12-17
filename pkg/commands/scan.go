package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/nais/salsa/pkg/intoto"
	"github.com/nais/salsa/pkg/scan"
)

type CmdConfig struct {
	workDir string
	cmd     string
	args    []string
}

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

	m, err := json.Marshal(s)
	if err != nil {
		fmt.Printf("failed: %v\n", err)
		os.Exit(1)
	}
	log.Print(string(m))
}

func createApp(name string, deps map[string]string) intoto.App {
	return intoto.App{
		Name:         name,
		BuilderId:    "todoId",
		BuildType:    "todoType",
		Dependencies: deps,
	}
}

type CommandOutput struct {
	Output string
}

func (c CmdConfig) ExecuteCommand() (CommandOutput, error) {
	cmd := exec.Command(c.cmd, c.args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Dir = c.workDir
	err := cmd.Run()

	if err != nil {
		log.Printf("cmd.Run: %s failed: %v\n", cmd, err)
		os.Exit(1)
	}
	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	if len(errStr) > 1 {
		fmt.Printf("out:\n%s\nerr:\n%s\n", outStr, errStr)
	}
	fmt.Printf(outStr)
	return CommandOutput{Output: outStr}, nil
}
