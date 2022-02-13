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
	context, err := CreateCIEnvironment(&parsedContext, &runnerContext)
	assert.NoError(t, err)

	assert.Equal(t, "ebe231e64736728ac2d6f3ae779fd29ad52d178f", context.GitHubContext.SHA)
	assert.Equal(t, "build", context.GitHubContext.Job)
	assert.Equal(t, "refs/heads/main", context.GitHubContext.Ref)
	assert.Equal(t, "nais/salsa", context.GitHubContext.Repository)
	assert.Equal(t, "nais", context.GitHubContext.RepositoryOwner)
	assert.Equal(t, "1691669140", context.GitHubContext.RunId)
	assert.Equal(t, "11", context.GitHubContext.RunNumber)
	assert.Equal(t, "jdoe", context.GitHubContext.Actor)

	assert.Equal(t, "https://github.com/nais/salsa", context.RepoUri())
	assert.Equal(t, "https://github.com/nais/salsa/Attestations/GitHubHostedActions@v1", context.BuilderId())
	assert.Equal(t, "https://github.com/nais/salsa/actions/runs/1691669140", context.BuildInvocationId())
	assert.Equal(t, json.RawMessage(nil), context.EventInputs())
	assert.Equal(t, "ebe231e64736728ac2d6f3ae779fd29ad52d178f", context.GithubSha())

}

func TestCreateRunnerContext(t *testing.T) {
	githubContext, err := os.ReadFile("testdata/github-context.json")
	assert.NoError(t, err)
	parsedContext := string(githubContext)
	context, err := CreateCIEnvironment(&parsedContext, &runnerContext)
	assert.NoError(t, err)
	assert.Equal(t, "Hosted Agent", context.RunnerContext.Name)
	assert.Equal(t, "Linux", context.RunnerContext.OS)
	assert.Equal(t, "X64", context.RunnerContext.Arch)
	assert.Equal(t, "/opt/hostedtoolcache", context.RunnerContext.ToolCache)
	assert.Equal(t, "/home/runner/work/_temp", context.RunnerContext.Temp)
}

var runnerContext = `{
		"os": "Linux",
		"arch": "X64",
		"name": "Hosted Agent",
		"tool_cache": "/opt/hostedtoolcache",
		"temp": "/home/runner/work/_temp",
		"workspace": "/home/runner/work/nais-salsa-action"
	  }`
