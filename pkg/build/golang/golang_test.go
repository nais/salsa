package golang

import (
	"github.com/nais/salsa/pkg/build"
	"github.com/nais/salsa/pkg/build/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGoDeps(t *testing.T) {
	got, err := GoDeps(goSumContents)
	assert.NoError(t, err)
	want := map[string]build.Dependency{}
	want["github.com/google/uuid"] = test.Dependency("github.com/google/uuid", "1.0.0", "sha256", "6f81a4fbb59d3ff7771d91fc109b19a6f57b12d0ce81a64bb6768d188bb569d0")
	want["github.com/pborman/uuid"] = test.Dependency("github.com/pborman/uuid", "1.2.1", "sha256", "f99648c39f2dfe8cdd8d169787fddac077e65916f363126801d349c5eff7a6fc")

	test.AssertEqual(t, got, want)
}

const goSumContents = `
github.com/google/uuid v1.0.0 h1:b4Gk+7WdP/d3HZH8EJsZpvV7EtDOgaZLtnaNGIu1adA=
github.com/google/uuid v1.0.0/go.mod h1:TIyPZe4MgqvfeYDBFedMoGGpEw/LqOeaOT+nhxU+yHo=
github.com/pborman/uuid v1.2.1 h1:+ZZIw58t/ozdjRaXh/3awHfmWRbzYxJoAdNJxe/3pvw=
github.com/pborman/uuid v1.2.1/go.mod h1:X/NO0urCmaxf9VXbdlT7C2Yzkj2IKimNn4k+gtPdI/k=
`

func TestBuildGo(t *testing.T) {
	tests := []test.IntegrationTest{
		{
			Name:      "find GOLANG build file and parse output",
			BuildType: BuildGo(),
			WorkDir:   "testdata/golang",
			BuildPath: "testdata/golang/go.sum",
			Cmd:       "go.sum",
			Want: test.Want{
				Key:     "github.com/Microsoft/go-winio",
				Version: "0.5.1",
				Algo:    "sha256",
				Digest:  "68f269d900fb38eae13b9b505ea42819225cf838c3b564c62ce98dc809ba1606",
			},
		},
		{
			Name:         "cant find GOLANG build file",
			BuildType:    BuildGo(),
			WorkDir:      "testdata/whatever",
			Error:        true,
			ErrorMessage: "could not find match, reading dir open testdata/whatever: no such file or directory",
		},
	}

	test.Run(t, tests)
}
