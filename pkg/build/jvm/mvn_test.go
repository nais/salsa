package jvm

import (
	"fmt"
	"github.com/nais/salsa/pkg/build"
	"github.com/nais/salsa/pkg/build/test"
	"github.com/nais/salsa/pkg/token"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestMavenDeps(t *testing.T) {
	got, err := MavenCompileAndRuntimeTimeDeps("testdata/target/dependency")
	assert.NoError(t, err)
	want := map[string]build.Dependency{}
	want["org.springframework:spring-core"] = test.Dependency(
		"org.springframework:spring-core",
		"5.3.16", "sha256",
		"0903d17e58654a2c79f4e46df79dc73ccaa49b6edbc7c3278359db403b687f6e",
	)
	want["org.yaml:snakeyaml"] = test.Dependency(
		"org.yaml:snakeyaml",
		"1.26",
		"sha256",
		"d87d607e500885356c03c1cae61e8c2e05d697df8787d5aba13484c2eb76a844",
	)

	test.AssertEqual(t, got, want)
}

func TestMavenMvnCmd(t *testing.T) {
	for _, mvnTest := range []struct {
		name  string
		token string
		key   string
		cmd   string
	}{
		{
			name:  "generate a maven cmd from input with empty '--token' flag",
			token: "",
			key:   token.DefaultGithubTokenEnvKey,
			cmd:   "/usr/local/bin/mvn dependency:copy-dependencies -DincludeScope=runtime -Dmdep.useRepositoryLayout=true",
		},
		{
			name:  "generate a maven cmd from input token",
			token: "token",
			key:   token.DefaultGithubTokenEnvKey,
			cmd:   "/usr/local/bin/mvn dependency:copy-dependencies -DincludeScope=runtime -Dmdep.useRepositoryLayout=true",
		},
	} {
		t.Run(mvnTest.name, func(t *testing.T) {
			mvn := Maven{
				BuildFilePatterns: nil,
				Settings: Settings{
					Auth: Auth{
						GithubToken: mvnTest.token,
					},
				},
			}

			cmd, err := mvn.mvnCmd()

			assert.Equal(t, mvnTest.cmd, fmt.Sprintf("%s", cmd))
			assert.Equal(t, mvnTest.token, os.Getenv(mvnTest.key))
			// unset for next test
			err = os.Unsetenv(mvnTest.key)
			assert.NoError(t, err)
		})
	}
}
