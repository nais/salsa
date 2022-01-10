package exec

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
)

type CmdCfg struct {
	WorkDir string
	Cmd     string
	Args    []string
}

type CommandOutput struct {
	Output string
	Error  string
}

func (cfg CmdCfg) Exec() (*CommandOutput, error) {
	cmd := exec.Command(cfg.Cmd, cfg.Args...)
	cmd.Dir = cfg.WorkDir

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("%s failed:\n%w\n", cmd, err)
	}
	outStr, errStr := string(stdoutBuf.Bytes()), string(stderrBuf.Bytes())
	return &CommandOutput{
		Output: outStr,
		Error:  errStr,
	}, nil
}
