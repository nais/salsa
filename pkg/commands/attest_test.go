package commands

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/nais/salsa/pkg/utils"
	"github.com/stretchr/testify/assert"
)

type FakeRunner struct{}

func TestAttest(t *testing.T) {
	PathFlags.Repo = "."
	PathFlags.Remote = false
	PathFlags.RepoDir = "."

	o := AttestOptions{
		Key:           "mykey",
		NoUpload:      false,
		RekorURL:      "http://rekor.example.com",
		PredicateFile: "file.json",
		PredicateType: "slsaprovenance",
	}
	out, err := o.Run([]string{"image"}, FakeRunner{})
	assert.NoError(t, err)

	expected := "[cosign attest --key mykey --predicate file.json --type slsaprovenance --rekor-url http://rekor.example.com --no-upload=false image]"
	assert.Equal(t, expected, out)
}

func TestDryRunCmd(t *testing.T) {
	if os.Getenv("GO_TEST_DRYRUN") != "1" {
		return
	}
	fmt.Printf("%s", flag.Args())
	os.Exit(0)
}

func (r FakeRunner) CreateCmd() utils.CreateCmd {
	return func(command string, args ...string) *exec.Cmd {
		cs := []string{"-test.run=TestDryRunCmd", "--", command}
		cs = append(cs, args...)
		cmd := exec.Command(os.Args[0], cs...)
		cmd.Env = []string{"GO_TEST_DRYRUN=1"}
		return cmd
	}
}
