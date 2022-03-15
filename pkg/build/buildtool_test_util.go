package build

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type Want struct {
	Key     string
	Version string
	Algo    string
	Digest  string
}

type IntegrationTest struct {
	Name         string
	BuildType    Tool
	WorkDir      string
	BuildPath    string
	Cmd          string
	Want         Want
	Error        bool
	ErrorMessage string
}

func RunTests(t *testing.T, tests []IntegrationTest) {
	for _, test := range tests {
		test.integrationTest(t)
	}
}

func (in IntegrationTest) integrationTest(t *testing.T) {
	t.Run(in.Name, func(t *testing.T) {
		tools := Tools{
			Tools: []Tool{in.BuildType},
		}

		// Check 1 random dependency is parsed dependencies.
		expected := map[string]Dependency{
			in.Want.Key: TestDependency(in.Want.Key, in.Want.Version, in.Want.Algo, in.Want.Digest),
		}

		deps, err := tools.DetectDeps(in.WorkDir)
		if in.Error {
			assert.EqualError(t, err, in.ErrorMessage)
		} else {
			assert.NoError(t, err)
			assert.NotNil(t, deps)
			assert.Equal(t, in.Cmd, deps.Cmd.CmdFlags)
			assert.NotEmpty(t, deps)
			assert.Equal(t, expected[in.Want.Key], deps.RuntimeDeps[in.Want.Key])
		}
	})
}

func TestDependency(coordinates, version, algo, checksum string) Dependency {
	return Dependence(coordinates, version,
		Verification(algo, checksum),
	)
}

func AssertEqual(t *testing.T, got, want map[string]Dependency) {
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %q, wanted %q", got, want)
	}
}
