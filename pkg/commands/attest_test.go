package commands

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/nais/salsa/pkg/utils"
	"github.com/stretchr/testify/assert"
)

type fakeRunner struct {
	envVar string
	cmd    string
}

func TestAttestCosignCommand(t *testing.T) {
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

	runner := fakeRunner{
		envVar: "GO_TEST_DRYRUN",
		cmd:    "-test.run=TestDryRunCmd",
	}

	out, err := o.Run([]string{"image"}, runner)
	assert.NoError(t, err)
	expectedCmd := "[cosign attest --key mykey --predicate file.json --type slsaprovenance --rekor-url http://rekor.example.com --no-upload=false image]\n"
	assert.Equal(t, expectedCmd, out)
}

func TestAttestVerifySuccess(t *testing.T) {
	workDir, err := os.MkdirTemp("testdata", "output")
	assert.NoError(t, err)
	defer os.RemoveAll(workDir)

	parts := strings.Split(workDir, "/")
	PathFlags.Repo = parts[1]
	PathFlags.Remote = false
	PathFlags.RepoDir = parts[0]

	o := AttestOptions{
		Key: "mykey",
	}
	verify = true
	runner := fakeRunner{
		envVar: "GO_TEST_ATTEST_VERIFY_OUTPUT",
		cmd:    "-test.run=TestAttestVerifyOutput",
	}

	_, err = o.Run([]string{"image"}, runner)
	assert.NoError(t, err)
}

func TestDryRunCmd(t *testing.T) {
	if os.Getenv("GO_TEST_DRYRUN") != "1" {
		return
	}
	fmt.Printf("%s\n", flag.Args())
	os.Exit(0)
}

func TestAttestVerifyOutput(t *testing.T) {
	if os.Getenv("GO_TEST_ATTEST_VERIFY_OUTPUT") != "1" {
		return
	}
	path := "../cosign-verify-output.txt"
	output, err := os.ReadFile(path)

	if err != nil {
		fmt.Print("fail")
		t.Fatalf("could not read testdata file: %s", err)
	}
	fmt.Printf("%s", output)
	os.Exit(0)
}

func (r fakeRunner) CreateCmd() utils.CreateCmd {
	return func(command string, args ...string) *exec.Cmd {
		cs := []string{r.cmd, "--", command}
		cs = append(cs, args...)
		cmd := exec.Command(os.Args[0], cs...)
		cmd.Env = []string{r.envVar + "=1"}
		return cmd
	}
}

func TestAttestCmd(t *testing.T) {
	PathFlags.Repo = "testdata"
	PathFlags.Remote = false
	PathFlags.RepoDir = "."
	o := AttestOptions{
		Key:           "mykey",
		NoUpload:      true,
		RekorURL:      "http://rekor.example.com",
		PredicateFile: "../testdata/cosgintest.provenance",
		PredicateType: "slsaprovenance",
		AllowInsecure: true,
		SkipVerify:    true,
	}

	args := []string{
		"ttl.image:1h",
	}
	_, err := o.Run(args, utils.ExecCmd{})
	assert.NoError(t, err)
}
