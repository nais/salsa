package commands

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

type CmdConfig struct {
	workDir string
	cmd     string
	args    []string
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
