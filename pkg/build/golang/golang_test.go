package golang

import (
	"github.com/nais/salsa/pkg/build/test"
	"testing"

	"github.com/nais/salsa/pkg/build"
)

func TestGoDeps(t *testing.T) {
	got := GoDeps(goSumContents)
	want := map[string]build.Dependency{}
	want["github.com/google/uuid"] = test.Dependency("github.com/google/uuid", "1.0.0", "sha256", "b4Gk+7WdP/d3HZH8EJsZpvV7EtDOgaZLtnaNGIu1adA=")
	want["github.com/pborman/uuid"] = test.Dependency("github.com/pborman/uuid", "1.2.1", "sha256", "+ZZIw58t/ozdjRaXh/3awHfmWRbzYxJoAdNJxe/3pvw=")

	test.AssertEqual(t, got, want)
}

const goSumContents = `
github.com/google/uuid v1.0.0 h1:b4Gk+7WdP/d3HZH8EJsZpvV7EtDOgaZLtnaNGIu1adA=
github.com/google/uuid v1.0.0/go.mod h1:TIyPZe4MgqvfeYDBFedMoGGpEw/LqOeaOT+nhxU+yHo=
github.com/pborman/uuid v1.2.1 h1:+ZZIw58t/ozdjRaXh/3awHfmWRbzYxJoAdNJxe/3pvw=
github.com/pborman/uuid v1.2.1/go.mod h1:X/NO0urCmaxf9VXbdlT7C2Yzkj2IKimNn4k+gtPdI/k=
`
