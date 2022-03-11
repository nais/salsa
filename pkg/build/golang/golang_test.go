package golang

import (
	"github.com/nais/salsa/pkg/build/test"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/nais/salsa/pkg/build"
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
