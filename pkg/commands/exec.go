package commands

import (
    "bytes"
    "fmt"
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

    log.Printf("cmd: %s", cmd)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Dir = c.workDir
	err := cmd.Run()

    outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	if err != nil {
		log.Printf("cmd.Run: %s failed: %v\n", cmd, err)
        log.Printf("stderr: %s", errStr)
        return CommandOutput{}, err
	}
	if len(errStr) > 1 {
		fmt.Printf("out:\n%s\nerr:\n%s\n", outStr, errStr)
	}
	fmt.Printf(outStr)
	return CommandOutput{Output: outStr}, nil
}
