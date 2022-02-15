package vcs

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestCreateGithubContext(t *testing.T) {
	githubContext, err := os.ReadFile("testdata/github-context.json")
	assert.NoError(t, err)
	parsedContext := string(githubContext)
	env := Environment{}
	err = ParseGithub(&parsedContext, &env)
	assert.NoError(t, err)

	assert.Equal(t, "90dc9f2bc4007d1099a941ba3d408d2c896fe8dd", env.GitHubContext.SHA)
	assert.Equal(t, "build", env.GitHubContext.Job)
	assert.Equal(t, "refs/heads/main", env.GitHubContext.Ref)
	assert.Equal(t, "nais/salsa", env.GitHubContext.Repository)
	assert.Equal(t, "nais", env.GitHubContext.RepositoryOwner)
	assert.Equal(t, "1839977840", env.GitHubContext.RunId)
	assert.Equal(t, "57", env.GitHubContext.RunNumber)
	assert.Equal(t, "jdoe", env.GitHubContext.Actor)

	assert.Equal(t, "https://github.com/nais/salsa", env.RepoUri())
	assert.Equal(t, "https://github.com/nais/salsa/Attestations/GitHubHostedActions@v1", env.BuilderId())
	assert.Equal(t, "https://github.com/nais/salsa/actions/runs/1839977840", env.BuildInvocationId())
	assert.Equal(t, json.RawMessage(nil), env.EventInputs())
	assert.Equal(t, "90dc9f2bc4007d1099a941ba3d408d2c896fe8dd", env.GithubSha())

}

func TestCreateRunnerContext(t *testing.T) {
	env := Environment{}
	err := ParseRunner(&runnerContext, &env)
	assert.NoError(t, err)
	assert.Equal(t, "Hosted Agent", env.RunnerContext.Name)
	assert.Equal(t, "Linux", env.RunnerContext.OS)
	assert.Equal(t, "X64", env.RunnerContext.Arch)
	assert.Equal(t, "/opt/hostedtoolcache", env.RunnerContext.ToolCache)
	assert.Equal(t, "/home/runner/work/_temp", env.RunnerContext.Temp)
}

var runnerContext = `{
		"os": "Linux",
		"arch": "X64",
		"name": "Hosted Agent",
		"tool_cache": "/opt/hostedtoolcache",
		"temp": "/home/runner/work/_temp",
		"workspace": "/home/runner/work/nais-salsa-action"
	  }`

func TestCreateCurrentEnvironmentContext(t *testing.T) {
	env := Environment{}
	expected := make(map[string]string)
	expected["GOVERSION"] = "1.17"
	expected["GOROOT"] = "/opt/hostedtoolcache/go/1.17.6/x64"
	err := ParseEnv(&envContext, &env)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(env.CurrentEnvironment.Envs))
}

var envContext = `{
  		"GOVERSION": "1.17",
		"GOROOT": "/opt/hostedtoolcache/go/1.17.6/x64"
	  }`
