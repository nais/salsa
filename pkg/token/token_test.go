package token

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestExportEnvironment(t *testing.T) {
	for _, exportTest := range []struct {
		name  string
		input string
		token string
		key   string
		error bool
	}{
		{
			name:  "export environment token with a single value in '--token' flag",
			input: "token",
			token: "token",
			key:   DefaultGithubTokenEnvKey,
			error: false,
		},
		{
			name:  "export environment token with a key=value in '--token' flag",
			input: "github.secret=specifiedKeyForTokenValue",
			token: "specifiedKeyForTokenValue",
			key:   "GITHUB_SECRET",
			error: false,
		},
		{
			name:  "no export of environment when no '--token' flag is specified",
			input: "",
			token: "",
			key:   DefaultGithubTokenEnvKey,
			error: false,
		},
		{
			name:  "wrong funcy format in '--token' flag",
			input: "=",
			token: "",
			key:   DefaultGithubTokenEnvKey,
			error: true,
		},
		{
			name:  "wrong funcy format 2 in '--token' flag",
			input: "tokenkey=tokenvalue=",
			token: "",
			key:   DefaultGithubTokenEnvKey,
			error: true,
		},
	} {
		t.Run(exportTest.name, func(t *testing.T) {
			err := Export(exportTest.input)

			if exportTest.error {
				assert.EqualError(t, err, "wrong format for --token flag; should be a single value or 'key=value'")
			} else {
				assert.Equal(t, exportTest.token, os.Getenv(exportTest.key))
				// unset for next test
				err = os.Unsetenv(exportTest.key)
				assert.NoError(t, err)
			}
		})
	}
}
