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
	context, err := CreateCIEnvironment(&s)
	assert.NoError(t, err)
	assert.Equal(t, "ebe231e64736728ac2d6f3ae779fd29ad52d178f", context.GitHubContext.SHA)
}
