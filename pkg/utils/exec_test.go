package utils

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

type FakeRunner struct{}

func TestRun(t *testing.T) {
	f := &Cmd{
		Name:    "cosign",
		SubCmd:  "",
		Flags:   []string{"--key", "key", "--predicate", "provenance.json"},
		Args:    []string{"image"},
		WorkDir: "",
		Runner:  &FakeRunner{},
	}
	out, _ := f.Run()

	assert.Equal(t, "[cosign --key key --predicate provenance.json image]", out)
}

func TestDryRunCmd(t *testing.T) {
	if os.Getenv("GO_TEST_DRYRUN") != "1" {
		return
	}
	fmt.Printf("%s", flag.Args())
	os.Exit(0)
}

func (r FakeRunner) CreateCmd() CreateCmd {
	return func(command string, args ...string) *exec.Cmd {
		cs := []string{"-test.run=TestDryRunCmd", "--", command}
		cs = append(cs, args...)
		cmd := exec.Command(os.Args[0], cs...)
		cmd.Env = []string{"GO_TEST_DRYRUN=1"}
		return cmd
	}
}
