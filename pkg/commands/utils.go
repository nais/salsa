package commands

import (
	"fmt"
	"os/exec"
)

func RequireCommand(cmd string) error {
	if _, err := exec.LookPath(cmd); err != nil {
		return fmt.Errorf("could not find required cmd: %w", err)
	}
	return nil
}
