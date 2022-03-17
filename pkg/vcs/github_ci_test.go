package vcs

import (
	"encoding/base64"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestCreateCIEnvironment(t *testing.T) {
	err := os.Setenv("GITHUB_ACTIONS", "true")
	assert.NoError(t, err)
	context := githubContext(t)
	runner := runnerContext()
	env := envC()
	ci, err := CreateGithubCIEnvironment(context, &runner, &env)
	assert.NoError(t, err)
	assert.Equal(t, "https://github.com/nais/salsa", ci.RepoUri())
	assert.Equal(t, "90dc9f2bc4007d1099a941ba3d408d2c896fe8dd", ci.Sha())
	assert.Equal(t, "https://github.com/nais/salsa/Attestations/GitHubHostedActions@v1", ci.BuilderId())
	assert.Equal(t, "https://github.com/nais/salsa/actions/runs/1839977840", ci.BuildInvocationId())

	metadata := Metadata{
		Arch: "X64",
		Env: map[string]string{
			"GO_ROOT":    "/opt/hostedtoolcache/go/1.17.6/x64",
			"GO_VERSION": "1.17",
		},
		Context: Context{
			Github: Github{
				RunId: "1839977840",
			}, Runner: Runner{
				Os:   "Linux",
				Temp: "/home/runner/work/_temp",
			},
		},
	}
	assert.Equal(t, &metadata, ci.NonReproducibleMetadata())

	current := map[string]string{
		"GO_ROOT":    "/opt/hostedtoolcache/go/1.17.6/x64",
		"GO_VERSION": "1.17",
	}
	assert.Equal(t, current, ci.CurrentFilteredEnvironment())

	result, err := ci.UserDefinedParameters().Inputs.MarshalJSON()
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, fmt.Sprintf("%s", result))
	assert.Equal(t, "Salsa CI", ci.Context())

}

func githubContext(t *testing.T) []byte {
	githubContext, err := os.ReadFile("testdata/github-context.json")
	assert.NoError(t, err)
	return githubContext
}

func runnerContext() string {
	return base64.StdEncoding.EncodeToString([]byte(RunnerTestContext))
}

func envC() string {
	return base64.StdEncoding.EncodeToString([]byte(envTestContext))
}

var envTestContext = `{
  		"GO_VERSION": "1.17",
		"GO_ROOT": "/opt/hostedtoolcache/go/1.17.6/x64"
	  }`

var RunnerTestContext = `{
		"os": "Linux",
		"arch": "X64",
		"name": "Hosted Agent",
		"tool_cache": "/opt/hostedtoolcache",
		"temp": "/home/runner/work/_temp",
		"workspace": "/home/runner/work/nais-salsa-action"
	  }`
