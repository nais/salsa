package github

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestParseGithubContext(t *testing.T) {
	githubContext, err := os.ReadFile("../testdata/github-context.json")
	assert.NoError(t, err)
	context, err := ParseContext(githubContext)
	assert.NoError(t, err)

	assert.Equal(t, "90dc9f2bc4007d1099a941ba3d408d2c896fe8dd", context.SHA)
	assert.Equal(t, "build", context.Job)
	assert.Equal(t, "refs/heads/main", context.Ref)
	assert.Equal(t, "nais/salsa", context.Repository)
	assert.Equal(t, "nais", context.RepositoryOwner)
	assert.Equal(t, "1839977840", context.RunId)
	assert.Equal(t, "57", context.RunNumber)
	assert.Equal(t, "jdoe", context.Actor)
}
