package utils

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
)

type CreateCmd = func(name string, arg ...string) *exec.Cmd

type Cmd struct {
	Name    string
	SubCmd  string
	Flags   []string
	Args    []string
	WorkDir string
	Runner  CmdRunner
}

func NewCmd(
	name string,
	subCmd string,
	flags []string,
	args []string,
	workDir string,
) Cmd {
	return Cmd{Name: name, SubCmd: subCmd, Flags: flags, Args: args, WorkDir: workDir, Runner: &ExecCmd{}}
}

func (c *Cmd) WithRunner(runner CmdRunner) {
	c.Runner = runner
}

type CmdRunner interface {
	CreateCmd() CreateCmd
}

type ExecCmd struct{}

func (c ExecCmd) CreateCmd() CreateCmd {
	return exec.Command
}

func (c Cmd) Run() (string, error) {
	args := make([]string, 0)
	if c.SubCmd != "" {
		args = append(args, c.SubCmd)
	}
	if c.Flags != nil {
		args = append(args, c.Flags...)
	}
	if c.Args != nil {
		args = append(args, c.Args...)
	}
	cmd := c.Runner.CreateCmd()(c.Name, args...)

	err := requireCommand(cmd.Path)
	if err != nil {
		return "", err
	}
	cmd.Dir = c.WorkDir
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)

	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("%s failed:\n%w\n", cmd, err)
	}
	outStr, _ := stdoutBuf.String(), stderrBuf.String()
	return outStr, nil
}

func requireCommand(cmd string) error {
	if _, err := exec.LookPath(cmd); err != nil {
		return fmt.Errorf("could not find required cmd: %w", err)
	}
	return nil
}
