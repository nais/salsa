package vcs

import (
	"encoding/base64"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestResolveBuildContext(t *testing.T) {
	context, err := os.ReadFile("testdata/github-context.json")
	runner := RunnerTestContext
	assert.NoError(t, err)

	for _, test := range []struct {
		name         string
		context      string
		runner       string
		buildEnv     bool
		error        bool
		manually     bool
		errorMessage string
	}{
		{
			name:     "Resolve a CI environment with proper context and runner",
			context:  encode(context),
			runner:   encode([]byte(runner)),
			buildEnv: true,
		},
		{
			name:     "CLI is run manually without build context or runner",
			context:  "",
			runner:   "",
			manually: true,
		},
		{
			name:         "Not a valid build context",
			context:      "yolo",
			runner:       "yolo",
			error:        true,
			errorMessage: "decoded build context is not in json format",
		},
		{
			name:         "Valid input json context and runner but not a supported context",
			context:      encode([]byte(`{"valid": "json"}`)),
			runner:       encode([]byte(`{"valid": "json"}`)),
			error:        true,
			errorMessage: "build context is not supported",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			if test.buildEnv {
				err = os.Setenv("GITHUB_ACTIONS", "true")
				assert.NoError(t, err)
			}

			resolved, err := ResolveBuildContext(&test.context, &test.runner, nil)

			switch true {
			case test.error:
				assert.EqualError(t, err, test.errorMessage)
			case test.manually:
				assert.NoError(t, err)
				assert.Nil(t, resolved)
			default:
				assert.NoError(t, err)
				assert.NotNil(t, resolved)
				assert.Equal(t, "tomato CI", resolved.Context())
				assert.Equal(t, map[string]string{}, resolved.CurrentFilteredEnvironment())
				assert.Equal(t, "https://github.com/bolo/tomato/actions/runs/1839977840", resolved.BuildInvocationId())
				assert.Equal(t, "https://github.com/Attestations/GitHubActionsWorkflow@v1", resolved.BuildType())
				assert.Equal(t, "https://github.com/bolo/tomato", resolved.RepoUri())
				assert.Equal(t, "https://github.com/bolo/tomato/Attestations/GitHubHostedActions@v1", resolved.BuilderId())
				assert.Equal(t, "90dc9f2bc4007d1099a941ba3d408d2c896fe8dd", resolved.Sha())
				assert.NotNil(t, resolved.UserDefinedParameters())
				assert.NotEmpty(t, resolved.UserDefinedParameters())
				assert.Equal(t, "Linux", resolved.NonReproducibleMetadata().Context.Runner.Os)
				assert.Equal(t, "/home/runner/work/_temp", resolved.NonReproducibleMetadata().Context.Runner.Temp)
				assert.Equal(t, "1839977840", resolved.NonReproducibleMetadata().Context.Github.RunId)
				assert.Empty(t, "", resolved.NonReproducibleMetadata().Env)
				assert.Equal(t, "X64", resolved.NonReproducibleMetadata().Arch)
				// Unset for next test
				err = os.Unsetenv("GITHUB_ACTIONS")
				assert.NoError(t, err)
			}
		})
	}
}

func encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}
