package utils

import (
	"bytes"
	"fmt"
	"github.com/briandowns/spinner"
	"io"
	"os"
	"os/exec"
	"time"
)

func RequireCommand(cmd string) error {
	if _, err := exec.LookPath(cmd); err != nil {
		return fmt.Errorf("could not find required cmd: %w", err)
	}
	return nil
}

func Exec(cmd *exec.Cmd) (string, error) {
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("%s failed:\n%w\n", cmd, err)
	}
	outStr, _ := string(stdoutBuf.Bytes()), string(stderrBuf.Bytes())
	return outStr, nil
}

func StartSpinner(message string) *spinner.Spinner {
	s := spinner.New(spinner.CharSets[11], 150*time.Millisecond)
	s.FinalMSG = message
	s.Start()
	return s
}
