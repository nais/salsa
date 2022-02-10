package vcs

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestCreateGithubContext(t *testing.T) {
	githubContext, err := os.ReadFile("testdata/github-context.json")
	assert.NoError(t, err)
	s := string(githubContext)
	context, err := CreateCIEnvironment(&s, nil)
	assert.NoError(t, err)

	assert.Equal(t, "ebe231e64736728ac2d6f3ae779fd29ad52d178f", context.GitHubContext.SHA)
	assert.Equal(t, "build", context.GitHubContext.Job)
	assert.Equal(t, "refs/heads/main", context.GitHubContext.Ref)
	assert.Equal(t, "nais/salsa", context.GitHubContext.Repository)
	assert.Equal(t, "nais", context.GitHubContext.RepositoryOwner)
	assert.Equal(t, "1691669140", context.GitHubContext.RunId)
	assert.Equal(t, "11", context.GitHubContext.RunNumber)
	assert.Equal(t, "jdoe", context.GitHubContext.Actor)
	outputJson, _ := context.GitHubContext.Event.MarshalJSON()
	assert.NotEmpty(t, string(outputJson))
}

func TestCreateRunnerContext(t *testing.T) {
	context, err := CreateCIEnvironment(nil, &runnerContext)
	assert.NoError(t, err)
	assert.Equal(t, "Linux", context.RunnerContext.OS)
	assert.Equal(t, "/opt/hostedtoolcache", context.RunnerContext.ToolCache)
	assert.Equal(t, "/home/runner/work/_temp", context.RunnerContext.Temp)
}

var runnerContext = `{
		"os": "Linux",
		"name": "Hosted Agent",
		"tool_cache": "/opt/hostedtoolcache",
		"temp": "/home/runner/work/_temp",
		"workspace": "/home/runner/work/nais-salsa-action"
	  }`
